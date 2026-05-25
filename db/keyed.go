package db

import (
	"fmt"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/util/bson"
	"go.etcd.io/bbolt"
)

// A value that is stored directly at Key in BoltDB.
// The method is not called Key because conf.Entry has
// a field named Key which references data in mongo
// but still needs to implement this interface.
// Naming is hard, but this is probably fine because
// they will most likely be returning a db.K anyway.
type Keyer interface {
	K() Key
	CollectionName() string
}

// RegisterKeyed creates a new keyed collection for the given type T.
// Requires *boltDatabase specifically (not the Database interface) because
// it needs to create buckets in the BoltDB instance.
func RegisterKeyed[T Keyer](b *boltDatabase) KeyedCollection[T] {
	b.Lock()
	defer b.Unlock()
	if b.db == nil {
		logging.Fatal("Tried to create BoltDB keyed database when disconnected.")
	}
	var zero T
	n := []byte(zero.CollectionName())
	err := b.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(n)
		return err
	})
	if err != nil {
		logging.Fatal("Creating BoltDB bucket failed: %v", err)
	}
	return &keyedBucket[T]{name: n, db: b.db}
}

type keyedBucket[T Keyer] struct {
	name   []byte
	db     *bbolt.DB
	debug_ bool
}

func (bucket *keyedBucket[T]) Debug(on bool) {
	bucket.debug_ = on
}

func (bucket *keyedBucket[T]) debug(f string, args ...any) {
	if bucket.debug_ {
		logging.Debug("%s."+f, append([]any{bucket.name}, args...)...)
	}
}

func (bucket *keyedBucket[T]) error(f string, args ...any) error {
	return fmt.Errorf("%s."+f, append([]any{bucket.name}, args...)...)
}

func (bucket *keyedBucket[T]) find(tx *bbolt.Tx, elems [][]byte) *valsBucket {
	b := tx.Bucket(bucket.name)
	for _, elem := range elems {
		if b = b.Bucket(elem); b == nil {
			bucket.debug("find(): bucket %q not found", elem)
			return nil
		}
	}
	return &valsBucket{Bucket: b}
}

func (bucket *keyedBucket[T]) create(tx *bbolt.Tx, elems [][]byte) (*valsBucket, error) {
	b := tx.Bucket(bucket.name)
	var err error
	for _, elem := range elems {
		if b, err = b.CreateBucketIfNotExists(elem); err != nil {
			return nil, fmt.Errorf("create bucket %q: %w", elem, err)
		}
	}
	return &valsBucket{Bucket: b}, nil
}

func (bucket *keyedBucket[T]) Get(key Key) (T, error) {
	var zero T
	elems, last := key.B()
	if len(last) == 0 {
		return zero, bucket.error("Get(): zero length key")
	}
	var result T
	var err error
	err = bucket.db.View(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, elems)
		if b == nil {
			return nil
		}
		data := b.Get(last)
		bucket.debug("Get(%s) = %q", key, data)
		if data == nil {
			return nil
		}
		return bson.Unmarshal(suffix(data), &result)
	})
	if err != nil {
		return zero, err
	}
	return result, nil
}

func (bucket *keyedBucket[T]) All(key Key) ([]T, error) {
	elems, last := key.B()
	if len(last) > 0 {
		elems = append(elems, last)
	}
	var result []T
	err := bucket.db.View(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, elems)
		if b == nil {
			return nil
		}
		var err error
		result, err = scanAll[T](b)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (bucket *keyedBucket[T]) Match(pred func(T) bool) ([]T, error) {
	var result []T
	err := bucket.db.View(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, nil)
		if b == nil {
			return nil
		}
		var err error
		result, err = scanMatch[T](b, pred)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (bucket *keyedBucket[T]) Put(value T) error {
	key := value.K()
	if err := key.Valid(); err != nil {
		return bucket.error("Put(): invalid key for %T: %w", value, err)
	}
	elems, last := key.B()
	if len(last) == 0 {
		return bucket.error("Put(): can't put value with empty key")
	}
	data, err := toBson(value)
	if err != nil {
		return err
	}
	bucket.debug("Put(%s) = %q", value.K(), data)
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		return bucket.putTx(tx, elems, last, data)
	})
}

func (bucket *keyedBucket[T]) BatchPut(values []T) error {
	type kvTuple struct {
		elems  [][]byte
		last   []byte
		data   []byte
	}
	tuples := make([]kvTuple, len(values))

	for i, v := range values {
		key := v.K()
		if err := key.Valid(); err != nil {
			return bucket.error("BatchPut(): invalid key for %T: %w", v, err)
		}
		elems, last := key.B()
		if len(last) == 0 {
			return bucket.error("BatchPut(): can't put value with empty key")
		}
		data, err := toBson(v)
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

func (bucket *keyedBucket[T]) putTx(tx *bbolt.Tx, elems [][]byte, key, value []byte) error {
	b, err := bucket.create(tx, elems)
	if err != nil {
		return err
	}
	return b.Put(key, value)
}

func (bucket *keyedBucket[T]) Del(value T) error {
	key := value.K()
	if err := key.Valid(); err != nil {
		return bucket.error("Del(): invalid key for %T: %w", value, err)
	}
	elems, last := key.B()
	if len(last) == 0 {
		return bucket.error("Del(): refusing to delete everything")
	}
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, elems)
		if b == nil {
			return nil
		}
		if b.Bucket.Bucket(last) != nil {
			// Allow partial keys to recursively delete nested buckets.
			return b.DeleteBucket(last)
		}
		return b.Delete(last)
	})
}

func (bucket *keyedBucket[T]) Fsck() error {
	all, err := bucket.All(K{})
	if err != nil {
		return err
	}
	if len(all) == 0 {
		return nil
	}

	var zero T
	vv, ok := any(zero).(ValueValidator[T])
	if !ok {
		return nil
	}
	del, err := vv.ValidateValues(all)
	if err != nil {
		return err
	}
	if len(del) == 0 {
		return nil
	}
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		for _, d := range del {
			k := d.K()
			elems, last := k.B()
			if len(last) == 0 {
				return bucket.error("Del(): refusing to delete everything")
			}
			b := bucket.find(tx, elems)
			if b == nil {
				continue
			}
			if b.Bucket.Bucket(last) != nil {
				// Allow partial keys to recursively delete nested buckets.
				if err := b.DeleteBucket(last); err != nil {
					return err
				}
			} else {
				if err := b.Delete(last); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (bucket *keyedBucket[T]) Next(k Key, set ...int) (int, error) {
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
