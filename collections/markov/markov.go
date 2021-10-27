package markov

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/util/diff"
	"github.com/fluffle/sp0rkle/util/markov"
	bolt "go.etcd.io/bbolt"
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

func (ml *MarkovLink) encodeUses() {
	ml.uses = make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(ml.uses, uint64(ml.Uses))
	ml.uses = ml.uses[:n]
}

func (ml *MarkovLink) decodeUses() {
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

func (mc *Collection) MigrateTo(newState db.MigrationState) error {
	if newState != db.MONGO_PRIMARY {
		return nil
	}
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
		// Bolt requires values to be valid for the life of the transaction,
		// so convert int -> Uvarint stored in a private struct field.
		// Mongo will never know!
		link.encodeUses()
		if err = sb.Put([]byte(link.Dest), link.uses); err != nil {
			return err
		}
		count++
	}
	logging.Debug("Migrated %d and skipped %d markov links for %s.", count, skipped, tag)
	return nil
}

func (mc *Collection) Check() db.MigrationState {
	return mc.checker.Checker.Check()
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
	mlink := New(source, dest, tag)
	blink := New(source, dest, tag)
	state := mc.Check()
	// Read current value.
	if state < db.BOLT_ONLY {
		if err := mc.mongo.Get(mlink.byTagSrcDest(), mlink); err != nil {
			mlink = New(source, dest, tag)
		}
	}
	if state > db.MONGO_ONLY {
		mc.bolt.View(func(tx *bolt.Tx) error {
			return mc.getUsesTx(tx, blink)
		})
	}
	// Diff if we're reading both.
	if (state == db.MONGO_PRIMARY || state == db.BOLT_PRIMARY) &&
		mlink.Uses != blink.Uses {
		logging.Warn("Markov link uses mismatch (%d != %d) for %s(%q->%q)",
			mlink.Uses, blink.Uses, tag, source, dest)
	}
	mlink.Uses++
	// Increment and write new value.
	if state < db.BOLT_ONLY {
		if err := mc.mongo.Put(mlink); err != nil {
			logging.Error("Failed to insert Mongo MarkovLink %s: %v",
				mlink, err)
		}
	}
	if state > db.MONGO_ONLY {
		err := mc.bolt.Update(func(tx *bolt.Tx) error {
			return mc.putUsesTx(tx, mlink)
		})
		if err != nil {
			logging.Error("Failed to insert Bolt MarkovLink %s(%q->%q): %v",
				tag, source, dest, err)
		}
	}
}

func (mc *Collection) getUsesTx(tx *bolt.Tx, link *MarkovLink) error {
	mb := tx.Bucket([]byte(COLLECTION))
	tb := mb.Bucket([]byte(link.Tag))
	if tb == nil {
		return fmt.Errorf("couldn't find bucket representing tag %q", link.Tag)
	}
	sb := tb.Bucket([]byte(link.Source))
	if sb == nil {
		return fmt.Errorf("couldn't find bucket representing source %q", link.Source)
	}
	v := sb.Get([]byte(link.Dest))
	if v == nil {
		return fmt.Errorf("couldn't find key representing dest %q", link.Dest)
	}
	uses, _ := binary.Uvarint(v)
	link.Uses = int(uses)
	return nil
}

func (mc *Collection) putUsesTx(tx *bolt.Tx, link *MarkovLink) error {
	mb := tx.Bucket([]byte(COLLECTION))
	tb, err := mb.CreateBucketIfNotExists([]byte(link.Tag))
	if err != nil {
		return err
	}
	sb, err := tb.CreateBucketIfNotExists([]byte(link.Source))
	if err != nil {
		return err
	}
	link.encodeUses()
	return sb.Put([]byte(link.Dest), link.uses)
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
	var mErr, bErr error
	if mc.Check() < db.BOLT_ONLY {
		_, mErr = mc.mongo.Mongo().RemoveAll(bson.M{"tag": tag})
	}
	if mc.Check() > db.MONGO_ONLY {
		bErr = mc.bolt.Update(func(tx *bolt.Tx) error {
			mb := tx.Bucket([]byte(COLLECTION))
			return mb.DeleteBucket([]byte(tag))
		})
	}
	if mErr == nil && bErr == nil {
		return nil
	}
	return fmt.Errorf("clearing markov tag %q: mongo=%v bolt=%v", tag, mErr, bErr)
}

type MarkovSource struct {
	*Collection
	tag string
}

func (mc *Collection) Source(tag string) markov.Source {
	return &MarkovSource{mc, tag}
}

func (ms *MarkovSource) GetLinks(source string) (markov.Links, error) {
	mLinks, bLinks := markov.Links{}, markov.Links{}
	var mErr, bErr error
	state := ms.Check()
	if state < db.BOLT_ONLY {
		// Read from mongo.
		key := &MarkovLink{
			Source: source,
			Tag:    ms.tag,
		}
		var mAll MarkovLinks
		mErr = ms.mongo.All(key.byTagSrc(), &mAll)
		for _, link := range mAll {
			if util.LooksURLish(link.Source) || util.LooksURLish(link.Dest) {
				// Avoid diffs due to URLs skipped during migration.
				continue
			}
			mLinks = append(mLinks, markov.Link{link.Dest, link.Uses})
		}
	}
	if state > db.MONGO_ONLY {
		// Read from bolt.
		bErr = ms.bolt.View(func(tx *bolt.Tx) error {
			return ms.getLinksTx(tx, []byte(ms.tag), []byte(source), &bLinks)
		})
	}
	// If either failed, bail out!
	if mErr != nil || bErr != nil {
		return nil, fmt.Errorf("markov getlinks(%q, %q): merr=%v berr=%v", ms.tag, source, mErr, bErr)
	}
	// Diff if we're reading from both.
	if state == db.MONGO_PRIMARY || state == db.BOLT_PRIMARY {
		if unified, err := diff.SortDiff(mLinks, bLinks); err == diff.ErrDiff {
			logging.Warn("markov getlinks(%q, %q): %v\n\n%s\n", ms.tag, source, err, strings.Join(unified, "\n"))
		}
	}
	if state >= db.BOLT_PRIMARY {
		return bLinks, nil
	}
	return mLinks, nil
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
