package db

// Wraps an mgo connection and db object for convenience
// Yes, these are globals. I'm undecided, but let's see how it goes.

import (
	"errors"
	"sync"

	"github.com/fluffle/golog/logging"
	"gopkg.in/mgo.v2"
)

const DATABASE string = "sp0rkle"

type mongoDatabase struct {
	sync.Mutex
	sessions []*mgo.Session
}

var Mongo Database = &mongoDatabase{}

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
	return &mongoCollection{s.DB(DATABASE).C(name)}
}

type mongoCollection struct {
	*mgo.Collection
}

func (m *mongoCollection) Get(key Key, value interface{}) error {
	return m.Collection.Find(key.M()).One(value)
}

func (m *mongoCollection) All(key Key, value interface{}) error {
	return m.Collection.Find(key.M()).All(value)
}

func (m *mongoCollection) Put(key Key, value interface{}) error {
	_, err := m.Collection.Upsert(key.M(), value)
	return err
}

func (m *mongoCollection) Del(key Key) error {
	return m.Collection.Remove(key.M())
}

func (m *mongoCollection) Mongo() *mgo.Collection {
	return m.Collection
}
