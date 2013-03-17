package stats

import (
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
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
}

func NewStat(n bot.Nick, c bot.Chan) *NickStat {
	return &NickStat{
		Nick:   n,
		Key:    strings.ToLower(string(n)),
		Chan:   c,
		Active: [7][24]int{},
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

func (ns *NickStat) Id() bson.M {
	return bson.M{"nick": ns.Nick, "chan": ns.Chan}
}

type Collection struct {
	*mgo.Collection
}

func Init() *Collection {
	sc := &Collection{db.Init().C(COLLECTION)}
	indexes := [][]string{
		{"chan", "key"},
		{"lines"},
	}
	for _, key := range indexes {
		if err := sc.EnsureIndex(mgo.Index{Key: key}); err != nil {
			logging.Error("Couldn't create %v index on sp0rkle.stats: %v", key, err)
		}
	}
	return sc
}

func (sc *Collection) StatsFor(nick, ch string) *NickStat {
	var res NickStat
	q := sc.Find(bson.M{
		"chan": ch,
		"key":  strings.ToLower(nick),
	})
	if err := q.One(&res); err == nil {
		return &res
	}
	return nil
}

func (sc *Collection) TopTen(ch string) []*NickStat {
	var res []*NickStat
	q := sc.Find(bson.M{"chan": ch}).Sort("-lines").Limit(10)
	if err := q.All(&res); err != nil {
		logging.Error("TopTen Find error for channel %s: %v", ch, err)
	}
	return res
}
