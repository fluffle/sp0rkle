package db

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"

	"github.com/fluffle/golog/logging"
	"go.etcd.io/bbolt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// A value that is stored directly at K{{"_id", ObjectId}}
// in the _vals bucket, with pointers for each Key in _idxs.
type Indexer interface {
	Id() bson.ObjectId
	Indexes() []Key
}

// Per https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
var indexerType reflect.Type = reflect.TypeOf((*Indexer)(nil)).Elem()

func isPointer(data []byte) bool {
	if len(data) < prefixLen {
		return false
	}
	return bytes.Equal(data[:prefixLen], idPrefix)
}

func toPointer(value Indexer) []byte {
	e := S{idTag, string(value.Id())}
	return e.Bytes()
}

func (b *boltDatabase) Indexed() Database {
	b.Lock()
	defer b.Unlock()
	if b.db == nil {
		logging.Fatal("Tried to create BoltDB indexed database when disconnected.")
	}
	return &indexedDatabase{db: b.db}
}

type indexedDatabase struct {
	db *bbolt.DB
}

func (i *indexedDatabase) C(name string) Collection {
	vals := append([]byte(name), []byte("_vals")...)
	idxs := append([]byte(name), []byte("_idxs")...)

	err := i.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(vals)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(idxs)
		return err
	})
	if err != nil {
		logging.Fatal("Creating BoltDB bucket failed: %v")
	}
	return &indexedBucket{name: name, vals: vals, idxs: idxs, db: i.db}
}

type indexedBucket struct {
	name   string
	vals   []byte
	idxs   []byte
	db     *bbolt.DB
	debug_ bool
}

func (bucket *indexedBucket) Debug(on bool) {
	bucket.debug_ = on
}

func (bucket *indexedBucket) debug(f string, args ...interface{}) {
	if bucket.debug_ {
		logging.Debug("%s."+f, append([]interface{}{bucket.name}, args...)...)
	}
}

func (bucket *indexedBucket) error(f string, args ...interface{}) error {
	return fmt.Errorf("%s."+f, append([]interface{}{bucket.name}, args...)...)
}

func (bucket *indexedBucket) values(tx *bbolt.Tx) *bbolt.Bucket {
	return tx.Bucket(bucket.vals)
}

func (bucket *indexedBucket) find(tx *bbolt.Tx, elems [][]byte) *bbolt.Bucket {
	b := tx.Bucket(bucket.idxs)
	for _, elem := range elems {
		if b = b.Bucket(elem); b == nil {
			bucket.debug("find(): bucket %q not found", elem)
			return nil
		}
	}
	return b
}

func (bucket *indexedBucket) create(tx *bbolt.Tx, elems [][]byte) (*bbolt.Bucket, error) {
	b := tx.Bucket(bucket.idxs)
	var err error
	for _, elem := range elems {
		if b, err = b.CreateBucketIfNotExists(elem); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (bucket *indexedBucket) Get(key Key, value interface{}) error {
	elems, last := key.B()
	if len(last) == 0 {
		return bucket.error("Get(): zero length key")
	}

	return bucket.db.View(func(tx *bbolt.Tx) error {
		if len(elems) >= 0 || !isPointer(last) {
			b := bucket.find(tx, elems)
			if b == nil {
				return nil
			}
			last = b.Get(last)
			if last == nil {
				return nil
			}
		}
		data := bucket.values(tx).Get(last)
		bucket.debug("Get(%s) = %q", key, data)
		if data == nil {
			return nil
		}
		return bson.Unmarshal(suffix(data), value)
	})
}

func (bucket *indexedBucket) All(key Key, value interface{}) error {
	elems, last := key.B()
	if len(last) == 0 {
		// A zero-length key will perform a scan over the vals bucket directly,
		// since this conveniently contains all the real data keyed by ID.
		scanner := allScanner{
			sp: newSlicePtr(value),
		}
		return bucket.db.View(func(tx *bbolt.Tx) error {
			err := scanTx(bucket.values(tx), scanner)
			bucket.debug("%s: found %d keys", scanner, scanner.sp.len())
			return err
		})
	}
	// All implies that the last key elem is also a bucket.
	elems = append(elems, last)
	return bucket.db.View(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, elems)
		if b == nil {
			return nil
		}
		scanner := indexScanner{
			sp:   newSlicePtr(value),
			vals: bucket.values(tx),
			seen: map[string]bool{},
		}
		err := scanTx(b, scanner)
		bucket.debug("%s: found %d keys", scanner, scanner.sp.len())
		return err
	})
}

func (bucket *indexedBucket) Match(field, re string, value interface{}) error {
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
		// Match always scans across all values.
		err := scanTx(bucket.values(tx), scanner)
		bucket.debug("%s: found %d keys", scanner, scanner.sp.len())
		return err
	})
}

func (bucket *indexedBucket) Put(value interface{}) error {
	indexer, ok := value.(Indexer)
	if !ok {
		return bucket.error("Put(): don't know how to put value %#v", value)
	}
	data, err := toBson(indexer)
	if err != nil {
		return err
	}
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		return bucket.putTx(tx, indexer, data)
	})
}

func (bucket *indexedBucket) BatchPut(value interface{}) error {
	// vv == value Value
	vv := reflect.ValueOf(value)
	if vv.Kind() != reflect.Slice || !vv.Type().Elem().Implements(indexerType) {
		return bucket.error("BatchPut(): can only put a slice of Indexers")
	}

	type kvTuple struct {
		value Indexer
		data  []byte
	}
	tuples := make([]kvTuple, vv.Len())

	for i := 0; i < vv.Len(); i++ {
		indexer, _ := vv.Index(i).Interface().(Indexer)
		data, err := toBson(vv.Index(i).Interface())
		if err != nil {
			return err
		}
		tuples[i] = kvTuple{indexer, data}
	}
	bucket.debug("BatchPut(): serialized %d items", len(tuples))

	return bucket.db.Update(func(tx *bbolt.Tx) error {
		for _, tuple := range tuples {
			if err := bucket.putTx(tx, tuple.value, tuple.data); err != nil {
				return err
			}
		}
		bucket.debug("BatchPut(): put %d items", len(tuples))
		return nil
	})
}

func (bucket *indexedBucket) putTx(tx *bbolt.Tx, value Indexer, data []byte) error {
	ptr := toPointer(value)
	v := bucket.values(tx).Get(ptr)
	if isBson(v) {
		// There's already a value here, probably being pointed at.
		// Jump through some hoops to clean up those index pointers.
		// TODO(fluffle): This makes some assumptions that may not
		// hold true, and might leave dangling index pointers, ugh.
		//   1) The old value is of the same type as the new one.
		//   2) The indexes derived from the old data are exactly
		//      the correct set that should be deleted to tidy up.
		//   3) We don't need to recursively clean up empty nested buckets.
		old := dupe(value).(Indexer)
		if err := bson.Unmarshal(suffix(v), old); err != nil {
			return err
		}
		if err := bucket.delIndex(tx, old); err != nil {
			return err
		}
	}
	bucket.debug("Put(%s) = %q", value.Id(), data)
	if err := bucket.values(tx).Put(ptr, data); err != nil {
		return err
	}
	return bucket.putIndex(tx, value)
}

func (bucket *indexedBucket) putIndex(tx *bbolt.Tx, value Indexer) error {
	ptr := toPointer(value)
	for _, key := range value.Indexes() {
		elems, last := key.B()
		b, err := bucket.create(tx, elems)
		if err != nil {
			return err
		}
		if err = b.Put(last, ptr); err != nil {
			return err
		}
		bucket.debug("putIndex(%s) = %q", key, ptr)
	}
	return nil
}

func (bucket *indexedBucket) delIndex(tx *bbolt.Tx, value Indexer) error {
	ptr := toPointer(value)
	for _, key := range value.Indexes() {
		elems, last := key.B()
		b := bucket.find(tx, elems)
		if b == nil {
			return nil
		}
		if err := b.Delete(last); err != nil {
			return err
		}
		bucket.debug("delIndex(%s) = %q", key, ptr)
	}
	return nil
}

func (bucket *indexedBucket) Del(value interface{}) error {
	indexer, ok := value.(Indexer)
	if !ok {
		return bucket.error("Del(): don't know how to delete value %#v", value)
	}
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		if err := bucket.values(tx).Delete(toPointer(indexer)); err != nil {
			return err
		}
		bucket.debug("Del(%s)", indexer.Id())
		return bucket.delIndex(tx, indexer)
	})
}

func (bucket *indexedBucket) Next(k Key, set ...int) (int, error) {
	var i uint64
	elems, last := k.B()
	// Next implies that the last key elem is also a bucket.
	if len(last) > 0 {
		elems = append(elems, last)
	}
	err := bucket.db.Update(func(tx *bbolt.Tx) error {
		// The empty key will increment the counter for the values
		// bucket, non-empty keys will be in the index buckets.
		b := bucket.values(tx)
		if len(elems) > 0 {
			b = bucket.find(tx, elems)
		}
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

func (bucket *indexedBucket) Mongo() *mgo.Collection {
	panic("you are bad at migrations")
}
