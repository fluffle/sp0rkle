package stats

import (
	"fmt"
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

type NickStats []*NickStat

type Collection struct {
	db.C
}

func Init() *Collection {
	sc := &Collection{}
	sc.Init(db.Bolt.Indexed(), COLLECTION, nil)
	return sc
}

func (sc *Collection) StatsFor(nick, ch string) *NickStat {
	res := NewStat(bot.Nick(nick), bot.Chan(ch))
	if err := sc.Get(res.byKey(), res); err == nil {
		return res
	}
	return nil
}

func (sc *Collection) TopTen(ch string) []*NickStat {
	var bRes NickStats
	if err := sc.All(db.K{db.S{"lines", ch}}, &bRes); err != nil {
		return nil
	}
	// Results from bolt are in ascending order.
	for i, j := 0, len(bRes)-1; i < j; i, j = i+1, j-1 {
		bRes[i], bRes[j] = bRes[j], bRes[i]
	}
	if len(bRes) > 10 {
		bRes = bRes[:10]
	}
	return bRes
}
