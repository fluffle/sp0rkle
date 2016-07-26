package db

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fluffle/golog/logging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	prefixLen = 4
)

var (
	// It is important that these are prefixLen long.
	// idPrefix is conveniently what K{{"_id", "stuff"}} serializes to.
	idPrefix   = append([]byte("_id"), USEP)
	bsonPrefix = append([]byte("_bs"), USEP)
)

func isBson(data []byte) bool {
	if len(data) < prefixLen {
		return false
	}
	return bytes.Equal(data[:prefixLen], bsonPrefix)
}

func toBson(value interface{}) ([]byte, error) {
	marshalled, err := bson.Marshal(value)
	if err != nil {
		return nil, err
	}
	data := bytes.NewBuffer(make([]byte, 0, prefixLen+len(marshalled)))
	data.Write(bsonPrefix)
	data.Write(marshalled)
	return data.Bytes(), nil
}

func isPointer(data []byte) bool {
	if len(data) < prefixLen {
		return false
	}
	return bytes.Equal(data[:prefixLen], idPrefix)
}

func toPointer(value Indexer) []byte {
	e := &Elem{string(idPrefix), string(value.Id())}
	return e.Bytes()
}

func fromPointer(data []byte) K {
	return K{{"_id", string(data[prefixLen:])}}
}

func suffix(data []byte) []byte {
	if len(data) < prefixLen {
		return nil
	}
	return data[prefixLen:]
}

func resolvePointer(data []byte, root *bolt.Bucket) ([]byte, error) {
	switch {
	case data == nil:
		// Key not found (or is nested bucket).
		return nil, nil
	case isBson(data):
		return suffix(data), nil
	case isPointer(data):
		// Follow the pointer; all _id keys are in the root bucket.
		// No need to use suffix here, the pointer data is the key.
		data = root.Get(data)
		if isBson(data) {
			return suffix(data), nil
		}
		return nil, fmt.Errorf("resolvePointer: not bson: %q", data)
	default:
		return nil, fmt.Errorf("resolvePointer: unknown prefix in data %q", data)
	}
}

func bucketFor(key Key, b *bolt.Bucket) (*bolt.Bucket, []byte, error) {
	var err error
	elems, last := key.B()
	for _, e := range elems {
		// CreateBucketIfNotExists requires a writeable transaction.
		if new := b.Bucket(e); new != nil {
			b = new
			continue
		}
		if b, err = b.CreateBucket(e); err != nil {
			return b, last, err
		}
	}
	return b, last, nil
}

type boltDatabase struct {
	sync.Mutex
	db    *bolt.DB
	dir   string
	every time.Duration
	quit  chan struct{}
}

var Bolt = &boltDatabase{}

func (b *boltDatabase) Init(path, backupDir string, backupEvery time.Duration) error {
	b.Lock()
	defer b.Unlock()
	if b.db != nil {
		return errors.New("init already called")
	}
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return err
	}
	b.db, b.dir, b.every, b.quit = db, backupDir, backupEvery, make(chan struct{})
	// Do a backup on startup and error if it is not successful.
	if err := os.MkdirAll(b.dir, 0700); err != nil {
		return fmt.Errorf("could not create backup dir %q: %v", b.dir, err)
	}
	if err := b.doBackup(); err != nil {
		return fmt.Errorf("could not perform initial backup: %v", err)
	}
	go b.backupLoop()
	return nil
}

func (b *boltDatabase) Close() {
	b.Lock()
	defer b.Unlock()

	if b.db == nil {
		return
	}
	if err := b.db.Close(); err != nil {
		logging.Error("Unable to close BoltDB: %v", err)
	}
	b.db = nil
	close(b.quit)
}

func (b *boltDatabase) C(name string) Collection {
	b.Lock()
	defer b.Unlock()

	if b.db == nil {
		logging.Fatal("Tried to create BoltDB bucket %q when disconnected.", name)
	}

	err := b.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logging.Fatal("Creating BoltDB bucket failed: %v")
	}
	return &boltBucket{name: []byte(name), db: b.db}
}

func (b *boltDatabase) backupLoop() {
	tick := time.NewTicker(b.every)
	for {
		select {
		case <-tick.C:
			if err := b.doBackup(); err != nil {
				logging.Error("Backup error: %v", err)
			}
		case <-b.quit:
			tick.Stop()
			return
		}
	}
}

func (b *boltDatabase) doBackup() error {
	fn := path.Join(b.dir, fmt.Sprintf("sp0rkle.boltdb.%s.gz",
		time.Now().Format("2006-01-02.15:04")))
	fh, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("could not create %q: %v", fn, err)
	}
	fz := gzip.NewWriter(fh)
	defer fz.Close()
	err = b.db.View(func(tx *bolt.Tx) error {
		return tx.Copy(fz)
	})
	if err != nil {
		os.Remove(fn)
		return fmt.Errorf("could not copy db to %q: %v", fn, err)
	}
	logging.Info("Wrote backup to %q.", fn)
	return nil
}

type boltBucket struct {
	sync.Mutex
	name  []byte
	db    *bolt.DB
	debug bool
}

func (bucket *boltBucket) Debug(on bool) {
	bucket.debug = on
}

func (bucket *boltBucket) Get(key Key, value interface{}) error {
	return bucket.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket(bucket.name)
		b, k, err := bucketFor(key, root)
		if err != nil {
			return err
		}
		if len(k) == 0 {
			return errors.New("get: zero length key")
		}
		data, err := resolvePointer(b.Get(k), root)
		if bucket.debug {
			logging.Debug("Get(%s): %s = %q, %v", bucket.name, key, data, err)
		}
		if data == nil {
			return err
		}
		return bson.Unmarshal(data, value)
	})
}

// TODO(fluffle): Dedupe this with Prefix when less hungover.
func (bucket *boltBucket) All(key Key, value interface{}) error {
	// This entirely stolen from mgo's Iter.All() \o/
	// vv == value Value
	vv := reflect.ValueOf(value)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Slice {
		panic("All() requires a pointer-to-slice.")
	}
	// sv == slice Value
	sv := vv.Elem()
	// Resize slice to capacity.
	sv = sv.Slice(0, sv.Cap())
	// et == (slice) element Type
	et := sv.Type().Elem()

	// Depending on the key passed, we may end up finding multiple pointers to
	// the data that we care about, instead of or in addition to that data.
	// Returning dupes would be unhelpful.
	set := map[string]bool{}
	seen := func(v []byte) bool {
		_, ok := set[string(v)]
		set[string(v)] = true
		return ok
	}

	return bucket.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket(bucket.name)
		b, last, err := bucketFor(key, root)
		if err != nil {
			return err
		}
		// All implies that the last key elem is also a bucket,
		// but we support a zero-length key to perform a scan
		// over the root bucket.
		cs := []*bolt.Cursor{b.Cursor()}
		if len(last) > 0 {
			if b = b.Bucket(last); b == nil {
				return bolt.ErrBucketNotFound
			}
			cs[0] = b.Cursor()
		}
		var i int
		var c *bolt.Cursor
		for len(cs) > 0 {
			c, cs = cs[0], cs[1:]
			for k, v := c.First(); k != nil; k, v = c.Next() {
				switch {
				case v == nil:
					// All flattens the nested buckets under key.
					if nest := b.Bucket(k); nest != nil {
						cs = append(cs, nest.Cursor())
					}
					continue
				case isPointer(v) && seen(v):
					continue
				case isPointer(k) && seen(k):
					continue
				case isPointer(v):
					if v, err = resolvePointer(v, root); err != nil {
						return err
					} else if v == nil {
						continue
					}
				case isBson(v):
					v = suffix(v)
				default:
					// Reasonably sure we shouldn't hit this condition.
					logging.Warn("all: unexpected data k=%q v=%q", k, v)
					continue
				}

				if sv.Len() == i {
					// Extend sv to hold more elements.
					ev := reflect.New(et)
					sv = reflect.Append(sv, ev.Elem())
					sv = sv.Slice(0, sv.Cap())
				}
				if err := bson.Unmarshal(v, sv.Index(i).Addr().Interface()); err != nil {
					return err
				}
				i++
			}
		}
		if bucket.debug {
			logging.Debug("All(%s): %s found %d items.", bucket.name, key, i)
		}
		vv.Elem().Set(sv.Slice(0, i))
		return nil
	})
}

func (bucket *boltBucket) Prefix(key Key, value interface{}) error {
	// This entirely stolen from mgo's Iter.All() \o/
	// vv == value Value
	vv := reflect.ValueOf(value)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Slice {
		panic("Prefix() requires a pointer-to-slice.")
	}
	// sv == slice Value
	sv := vv.Elem()
	// Resize slice to capacity.
	sv = sv.Slice(0, sv.Cap())
	// et == (slice) element Type
	et := sv.Type().Elem()

	return bucket.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket(bucket.name)
		b, prefix, err := bucketFor(key, root)
		if err != nil {
			return err
		}
		if len(prefix) == 0 {
			logging.Warn("zero-length prefix scan for key %s.", key)
		}
		i := 0
		c := b.Cursor()
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			switch {
			case v == nil:
				// Prefix ignores nested buckets.
				continue
			case isPointer(v):
				if v, err = resolvePointer(v, root); err != nil {
					return err
				} else if v == nil {
					continue
				}
			case isBson(v):
				// TODO(fluffle): There is a corner case here where a zero
				// length prefix scan over the root bucket where both indexes
				// and data are stored will return duplicates.
				// Not sure I need to care right now.
				v = suffix(v)
			default:
				// Reasonably sure we shouldn't hit this condition.
				logging.Warn("prefix: unexpected data k=%q v=%q", k, v)
				continue
			}

			if sv.Len() == i {
				// Extend sv to hold more elements.
				ev := reflect.New(et)
				sv = reflect.Append(sv, ev.Elem())
				sv = sv.Slice(0, sv.Cap())
			}
			if err := bson.Unmarshal(v, sv.Index(i).Addr().Interface()); err != nil {
				return err
			}
			i++
		}
		if bucket.debug {
			logging.Debug("Prefix(%s): %s found %d items.", bucket.name, key, i)
		}
		vv.Elem().Set(sv.Slice(0, i))
		return nil
	})
}

func (bucket *boltBucket) Put(value interface{}) error {
	data, err := toBson(value)
	if err != nil {
		return err
	}
	switch value := value.(type) {
	case Keyer:
		return bucket.putKeyer(value, data)
	case Indexer:
		return bucket.putIndexer(value, data)
	}
	return fmt.Errorf("put: don't know how to put value %#v", value)
}

func (bucket *boltBucket) putKeyer(value Keyer, data []byte) error {
	return bucket.db.Update(func(tx *bolt.Tx) error {
		b, k, err := bucketFor(value.K(), tx.Bucket(bucket.name))
		if err != nil {
			return err
		}
		if len(k) == 0 {
			return errors.New("put: zero length key")
		}
		if bucket.debug {
			logging.Debug("Put(%s): %s = %q", bucket.name, value.K(), data)
		}
		return b.Put(k, data)
	})
}

func (bucket *boltBucket) putIndexer(value Indexer, data []byte) error {
	return bucket.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket(bucket.name)
		ptr := toPointer(value)
		v := root.Get(ptr)
		if isBson(v) {
			// There's already a value here, probably being pointed at.
			// Jump through some hoops to clean up those index pointers.
			// TODO(fluffle): This makes some assumptions that may not
			// hold true, and might leave dangling index pointers, ugh.
			//   1) The old value is of the same type as the new one.
			old := dupe(value).(Indexer)
			if err := bson.Unmarshal(suffix(v), old); err != nil {
				return err
			}
			//   2) The indexes derived from the old data are exactly
			//      the correct set that should be deleted to tidy up.
			for _, key := range old.Indexes() {
				b, k, err := bucketFor(key, root)
				if err != nil {
					return err
				}
				if bucket.debug {
					logging.Debug("Clean index(%s): %s = %q", bucket.name, key, ptr)
				}
				if err = b.Delete(k); err != nil {
					return err
				}
			}
		}
		if bucket.debug {
			logging.Debug("Put(%s): %s = %q", bucket.name, value.Id(), data)
		}
		if err := root.Put(ptr, data); err != nil {
			return err
		}
		for _, key := range value.Indexes() {
			b, k, err := bucketFor(key, root)
			if err != nil {
				return err
			}
			if bucket.debug {
				logging.Debug("Put index(%s): %s = %q", bucket.name, key, ptr)
			}
			if err = b.Put(k, ptr); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bucket *boltBucket) Del(value interface{}) error {
	switch value := value.(type) {
	case Keyer:
		return bucket.db.Update(func(tx *bolt.Tx) error {
			b, k, err := bucketFor(value.K(), tx.Bucket(bucket.name))
			if err != nil {
				return err
			}
			return b.Delete(k)
		})
	case Indexer:
		return bucket.db.Update(func(tx *bolt.Tx) error {
			root := tx.Bucket(bucket.name)
			if err := root.Delete(toPointer(value)); err != nil {
				return err
			}
			for _, key := range value.Indexes() {
				b, k, err := bucketFor(key, root)
				if err != nil {
					return err
				}
				if err = b.Delete(k); err != nil {
					return err
				}
			}
			return nil
		})
	}
	return fmt.Errorf("del: don't know how to delete value %#v", value)
}

func (bucket *boltBucket) Mongo() *mgo.Collection {
	panic("you are bad at migrations")
}
