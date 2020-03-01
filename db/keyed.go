package db

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/fluffle/golog/logging"
	"go.etcd.io/bbolt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// A value that is stored directly at Key in BoltDB.
// The method is not called Key because conf.Entry has
// a field named Key which references data in mongo
// but still needs to implement this interface.
// Naming is hard, but this is probably fine because
// they will most likely be returning a db.K anyway.
type Keyer interface {
	K() Key
}

// Per https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
var keyerType reflect.Type = reflect.TypeOf((*Keyer)(nil)).Elem()

func (b *boltDatabase) Keyed() Database {
	b.Lock()
	defer b.Unlock()
	if b.db == nil {
		logging.Fatal("Tried to create BoltDB keyed database when disconnected.")
	}
	return &keyedDatabase{db: b.db}
}

type keyedDatabase struct {
	db *bbolt.DB
}

func (k *keyedDatabase) C(name string) Collection {
	n := []byte(name)
	err := k.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(n)
		return err
	})
	if err != nil {
		logging.Fatal("Creating BoltDB bucket failed: %v")
	}
	return &keyedBucket{name: n, db: k.db}
}

type keyedBucket struct {
	name   []byte
	db     *bbolt.DB
	debug_ bool
}

func (bucket *keyedBucket) Debug(on bool) {
	bucket.debug_ = on
}

func (bucket *keyedBucket) debug(f string, args ...interface{}) {
	if bucket.debug_ {
		logging.Debug("%s."+f, append([]interface{}{bucket.name}, args...)...)
	}
}

func (bucket *keyedBucket) error(f string, args ...interface{}) error {
	return fmt.Errorf("%s."+f, append([]interface{}{bucket.name}, args...)...)
}

func (bucket *keyedBucket) find(tx *bbolt.Tx, elems [][]byte) *bbolt.Bucket {
	b := tx.Bucket(bucket.name)
	for _, elem := range elems {
		if b = b.Bucket(elem); b == nil {
			bucket.debug("find(): bucket %q not found", elem)
			return nil
		}
	}
	return b
}

func (bucket *keyedBucket) create(tx *bbolt.Tx, elems [][]byte) (*bbolt.Bucket, error) {
	b := tx.Bucket(bucket.name)
	var err error
	for _, elem := range elems {
		if b, err = b.CreateBucketIfNotExists(elem); err != nil {
			return nil, fmt.Errorf("create bucket %q: %w", elem, err)
		}
	}
	return b, nil
}

func (bucket *keyedBucket) Get(key Key, value interface{}) error {
	elems, last := key.B()
	if len(last) == 0 {
		return bucket.error("Get(): zero length key")
	}
	return bucket.db.View(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, elems)
		if b == nil {
			return nil
		}
		data := b.Get(last)
		bucket.debug("Get(%s) = %q", key, data)
		if data == nil {
			return nil
		}
		return bson.Unmarshal(suffix(data), value)
	})
}

func (bucket *keyedBucket) All(key Key, value interface{}) error {
	elems, last := key.B()
	// All implies that the last key elem is also a bucket.
	// We support a zero-length key to perform a scan over the root bucket.
	if len(last) > 0 {
		elems = append(elems, last)
	}
	scanner := allScanner{
		sp: newSlicePtr(value),
	}

	return bucket.db.View(func(tx *bbolt.Tx) error {
		if b := bucket.find(tx, elems); b != nil {
			err := scanTx(b, scanner)
			bucket.debug("%s: found %d keys", scanner, scanner.sp.len())
			return err
		}
		return nil
	})
}

func (bucket *keyedBucket) Match(field, re string, value interface{}) error {
	if re == "" {
		return bucket.error("Match(): zero-length regex match")
	}
	rx, err := regexp.Compile("(?i)" + re)
	if err != nil {
		return err
	}
	scanner := matchScanner{
		re:    re,
		rx:    rx,
		sp:    newSlicePtr(value),
		field: field,
	}

	// The slice elements may be pointers, we need the struct.
	cev := scanner.sp.newStruct()
	if cev.Kind() != reflect.Struct || cev.FieldByName(field).Kind() != reflect.String {
		return bucket.error("Match(): value kind is %s not struct, or field %s is kind %s not string (%#v)",
			cev.Kind(), field, cev.FieldByName(field).Kind(), value)
	}

	return bucket.db.View(func(tx *bbolt.Tx) error {
		if b := bucket.find(tx, nil); b != nil {
			err := scanTx(b, scanner)
			bucket.debug("%s: found %d keys", scanner, scanner.sp.len())
			return err
		}
		return nil
	})
}

func (bucket *keyedBucket) Put(value interface{}) error {
	keyer, ok := value.(Keyer)
	if !ok {
		return bucket.error("Put(): don't know how to put value %#v", value)
	}
	elems, last := keyer.K().B()
	if len(last) == 0 {
		return bucket.error("Put(): can't put value with empty key")
	}
	data, err := toBson(value)
	if err != nil {
		return err
	}
	bucket.debug("Put(%s) = %q", keyer.K(), data)
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		return bucket.putTx(tx, elems, last, data)
	})
}

func (bucket *keyedBucket) BatchPut(value interface{}) error {
	// vv == value Value
	vv := reflect.ValueOf(value)
	if vv.Kind() != reflect.Slice || !vv.Type().Elem().Implements(keyerType) {
		return bucket.error("BatchPut(): can only put a slice of Keyers")
	}

	// Do as much work as possible before the transaction.
	type kvTuple struct {
		elems      [][]byte
		last, data []byte
	}
	tuples := make([]kvTuple, vv.Len())

	for i := 0; i < vv.Len(); i++ {
		keyer, _ := vv.Index(i).Interface().(Keyer)
		elems, last := keyer.K().B()
		if len(last) == 0 {
			return bucket.error("BatchPut(): can't put value with empty key")
		}
		data, err := toBson(vv.Index(i).Interface())
		if err != nil {
			return err
		}
		tuples[i] = kvTuple{elems, last, data}
	}
	bucket.debug("BatchPut(): serialized %d items", len(tuples))

	return bucket.db.Update(func(tx *bbolt.Tx) error {
		for _, tuple := range tuples {
			if err := bucket.putTx(tx, tuple.elems, tuple.last, tuple.data); err != nil {
				return fmt.Errorf("BatchPut(%q): %w", tuple.last, err)
			}
		}
		bucket.debug("BatchPut(): put %d items", len(tuples))
		return nil
	})
}

func (bucket *keyedBucket) putTx(tx *bbolt.Tx, elems [][]byte, key, value []byte) error {
	b, err := bucket.create(tx, elems)
	if err != nil {
		return err
	}
	return b.Put(key, value)
}

func (bucket *keyedBucket) Del(value interface{}) error {
	keyer, ok := value.(Keyer)
	if !ok {
		return bucket.error("Del(): don't know how to delete value %#v", value)
	}
	elems, last := keyer.K().B()
	if len(last) == 0 {
		return bucket.error("Del(): refusing to delete everything")
	}
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, elems)
		if b == nil {
			// Parent bucket already doesn't exist.
			return nil
		}
		// Allow partial keys to recursively delete nested buckets.
		if b.Bucket(last) != nil {
			return b.DeleteBucket(last)
		}
		return b.Delete(last)
	})
}

func (bucket *keyedBucket) Next(k Key, set ...int) (int, error) {
	var i uint64
	elems, last := k.B()
	// Next implies that the last key elem is also a bucket.
	if len(last) > 0 {
		elems = append(elems, last)
	}
	err := bucket.db.Update(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, elems)
		if b == nil {
			return bbolt.ErrBucketNotFound
		}

		var err error
		if len(set) > 0 {
			i, err = uint64(set[0]), b.SetSequence(uint64(set[0]))
		} else {
			i, err = b.NextSequence()
		}
		return err
	})
	return int(i), err
}

func (bucket *keyedBucket) Mongo() *mgo.Collection {
	panic("you are bad at migrations")
}
