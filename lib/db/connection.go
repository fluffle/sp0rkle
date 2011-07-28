package db

// Wraps an mgo connection and db object for convenience

import (
	"launchpad.net/mgo"
	"os"
)

const DATABASE string = "sp0rkle"

type Database struct {
	// We're wrapping mgo.Database here so we can provide our own methods.
	mgo.Database

	// But unlike mgo.Database, it'd be useful to keep an internal session
	// reference around, so we can close things out later.
	Session *mgo.Session
}

// Wraps connecting to mongo and selecting the "sp0rkle" database.
func Connect(resource string) (*Database, os.Error) {
	sess, err := mgo.Mongo(resource)
	if err != nil {
		return nil, err
	}
	return &Database{Database: sess.DB(DATABASE), Session: sess}, nil
}
