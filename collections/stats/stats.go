package stats

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/bson"
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
		// TODO: This index causes fsck churn, because it's entirely possible
		// for many users in a channel to have said the same number of lines.
		// As the fsckValues iterator finds each one of them it repoints the
		// lines index for that line to that value, so the last-iterated
		// NickStat with a non-unique line count wins.
		db.K{db.S{"lines", string(ns.Chan)}, db.I{"lines", uint64(ns.Lines)}},
	}
}

func (ns *NickStat) Id() bson.ObjectId {
	return ns.Id_
}

func (ns *NickStat) byKey() db.Key {
	return db.K{db.S{"chan", string(ns.Chan)}, db.S{"key", ns.Key}}
}

func (ns *NickStat) CollectionName() string { return "stats" }

type NickStats []*NickStat

type Collection struct {
	collection db.IndexedCollection[*NickStat]
}

func New(collection db.IndexedCollection[*NickStat]) *Collection {
	return &Collection{collection: collection}
}

func (sc *Collection) StatsFor(nick, ch string) *NickStat {
	res := NewStat(bot.Nick(nick), bot.Chan(ch))
	stat, err := sc.collection.Get(res.byKey())
	if err != nil || stat == nil {
		return nil
	}
	return stat
}

func (sc *Collection) TopTen(ch string) []*NickStat {
	bRes, err := sc.collection.All(db.K{db.S{"lines", ch}})
	if err != nil {
		return nil
	}
	// Sort descending by Lines count
	sort.Slice(bRes, func(i, j int) bool {
		return bRes[i].Lines > bRes[j].Lines
	})
	if len(bRes) > 10 {
		bRes = bRes[:10]
	}
	return bRes
}

func (sc *Collection) Put(value *NickStat) error {
	return sc.collection.Put(value)
}

func (sc *Collection) Del(value *NickStat) error {
	return sc.collection.Del(value)
}
