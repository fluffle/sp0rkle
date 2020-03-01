package db

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/fluffle/golog/logging"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/mgo.v2/bson"
)

const (
	prefixLen = 4
	idTag     = "_id"
	bsonTag   = "_bs"
)

var (
	// It is important that these are prefixLen long.
	// idPrefix is conveniently what K{{"_id", "stuff"}} serializes to.
	idPrefix   = append([]byte(idTag), USEP)
	bsonPrefix = append([]byte(bsonTag), USEP)
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
		return nil, fmt.Errorf("bson marshal: %w", err)
	}
	data := bytes.NewBuffer(make([]byte, 0, prefixLen+len(marshalled)))
	data.Write(bsonPrefix)
	data.Write(marshalled)
	return data.Bytes(), nil
}

func suffix(data []byte) []byte {
	if len(data) < prefixLen {
		return nil
	}
	return data[prefixLen:]
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

func (b *boltDatabase) DB() *bolt.DB {
	return b.db
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
