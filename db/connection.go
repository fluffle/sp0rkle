package db

// Wraps an mgo connection and db object for convenience
// Yes, these are globals. I'm undecided, but let's see how it goes.

import (
	"flag"
	"github.com/fluffle/golog/logging"
	"labix.org/v2/mgo"
	"sync"
)

const DATABASE string = "sp0rkle"

var database *string = flag.String("database", "localhost",
	"Address of MongoDB server to connect to, defaults to localhost.")

var lock sync.Mutex
var db *mgo.Database
var session *mgo.Session

// Wraps connecting to mongo and selecting the "sp0rkle" database.
func Init() *mgo.Database {
	lock.Lock()
	defer lock.Unlock()
	if db != nil {
		return db
	}
	s, err := mgo.Dial(*database)
	if err != nil {
		logging.Fatal("Unable to connect to MongoDB: %s", err)
	}
	session, db = s, s.DB(DATABASE)
	return db
}

func Close() {
	lock.Lock()
	defer lock.Unlock()
	if db == nil {
		return
	}
	session.Close()
	session, db = nil, nil
}
