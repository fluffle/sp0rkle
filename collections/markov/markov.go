package markov

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/markov"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
)

const COLLECTION = "markov"

type MarkovLink struct {
	Source, Dest string
	Uses         int
	Tag          string
	Id           bson.ObjectId `bson:"_id,omitempty"`
}

func New(source, dest, tag string) *MarkovLink {
	return &MarkovLink{
		Source: source,
		Dest:   dest,
		Tag:    strings.ToLower(tag),
		Id:     bson.NewObjectId(),
	}
}

// Markov links are stored in a mongo collection
type Collection struct {
	*mgo.Collection
}

// Wrapper to get hold of a factoid collection handle
func Init() *Collection {
	mc := &Collection{db.Init().C(COLLECTION)}
	if err := mc.EnsureIndex(mgo.Index{
		Key: []string{"tag", "source", "dest"},
	}); err != nil {
		logging.Error("Couldn't create an index on markov: %s", err)
	}
	return mc
}

func (mc *Collection) Get(source, dest, tag string) *MarkovLink {
	var res MarkovLink
	if err := mc.Find(bson.M{
		"tag":    strings.ToLower(tag),
		"source": source,
		"dest":   dest,
	}).One(&res); err == nil {
		return &res
	}
	return nil
}

func (mc *Collection) incUses(source, dest, tag string) {
	link := mc.Get(source, dest, tag)
	if link == nil {
		link = New(source, dest, tag)
	}
	link.Uses++

	if _, err := mc.UpsertId(link.Id, link); err != nil {
		logging.Error("Failed to insert MarkovLink %s(%s->%s): %s",
			tag, source, dest, err)
	}
}

func (mc *Collection) AddAction(action, tag string) {
	mc.Add(markov.ACTION_START, action, tag)
}

func (mc *Collection) AddSentence(sentence, tag string) {
	mc.Add(markov.SENTENCE_START, sentence, tag)
}

func (mc *Collection) Add(source, data, tag string) {
	for _, dest := range strings.Fields(data) {
		mc.incUses(source, dest, tag)
		source = dest
	}
	mc.incUses(source, markov.SENTENCE_END, tag)
}

func (mc *Collection) ClearTag(tag string) error {
	if _, err := mc.RemoveAll(bson.M{"tag": tag}); err != nil {
		return err
	}
	return nil
}

type MarkovSource struct {
	*Collection
	tag string
}

func (mc *Collection) Source(tag string) *MarkovSource {
	return &MarkovSource{mc, tag}
}

func (ms *MarkovSource) GetLinks(value string) ([]markov.Link, error) {
	q := ms.Find(bson.M{
		"source": value,
		"tag":  ms.tag,
	})
	num, err := q.Count()
	if err != nil {
		return nil, err
	}

	output, iter := make([]markov.Link, 0, num), q.Iter()
	var result MarkovLink
	for iter.Next(&result) {
		output = append(output, markov.Link{result.Dest, result.Uses})
	}

	if iter.Err() != nil {
		return output, iter.Err()
	}
	return output, nil
}
