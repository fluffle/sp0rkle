package markov

import (
	//	"github.com/fluffle/golog/logging"
	"fmt"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/markov"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const COLLECTION string = "markov"

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
		Tag:    tag,
		Id:     bson.NewObjectId(),
	}
}

// Factoids are stored in a mongo collection of Factoid structs
type Collection struct {
	// We're wrapping mgo.Collection so we can provide our own methods.
	*mgo.Collection
}

// Wrapper to get hold of a factoid collection handle
func Init() *Collection {
	fc := &Collection{
		Collection: db.Init().C(COLLECTION),
	}
	return fc
}

func (fc *Collection) Get(source, dest, tag string) (result *MarkovLink) {
	if q := fc.Find(bson.M{"source": source,
		"tag":  tag,
		"dest": dest}).One(&result); q == nil {
		return
	}
	return nil
}

func (mc *Collection) incUses(source, dest, tag string) error {
	link := mc.Get(source, dest, tag)
	if link == nil {
		link = New(source, dest, tag)
	}
	fmt.Printf("%v\n", link)
	link.Uses++

	if _, err := mc.UpsertId(link.Id, link); err != nil {
		return fmt.Errorf("Failed to insert MarkovLink %s->%s as %s: 5s",
			source, dest, tag, err)
	}
	return nil
}

func (mc *Collection) AddSentence(sentence []string, tag string) error {
	prev := markov.SENTENCE_START
	for i := 0; i < len(sentence); i++ {
		dest := sentence[i]
		mc.incUses(prev, dest, tag)
		prev = dest
	}
	mc.incUses(prev, markov.SENTENCE_END, tag)
	return nil
}

type CollectionSource struct {
	collection  *Collection
	tag         string
	isTagPrefix bool
}

func (cs *CollectionSource) getSearchExpression(value string) (ret bson.M) {
	ret = bson.M{"source": value}

	if cs.isTagPrefix {
		ret["tag"] = bson.RegEx{cs.tag + ".*", ""}
	} else {
		ret["tag"] = cs.tag
	}
	return
}

func (cs *CollectionSource) GetLinks(value string) ([]markov.Link, error) {
	search := cs.getSearchExpression(value)

	num, err := cs.collection.Find(search).Count()
	if err != nil {
		return nil, err
	}

	output := make([]markov.Link, 0, num)

	var result MarkovLink
	iter := cs.collection.Find(search).Iter()
	for iter.Next(&result) {
		output = append(output,
			markov.Link{result.Dest, result.Uses})
	}
	if iter.Err() != nil {
		return output, iter.Err()
	}

	return output, nil
}

func (fc *Collection) CreateSourceForTag(tag string) markov.Source {
	return &CollectionSource{
		collection:  fc,
		tag:         tag,
		isTagPrefix: false}
}

func (fc *Collection) CreateSourceForTagPrefix(tag string) markov.Source {
	return &CollectionSource{
		collection:  fc,
		tag:         tag,
		isTagPrefix: true}
}
