package pushes

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/db"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION = "push"

type State struct {
	Nick    string        `json:"nick"`
	Aliases []string      `json:"aliases,omitempty"`
	Iden    string        `json:"iden,omitempty"`
	Pin     string        `json:"pin"`
	Token   *oauth2.Token `json:"token,omitempty"`
	Done    bool          `json:"done"`
	Time    time.Time     `json:"time"`
	Id      bson.ObjectId `bson:"_id,omitempty"`
}

func (s *State) AuthWindowExpired() bool {
	// We have an hour's grace time to complete the auth flow.
	return s == nil || (!s.CanPush() &&
		time.Now().After(s.Time.Add(time.Hour)))
}

func (s *State) CanConfirm() bool {
	return s != nil && s.Token != nil && s.Iden != "" && !s.Done
}

func (s *State) CanPush() bool {
	return s != nil && s.Token != nil && s.Iden != "" && s.Done
}

func (s *State) State() string {
	return base64.URLEncoding.EncodeToString([]byte(s.Id))
}

func (s *State) HasAlias(alias string) bool {
	return s.aliasIndex(alias) != -1
}

func (s *State) AddAlias(alias string) {
	s.Aliases = append(s.Aliases, strings.ToLower(alias))
}

func (s *State) DelAlias(alias string) {
	idx := s.aliasIndex(alias)
	if idx == -1 {
		return
	}
	s.Aliases = append(s.Aliases[:idx], s.Aliases[idx+1:]...)
}

func (s *State) aliasIndex(alias string) int {
	lc := strings.ToLower(alias)
	for i, a := range s.Aliases {
		if a == lc {
			return i
		}
	}
	return -1
}

type Collection struct {
	*mgo.Collection
}

func Init() *Collection {
	pc := &Collection{db.Mongo.C(COLLECTION).Mongo()}
	if err := pc.EnsureIndex(mgo.Index{
		Key:    []string{"nick"},
		Unique: true,
	}); err != nil {
		logging.Error("Couldn't create an index on push: %s", err)
	}
	return pc
}

func (pc *Collection) NewState(nick string) (*State, error) {
	s := &State{
		Nick: strings.ToLower(nick),
		Time: time.Now(),
		Id:   bson.NewObjectId(),
	}
	if err := pc.Insert(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (pc *Collection) SetState(s *State) error {
	if _, err := pc.UpsertId(s.Id, s); err != nil {
		return err
	}
	return nil
}

func (pc *Collection) DelState(s *State) error {
	return pc.RemoveId(s.Id)
}

func (pc *Collection) GetByB64(b64 string) *State {
	id, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		return nil
	}
	s := &State{}
	if err := pc.FindId(bson.ObjectId(id)).One(s); err != nil {
		return nil
	}
	if s.AuthWindowExpired() {
		pc.RemoveId(s.Id)
		return nil
	}
	return s
}

func (pc *Collection) GetByNick(nick string, aliases ...bool) *State {
	s := &State{}
	query := bson.M{"$or": []bson.M{
		{"nick": strings.ToLower(nick)},
		{"aliases": strings.ToLower(nick)},
	}}
	if len(aliases) > 0 && !aliases[0] {
		// If aliases is explicitly false do not allow aliases to match when querying.
		query = bson.M{"nick": strings.ToLower(nick)}
	}
	if err := pc.Find(query).One(s); err != nil {
		return nil
	}
	if s.AuthWindowExpired() {
		pc.RemoveId(s.Id)
		return nil
	}
	return s
}
