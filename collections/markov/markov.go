package markov

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/util/diff"
	"github.com/fluffle/sp0rkle/util/markov"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION = "markov"

type MarkovLink struct {
	Source, Dest string
	Uses         int
	uses         []byte
	Tag          string
	Id_          bson.ObjectId `bson:"_id,omitempty"`
}

var _ db.Indexer = (*MarkovLink)(nil)

func New(source, dest, tag string) *MarkovLink {
	return &MarkovLink{
		Source: source,
		Dest:   dest,
		Tag:    strings.ToLower(tag),
		Id_:    bson.NewObjectId(),
	}
}

func (ml *MarkovLink) String() string {
	return fmt.Sprintf("%s(%q->%q):%d", ml.Tag, ml.Source, ml.Dest, ml.Uses)
}

func (ml *MarkovLink) Indexes() []db.Key {
	return []db.Key{
		db.K{db.S{"tag", ml.Tag}, db.S{"source", ml.Source}, db.S{"dest", ml.Dest}},
	}
}

func (ml *MarkovLink) Id() bson.ObjectId {
	return ml.Id_
}

func (ml *MarkovLink) byTagSrcDest() db.Key {
	return db.K{db.S{"tag", ml.Tag}, db.S{"source", ml.Source}, db.S{"dest", ml.Dest}}
}

func (ml *MarkovLink) byTagSrc() db.Key {
	return db.K{db.S{"tag", ml.Tag}, db.S{"source", ml.Source}}
}

type MarkovLinks []*MarkovLink

func (mls MarkovLinks) Strings() []string {
	s := make([]string, len(mls))
	for i, ml := range mls {
		s[i] = ml.String()
	}
	return s
}

type Collection struct {
	checker db.M
	mongo   db.C
	// Markov is a bit special. Because of the quantity of data
	// stored, it is desirable to be able to use the boltdb key
	// for storage too -- something that the standard Collection
	// interface does not provide for, in the possibly misguided
	// name of API simplicity. So instead of hacking this in, we
	// skip the entire db layer and deal with it here instead.
	bolt *bolt.DB
}

// Wrapper to get hold of a factoid collection handle
func Init() *Collection {
	mc := &Collection{}
	mc.mongo.Init(db.Mongo, COLLECTION, mongoIndexes)
	mc.bolt = db.Bolt.DB()
	mc.checker.Init(mc, COLLECTION)

	err := mc.bolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(COLLECTION))
		return err
	})
	if err != nil {
		logging.Fatal("Creating Markov BoltDB bucket failed: %v", err)
	}
	return mc
}

func (mc *Collection) Migrate() error {
	m := mc.mongo.Mongo()

	// Migrate each tag separately.
	tags := make([]string, 0)
	if err := m.Find(bson.M{}).Distinct("tag", &tags); err != nil {
		return err
	}
	for _, tag := range tags {
		var links MarkovLinks
		if err := m.Find(bson.M{"tag": tag}).All(&links); err != nil {
			return err
		}
		// Bolt requires values to be valid for the life of the transaction,
		// so convert int -> Uvarint stored in a private struct field.
		// Mongo will never know!
		for i := range links {
			links[i].uses = make([]byte, binary.MaxVarintLen64)
			n := binary.PutUvarint(links[i].uses, uint64(links[i].Uses))
			links[i].uses = links[i].uses[:n]
		}
		err := mc.bolt.Update(func(tx *bolt.Tx) error {
			return mc.importLinksTx(tx, []byte(tag), links)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (mc *Collection) importLinksTx(tx *bolt.Tx, tag []byte, links MarkovLinks) error {
	// Each tag is a new sub-bucket of the markov collection...
	mb := tx.Bucket([]byte(COLLECTION))
	tb, err := mb.CreateBucketIfNotExists(tag)
	if err != nil {
		return err
	}
	count, skipped := 0, 0
	for _, link := range links {
		// There's a tonne of URLs in the markov data. Skip migrating them.
		if util.LooksURLish(link.Source) || util.LooksURLish(link.Dest) {
			skipped++
			continue
		}
		// And each source is a sub-bucket for the tag,
		sb, err := tb.CreateBucketIfNotExists([]byte(link.Source))
		if err != nil {
			return err
		}
		// Which contains a mapping of dest -> use count.
		if err = sb.Put([]byte(link.Dest), link.uses); err != nil {
			return err
		}
		count++
	}
	logging.Debug("Migrated %d and skipped %d markov links for %s.", count, skipped, tag)
	return nil
}

func (mc *Collection) Migrated() bool {
	return mc.checker.Checker.Migrated()
}

func mongoIndexes(c db.Collection) {
	if err := c.Mongo().EnsureIndex(mgo.Index{
		Key: []string{"tag", "source", "dest"},
	}); err != nil {
		logging.Error("Couldn't create an index on markov: %s", err)
	}
}

func (mc *Collection) incUses(source, dest, tag string) {
	if util.LooksURLish(source) || util.LooksURLish(dest) {
		// Skip URLs entirely.
		return
	}
	// Mongo.
	mlink := New(source, dest, tag)
	if err := mc.mongo.Get(mlink.byTagSrcDest(), mlink); err != nil {
		mlink = New(source, dest, tag)
	}
	mlink.Uses++
	if err := mc.mongo.Put(mlink); err != nil {
		logging.Error("Failed to insert Mongo MarkovLink %s: %v",
			mlink, err)
	}

	// Bolt.
	err := mc.bolt.Update(func(tx *bolt.Tx) error {
		return mc.incUsesTx(tx, uint64(mlink.Uses), []byte(tag), []byte(source), []byte(dest))
	})
	if err != nil {
		logging.Error("Failed to insert Bolt MarkovLink %s(%q->%q): %v",
			tag, source, dest, err)
	}
}

func (mc *Collection) incUsesTx(tx *bolt.Tx, muses uint64, tag, source, dest []byte) error {
	mb := tx.Bucket([]byte(COLLECTION))
	tb, err := mb.CreateBucketIfNotExists(tag)
	if err != nil {
		return err
	}
	sb, err := tb.CreateBucketIfNotExists(source)
	if err != nil {
		return err
	}
	b, n := make([]byte, binary.MaxVarintLen64), 0
	if mc.Migrated() {
		uses, _ := binary.Uvarint(sb.Get(dest))
		uses++
		if muses != uses {
			logging.Warn("Markov mismatch (%d != %d) for %s(%q->%q)",
				uses, muses, tag, source, dest)
		}
		n = binary.PutUvarint(b, uses)
	} else {
		n = binary.PutUvarint(b, muses)
	}
	return sb.Put(dest, b[:n])
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
	_, merr := mc.mongo.Mongo().RemoveAll(bson.M{"tag": tag})
	berr := mc.bolt.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(COLLECTION))
		return mb.DeleteBucket([]byte(tag))
	})
	if merr == nil && berr == nil {
		return nil
	}
	return fmt.Errorf("clearing markov tag %q: merr=%v berr=%v", tag, merr, berr)
}

type MarkovSource struct {
	*Collection
	tag string
}

func (mc *Collection) Source(tag string) markov.Source {
	return &MarkovSource{mc, tag}
}

func (ms *MarkovSource) GetLinks(source string) (markov.Links, error) {
	// Mongo.
	key := &MarkovLink{
		Source: source,
		Tag:    ms.tag,
	}
	var mall MarkovLinks
	merr := ms.mongo.All(key.byTagSrc(), &mall)
	mlinks := make(markov.Links, 0, len(mall))
	for _, link := range mall {
		if util.LooksURLish(link.Source) || util.LooksURLish(link.Dest) {
			// Avoid diffs due to URLs skipped during migration.
			continue
		}
		mlinks = append(mlinks, markov.Link{link.Dest, link.Uses})
	}

	// Bolt.
	blinks := markov.Links{}
	berr := ms.bolt.View(func(tx *bolt.Tx) error {
		return ms.getLinksTx(tx, []byte(ms.tag), []byte(source), &blinks)
	})

	// Diff.
	if merr != nil || berr != nil {
		return nil, fmt.Errorf("markov getlinks(%q, %q): merr=%v berr=%v", ms.tag, source, merr, berr)
	}
	mstr := mlinks.Strings()
	bstr := blinks.Strings()
	sort.Strings(mstr)
	sort.Strings(bstr)
	if d, err := diff.Unified(mstr, bstr); err != nil {
		logging.Warn("markov getlinks(%q, %q): %v\n\n%s\n", ms.tag, source, err, strings.Join(d, "\n"))
	}
	if ms.Migrated() {
		return blinks, nil
	}
	return mlinks, nil
}

func (ms *MarkovSource) getLinksTx(tx *bolt.Tx, tag, source []byte, blinks *markov.Links) error {
	mb := tx.Bucket([]byte(COLLECTION))
	tb := mb.Bucket(tag)
	if tb == nil {
		return fmt.Errorf("couldn't find bucket representing tag %q", tag)
	}
	sb := tb.Bucket(source)
	if sb == nil {
		return fmt.Errorf("couldn't find bucket representing source %q", source)
	}
	return sb.ForEach(func(k, v []byte) error {
		uses, _ := binary.Uvarint(v)
		*blinks = append(*blinks, markov.Link{string(k), int(uses)})
		return nil
	})
}
