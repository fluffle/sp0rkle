package db

import (
	"bytes"
	"fmt"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/util/bson"
	"go.etcd.io/bbolt"
)

// A value that is stored directly at K{{"_id", ObjectId}}
// in the _vals bucket, with pointers for each Key in _idxs.
type Indexer interface {
	Id() bson.ObjectId
	Indexes() []Key
	CollectionName() string
}

// RegisterIndexed creates a new indexed collection for the given type T.
// Requires *boltDatabase specifically (not the Database interface) because
// it needs to create buckets in the BoltDB instance.
func RegisterIndexed[T Indexer](b *boltDatabase) IndexedCollection[T] {
	b.Lock()
	defer b.Unlock()
	if b.db == nil {
		logging.Fatal("Tried to create BoltDB indexed database when disconnected.")
	}
	var zero T
	n := []byte(zero.CollectionName())
	vals := append(n, []byte("_vals")...)
	idxs := append(n, []byte("_idxs")...)

	err := b.db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(vals); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(idxs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logging.Fatal("Creating BoltDB bucket for collection %q failed: %v", zero.CollectionName(), err)
	}
	bucket := &indexedBucket[T]{name: n, vals: vals, idxs: idxs, db: b.db}
	return bucket
}

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

// indexedBucket is a generic, type-safe indexed store.
type indexedBucket[T Indexer] struct {
	name   []byte
	vals   []byte
	idxs   []byte
	db     *bbolt.DB
	debug_ bool
}

func (bucket *indexedBucket[T]) Debug(on bool) {
	bucket.debug_ = on
}

func (bucket *indexedBucket[T]) debug(f string, args ...any) {
	if bucket.debug_ {
		logging.Debug("%s."+f, append([]any{string(bucket.name)}, args...)...)
	}
}

func (bucket *indexedBucket[T]) error(f string, args ...any) error {
	return fmt.Errorf("%s."+f, append([]any{string(bucket.name)}, args...)...)
}

func (bucket *indexedBucket[T]) values(tx *bbolt.Tx) *valsBucket {
	return &valsBucket{Bucket: tx.Bucket(bucket.vals)}
}

func (bucket *indexedBucket[T]) indexes(tx *bbolt.Tx) *idxsBucket {
	return &idxsBucket{Bucket: tx.Bucket(bucket.idxs)}
}

func (bucket *indexedBucket[T]) find(tx *bbolt.Tx, elems [][]byte) *idxsBucket {
	b := tx.Bucket(bucket.idxs)
	for _, elem := range elems {
		if b = b.Bucket(elem); b == nil {
			bucket.debug("find(): bucket %q not found", elem)
			return nil
		}
	}
	return &idxsBucket{Bucket: b}
}

func (bucket *indexedBucket[T]) create(tx *bbolt.Tx, elems [][]byte) (*idxsBucket, error) {
	b := tx.Bucket(bucket.idxs)
	var err error
	for _, elem := range elems {
		if b, err = b.CreateBucketIfNotExists(elem); err != nil {
			return nil, err
		}
	}
	return &idxsBucket{Bucket: b}, nil
}

func (bucket *indexedBucket[T]) Get(key Key) (T, error) {
	var zero T
	elems, last := key.B()
	if len(last) == 0 {
		return zero, bucket.error("Get(): zero length key")
	}
	var result T
	err := bucket.db.View(func(tx *bbolt.Tx) error {
		bucket.debug("Get(%s) looking up bucket key %q", key, last)
		// Dual-path pointer resolution:
		// - If there are index elements to traverse (elems), look up in the index bucket.
		// - If the key is not already a pointer (isPointer check), also look it up in the index bucket.
		// - Otherwise, the key is a direct _id pointer and we skip the index lookup.
		if len(elems) > 0 || !isPointer(last) {
			b := bucket.find(tx, elems)
			if b == nil {
				return nil
			}
			last = b.Get(last)
			bucket.debug("Get(%s) pointer = %q", key, last)
			if last == nil {
				return nil
			}
		}
		data := bucket.values(tx).Get(last)
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

func (bucket *indexedBucket[T]) All(key Key) ([]T, error) {
	var result []T
	elems, last := key.B()
	if len(last) == 0 {
		err := bucket.db.View(func(tx *bbolt.Tx) error {
			var err error
			result, err = scanAll[T](bucket.values(tx))
			return err
		})
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	elems = append(elems, last)
	err := bucket.db.View(func(tx *bbolt.Tx) error {
		b := bucket.find(tx, elems)
		if b == nil {
			return nil
		}
		var err error
		result, err = scanIndexedAll[T](b, bucket.values(tx))
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (bucket *indexedBucket[T]) Match(pred func(T) bool) ([]T, error) {
	var result []T
	err := bucket.db.View(func(tx *bbolt.Tx) error {
		var err error
		result, err = scanMatch[T](bucket.values(tx), pred)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (bucket *indexedBucket[T]) Put(value T) error {
	data, err := toBson(value)
	if err != nil {
		return err
	}
	for _, key := range value.Indexes() {
		if err := key.Valid(); err != nil {
			return bucket.error("Put(): invalid key for %T: %w", value, err)
		}
	}
	if !value.Id().Valid() {
		return bucket.error("Put(): invalid ObjectId for %T", value)
	}
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		return bucket.putTx(tx, value, data)
	})
}

func (bucket *indexedBucket[T]) BatchPut(values []T) error {
	type kvTuple struct {
		value T
		data  []byte
	}
	tuples := make([]kvTuple, len(values))

	for i, v := range values {
		for _, key := range v.Indexes() {
			if err := key.Valid(); err != nil {
				return bucket.error("BatchPut(): invalid key for %T: %w", v, err)
			}
		}
		data, err := toBson(v)
		if err != nil {
			return err
		}
		tuples[i] = kvTuple{v, data}
	}
	bucket.debug("BatchPut(): serialized %d items", len(tuples))

	return bucket.db.Update(func(tx *bbolt.Tx) error {
		for _, tuple := range tuples {
			if err := bucket.putTx(tx, tuple.value, tuple.data); err != nil {
				return fmt.Errorf("BatchPut(%s): %w", tuple.value.Id(), err)
			}
		}
		bucket.debug("BatchPut(): put %d items", len(tuples))
		return nil
	})
}

func (bucket *indexedBucket[T]) Del(value T) error {
	for _, key := range value.Indexes() {
		if err := key.Valid(); err != nil {
			return bucket.error("Del(): invalid key for %T: %w", value, err)
		}
	}
	return bucket.db.Update(func(tx *bbolt.Tx) error {
		if err := bucket.values(tx).Delete(toPointer(value)); err != nil {
			return err
		}
		bucket.debug("Del(%s)", value.Id())
		return bucket.delIndex(tx, value)
	})
}

func (bucket *indexedBucket[T]) Fsck() error {
	all, err := bucket.All(K{})
	if err != nil {
		return err
	}
	if len(all) > 0 {
		var zero T
		if vv, ok := interface{}(zero).(ValueValidator[T]); ok {
			del, err := vv.ValidateValues(all)
			if err != nil {
				return err
			}
			if len(del) > 0 {
				return bucket.db.Update(func(tx *bbolt.Tx) error {
					vals := bucket.values(tx)
					for _, d := range del {
						if err := vals.Delete(toPointer(d)); err != nil {
							return err
						}
						if err := bucket.delIndex(tx, d); err != nil {
							return err
						}
					}
					return nil
				})
			}
		}
	}

	return bucket.db.Update(func(tx *bbolt.Tx) error {
		vals := bucket.values(tx)
		idxs := bucket.indexes(tx)
		if err := fsckIndexPass[T](idxs, vals); err != nil {
			return err
		}
		if err := fsckValuePass[T](idxs, vals); err != nil {
			return err
		}
		return nil
	})
}

func (bucket *indexedBucket[T]) putTx(tx *bbolt.Tx, value T, data []byte) error {
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
		var old T
		if err := bson.Unmarshal(suffix(v), &old); err != nil {
			return err
		}
		if err := bucket.delIndex(tx, old); err != nil {
			return err
		}
	}
	bucket.debug("Put(%s, %s) = %q", value.Id(), ptr, data)
	if err := bucket.values(tx).Put(ptr, data); err != nil {
		return err
	}
	return bucket.putIndex(tx, value)
}

func (bucket *indexedBucket[T]) putIndex(tx *bbolt.Tx, value T) error {
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
		bucket.debug("putIndex(%s) to (%s) = %q", key, last, ptr)
	}
	return nil
}

func (bucket *indexedBucket[T]) delIndex(tx *bbolt.Tx, value T) error {
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

func (bucket *indexedBucket[T]) Next(k Key, set ...int) (int, error) {
	var i uint64
	elems, last := k.B()
	if len(last) > 0 {
		elems = append(elems, last)
	}
	err := bucket.db.Update(func(tx *bbolt.Tx) error {
		// seqBucket is a minimal interface for sequence operations used by Next().
		// This is necessary because we could be operating on an index bucket
		// or the values bucket depending on the key.
		type seqBucket interface {
			SetSequence(uint64) error
			NextSequence() (uint64, error)
		}
		var b seqBucket = bucket.values(tx)
		if len(elems) > 0 {
			b = bucket.find(tx, elems)
		}
		if b == nil {
			return bbolt.ErrBucketNotFound
		}

		// The empty key will increment the counter for the values
		// bucket, non-empty keys will be in the index buckets.
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
