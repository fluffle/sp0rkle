package factoids

// This might get ODM-ish in the future.

import (
	"github.com/garyburd/go-mongo"
	"log"
	"os"
	"time"
)

type FactoidType int

const (
	// Factoids can be of these types
	F_FACT FactoidType = iota
	F_ACTION
	F_REPLY
	F_URL
)

// A factoid maps a key to a value, and keeps some stats about it
type Factoid struct {
	Key, Value                  string
	Type                        FactoidType
	Created, Modified, Accessed *FactoidStat
	Perms                       *FactoidPerms
	Id                          mongo.ObjectId `bson:"_id"`
}

// Represent info about things that happened to the factoid
type FactoidStat struct {
	// When <thing> happened
	Timestamp *time.Time
	// Who did <thing>
	Nick, Ident, Host string
	// Where they did <thing>
	Chan string
	// How many times <thing> has been done before
	Count int
}

// Represent info about things that can be done to the factoid
type FactoidPerms struct {
	ReadOnly bool
	Owner    string
}

// Factoids are stored in a mongo collection of Factoid structs
type FactoidCollection struct {
	*mongo.Collection
}

// Wrapper to get hold of a factoid collection handle
func Collection(conn mongo.Conn) (*FactoidCollection, os.Error) {
	fc := &FactoidCollection{
		&mongo.Collection{
			Conn:         conn,
			Namespace:    "sp0rkle.factoids",
			LastErrorCmd: mongo.DefaultLastErrorCmd,
		},
	}
	err := fc.CreateIndex(mongo.D{mongo.DocItem{Key: "Key", Value: 1}}, nil)
	if err != nil {
		log.Printf("Couldn't create index on sp0rkle.factoids: %v", err)
		return nil, err
	}
	return fc, nil
}

func (fc FactoidCollection) GetOne(key string) (res *Factoid) {
	fc.Find(mongo.M{"Key": key}).One(res)
	return
}
