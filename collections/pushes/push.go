package pushes

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/datetime"
	"golang.org/x/oauth2"
	"github.com/fluffle/sp0rkle/util/bson"
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
	Id_     bson.ObjectId `bson:"_id,omitempty"`
}

var _ db.Indexer = (*State)(nil)

func (s *State) String() string {
	return fmt.Sprintf("Push for %q (%d aliases); done=%t at %s; iden=%q pin=%q tok=%q",
		s.Nick, len(s.Aliases), s.Done, datetime.Format(s.Time), s.Iden, s.Pin, s.Token)
}

func (s *State) Id() bson.ObjectId {
	return s.Id_
}

func (s *State) Exists() bool {
	return s != nil && len(s.Id_) > 0
}

func (s *State) Indexes() []db.Key {
	k := []db.Key{db.K{db.S{"nick", s.Nick}}}
	for _, alias := range s.Aliases {
		k = append(k, db.K{db.S{"aliases", alias}})
	}
	return k
}

func (s *State) byId() db.K {
	return db.K{db.ID{s.Id_}}
}

func byNick(nick string) db.K {
	return db.K{db.S{"nick", nick}}
}

func byAlias(alias string) db.K {
	return db.K{db.S{"aliases", alias}}
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
	return base64.URLEncoding.EncodeToString([]byte(s.Id_))
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
	db.C
}

func Init() *Collection {
	pc := &Collection{}
	pc.Init(db.Bolt.Indexed(), COLLECTION, nil)
	if err := pc.Fsck(&State{}); err != nil {
		logging.Fatal("pushes fsck failed: %v", err)
	}
	return pc
}

func (pc *Collection) NewState(nick string) (*State, error) {
	s := &State{
		Nick: strings.ToLower(nick),
		Time: time.Now(),
		Id_:  bson.NewObjectId(),
	}
	if err := pc.Put(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (pc *Collection) GetByB64(b64 string) *State {
	id, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		logging.Error("Decoding base64 string %q: %v", b64, err)
		return nil
	}
	s := &State{Id_: bson.ObjectId(id)}
	if err := pc.Get(s.byId(), s); err != nil {
		logging.Error("Looking up state with id=%q: %v", id, err)
		return nil
	}
	if s.AuthWindowExpired() {
		if err := pc.Del(s); err != nil {
			logging.Error("Deleting state with id=%q: %v", id, err)
		}
		return nil
	}
	return s
}

func (pc *Collection) GetByNick(nick string, checkAliases bool) *State {
	s := &State{}
	if err := pc.Get(byNick(nick), s); err != nil {
		logging.Error("Looking up state with nick=%q: %v", nick, err)
		return nil
	}
	if !s.Exists() && checkAliases {
		// Not found by nick, check aliases.
		if err := pc.Get(byAlias(nick), s); err != nil {
			logging.Error("Looking up state with alias=%q: %v", nick, err)
			return nil
		}
	}
	if !s.Exists() {
		return nil
	}
	if s.AuthWindowExpired() {
		if err := pc.Del(s); err != nil {
			logging.Error("Deleting state with id=%q: %v", s.Id_, err)
		}
		return nil
	}
	return s
}
