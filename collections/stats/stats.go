package stats

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION string = "stats"

type NickStat struct {
	Nick   bot.Nick
	Key    string
	Chan   bot.Chan
	Lines  int
	Words  int
	Chars  int
	Active [7][24]int
	Id_    bson.ObjectId `bson:"_id,omitempty"`
}

var _ db.Indexer = (*NickStat)(nil)

func NewStat(n bot.Nick, c bot.Chan) *NickStat {
	return &NickStat{
		Nick:   n,
		Key:    strings.ToLower(string(n)),
		Chan:   c,
		Active: [7][24]int{},
		Id_:    bson.NewObjectId(),
	}
}

func (ns *NickStat) Update(line string) {
	ns.Lines++
	ns.Words += len(strings.Fields(line))
	ns.Chars += len(line)
	t := time.Now()
	ns.Active[int(t.Weekday())][t.Hour()]++
}

func (ns *NickStat) MostActive() (day time.Weekday, hour int, count int) {
	for d, times := range ns.Active {
		for h, c := range times {
			if c > count {
				day = time.Weekday(d)
				hour = h
				count = c
			}
		}
	}
	return
}

func (ns *NickStat) String() string {
	day, hour, count := ns.MostActive()
	wordc := float64(ns.Words) / float64(ns.Lines)
	charc := float64(ns.Chars) / float64(ns.Lines)
	return fmt.Sprintf("%s has said %d words and %d lines in %s. "+
		"Each line averaged %.2f words and %.2f chars. "+
		"They are most active on %ss at around %d:00, "+
		"saying %d lines in that hour.",
		ns.Nick, ns.Words, ns.Lines, ns.Chan,
		wordc, charc, day, hour, count)
}

func (ns *NickStat) Indexes() []db.Key {
	return []db.Key{
		db.K{db.S{"chan", string(ns.Chan)}, db.S{"key", ns.Key}},
		db.K{db.S{"lines", string(ns.Chan)}, db.I{"lines", ns.Lines}},
	}
}

func (ns *NickStat) Id() bson.ObjectId {
	return ns.Id_
}

func (ns *NickStat) byKey() db.Key {
	return db.K{db.S{"chan", string(ns.Chan)}, db.S{"key", ns.Key}}
}

type NickStats []*NickStat

func (ns NickStats) Strings() []string {
	s := make([]string, len(ns))
	for i, n := range ns {
		// Can't use String() here since it doesn't
		// contain all the relevant info
		s[i] = fmt.Sprintf("%#v", n)
	}
	return s
}

type migrator struct {
	mongo, bolt db.Collection
}

func (m *migrator) Migrate() error {
	var all []*NickStat
	if err := m.mongo.All(db.K{}, &all); err != nil {
		return err
	}
	if err := m.bolt.BatchPut(all); err != nil {
		logging.Error("Migrating stats entries: %v", err)
		return err
	}
	logging.Info("Migrated %d stats entries.", len(all))
	return nil
}

func (m *migrator) Diff() ([]string, []string, error) {
	var mAll, bAll NickStats
	if err := m.mongo.All(db.K{}, &mAll); err != nil {
		return nil, nil, err
	}
	if err := m.bolt.All(db.K{}, &bAll); err != nil {
		return nil, nil, err
	}
	return mAll.Strings(), bAll.Strings(), nil
}

type Collection struct {
	db.Both
}

func Init() *Collection {
	sc := &Collection{db.Both{}}
	sc.Both.MongoC.Init(db.Mongo, COLLECTION, mongoIndexes)
	sc.Both.BoltC.Init(db.Bolt, COLLECTION, nil)
	m := &migrator{
		mongo: sc.Both.MongoC,
		bolt:  sc.Both.BoltC,
	}
	sc.Both.Checker.Init(m, COLLECTION)
	return sc
}

func mongoIndexes(c db.Collection) {
	indexes := [][]string{
		{"chan", "key"},
		{"lines"},
	}
	for _, key := range indexes {
		if err := c.Mongo().EnsureIndex(mgo.Index{Key: key}); err != nil {
			logging.Error("Couldn't create %v index on sp0rkle.stats: %v", key, err)
		}
	}
}

func (sc *Collection) StatsFor(nick, ch string) *NickStat {
	res := NewStat(bot.Nick(nick), bot.Chan(ch))
	if err := sc.Get(res.byKey(), res); err == nil {
		return res
	}
	return nil
}

func (sc *Collection) TopTen(ch string) []*NickStat {
	var mRes, bRes []*NickStat
	q := sc.Mongo().Find(bson.M{"chan": ch}).Sort("-lines").Limit(10)
	if err := q.All(&mRes); err != nil {
		logging.Error("Mongo TopTen Find error for channel %s: %v", ch, err)
	}
	if err := sc.Both.BoltC.All(db.K{db.S{"lines", ch}}, &bRes); err != nil {
		logging.Error("Bolt TopTen All error for channel %s: %v", ch, err)
	}
	// TODO(fluffle): Results from Bolt are in ascending order, meh.
	// TODO(fluffle): Consider supporting asc/desc/limit in db.C interface.
	for i, j := 0, len(bRes)-1; i < j; i, j = i+1, j-1 {
		bRes[i], bRes[j] = bRes[j], bRes[i]
	}
	if len(bRes) > 10 {
		bRes = bRes[:10]
	}
	if !reflect.DeepEqual(mRes, bRes) {
		logging.Warn("TopTen mismatch for channel %s.", ch)
		for i, v := range mRes {
			logging.Debug("Mongo %d: %#v", i, v)
		}
		for i, v := range bRes {
			logging.Debug("Bolt %d: %#v", i, v)
		}
	}
	if sc.Migrated() {
		return bRes
	}
	return mRes
}
