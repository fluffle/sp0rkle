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
var sessions []*mgo.Session

// Wraps connecting to mongo and selecting the "sp0rkle" database.
func Init() *mgo.Database {
	lock.Lock()
	defer lock.Unlock()
	if sessions != nil {
		// Give each caller a distinct session to avoid contention.
		s := sessions[0].Copy()
		sessions = append(sessions, s)
		return s.DB(DATABASE)
	}
	sessions = make([]*mgo.Session, 1)
	s, err := mgo.Dial(*database)
	if err != nil {
		logging.Fatal("Unable to connect to MongoDB: %s", err)
	}
	// Let's be explicit about requiring journaling, ehh?
	s.EnsureSafe(&mgo.Safe{J: true})
	sessions[0] = s
	return s.DB(DATABASE)
}

func Close() {
	lock.Lock()
	defer lock.Unlock()
	if sessions == nil {
		return
	}
	for _, s := range sessions {
		s.Close()
	}
	sessions = nil
}
