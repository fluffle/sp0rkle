package db

// Wraps an mgo connection and db object for convenience
// Yes, these are globals. I'm undecided, but let's see how it goes.

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/fluffle/golog/logging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const DATABASE string = "sp0rkle"

type mongoDatabase struct {
	sync.Mutex
	sessions []*mgo.Session
}

var Mongo = &mongoDatabase{}

func (m *mongoDatabase) Init(db string) error {
	m.Lock()
	defer m.Unlock()
	if m.sessions != nil {
		return errors.New("init already called")
	}
	s, err := mgo.Dial(db)
	if err != nil {
		return err
	}
	// Let's be explicit about requiring journaling, ehh?
	s.EnsureSafe(&mgo.Safe{J: true})
	m.sessions = []*mgo.Session{s}
	return nil
}

func (m *mongoDatabase) Close() {
	m.Lock()
	defer m.Unlock()
	if m.sessions == nil {
		return
	}
	for _, s := range m.sessions {
		s.Close()
	}
	m.sessions = nil
}

func (m *mongoDatabase) C(name string) Collection {
	m.Lock()
	defer m.Unlock()
	if m.sessions == nil {
		logging.Fatal("Tried to create MongoDB collection %q when disconnected.", name)
	}
	s := m.sessions[0].Copy()
	m.sessions = append(m.sessions, s)
	return &mongoCollection{Collection: s.DB(DATABASE).C(name)}
}

type mongoCollection struct {
	*mgo.Collection
	debug bool
}

func (m *mongoCollection) Debug(on bool) {
	m.debug = on
}

func (m *mongoCollection) Get(key Key, value interface{}) error {
	return m.Collection.Find(key.M()).One(value)
}

func (m *mongoCollection) Match(key, regex string, value interface{}) error {
	q := bson.M{strings.ToLower(key): bson.M{"$regex": regex, "$options": "i"}}
	return m.Collection.Find(q).All(value)
}

func (m *mongoCollection) All(key Key, value interface{}) error {
	return m.Collection.Find(key.M()).All(value)
}

func (m *mongoCollection) Put(value interface{}) (err error) {
	switch value := value.(type) {
	case Keyer:
		_, err = m.Collection.Upsert(value.K().M(), value)
	case Indexer:
		_, err = m.Collection.UpsertId(value.Id(), value)
	default:
		return fmt.Errorf("put: don't know how to put value %#v", value)
	}
	return err
}

func (m *mongoCollection) BatchPut(value interface{}) error {
	panic("no batch puts for you")
}

func (m *mongoCollection) Del(value interface{}) error {
	switch value := value.(type) {
	case Keyer:
		return m.Collection.Remove(value.K().M())
	case Indexer:
		return m.Collection.RemoveId(string(value.Id()))
	}
	return fmt.Errorf("del: don't know how to delete value %#v", value)
}

func (m *mongoCollection) Next(k Key, set ...uint64) (int, error) {
	panic("no autoincrements for you")
}

func (m *mongoCollection) Mongo() *mgo.Collection {
	return m.Collection
}
