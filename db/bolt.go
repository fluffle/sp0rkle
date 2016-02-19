package db

import (
	"errors"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fluffle/golog/logging"
)

type boltDatabase struct {
	sync.Mutex
	db *bolt.DB
}

var Bolt Database = &boltDatabase{}

type boltBucket struct {
	sync.Mutex
	name []byte
}

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
	return &boltBucket{name: []byte(name)}
}
