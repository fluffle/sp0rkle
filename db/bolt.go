package db

import (
	"bytes"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fluffle/golog/logging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type boltDatabase struct {
	sync.Mutex
	db *bolt.DB
}

var Bolt Database = &boltDatabase{}

func (b *boltDatabase) Init(path string) error {
	b.Lock()
	defer b.Unlock()
	if b.db != nil {
		return errors.New("init already called")
	}
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return err
	}
	b.db = db
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

type boltBucket struct {
	sync.Mutex
	name []byte
	db   *bolt.DB
}

func (bucket *boltBucket) Get(key Key, value interface{}) error {
	return bucket.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.name)
		data := b.Get(key.B())
		if data == nil || len(data) == 0 {
			return nil
		}
		return bson.Unmarshal(data, value)
	})
}

func (bucket *boltBucket) All(key Key, value interface{}) error {
	// This entirely stolen from mgo's Iter.All() \o/
	vv := reflect.ValueOf(value)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Slice {
		panic("All() requires a pointer-to-slice.")
	}
	sv := vv.Elem()
	sv = sv.Slice(0, sv.Cap())
	et := sv.Type().Elem()
	prefix := key.B()

	return bucket.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucket.name).Cursor()
		i := 0
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			if sv.Len() == i {
				ev := reflect.New(et)
				if err := bson.Unmarshal(v, ev.Interface()); err != nil {
					return err
				}
				sv = reflect.Append(sv, ev.Elem())
				sv = sv.Slice(0, sv.Cap())
			} else {
				if err := bson.Unmarshal(v, sv.Index(i).Addr().Interface()); err != nil {
					return err
				}
			}
			i++
		}
		vv.Elem().Set(sv.Slice(0, i))
		return nil
	})
}

func (bucket *boltBucket) Put(key Key, value interface{}) error {
	return bucket.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.name)
		data, err := bson.Marshal(value)
		if err != nil {
			return err
		}
		return b.Put(key.B(), data)
	})
}

func (bucket *boltBucket) Del(key Key) error {
	return bucket.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket.name).Delete(key.B())
	})
}

func (bucket *boltBucket) Mongo() *mgo.Collection {
	panic("you are bad at migrations")
}
