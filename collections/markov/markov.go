package markov

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/util/markov"
	"go.etcd.io/bbolt"
)

const COLLECTION = "markov"

type Collection struct {
	// Markov is a bit special. Because of the quantity of data
	// stored, it is desirable to be able to use the boltdb key
	// for storage too -- something that the standard Collection
	// interface does not provide for, in the possibly misguided
	// name of API simplicity. So instead of hacking this in, we
	// skip the entire db layer and deal with it here instead.
	bolt *bbolt.DB
}

// Wrapper to get hold of a factoid collection handle
func Init() *Collection {
	mc := &Collection{}
	mc.bolt = db.Bolt.DB()

	err := mc.bolt.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(COLLECTION))
		return err
	})
	if err != nil {
		logging.Fatal("Creating Markov BoltDB bucket failed: %v", err)
	}
	return mc
}

func (mc *Collection) incUses(source, dest, tag string) {
	if util.LooksURLish(source) || util.LooksURLish(dest) {
		// Skip URLs entirely.
		return
	}
	err := mc.bolt.Update(func(tx *bbolt.Tx) error {
		mb := tx.Bucket([]byte(COLLECTION))
		tb, err := mb.CreateBucketIfNotExists([]byte(tag))
		if err != nil {
			return err
		}
		sb, err := tb.CreateBucketIfNotExists([]byte(source))
		if err != nil {
			return err
		}

		v := sb.Get([]byte(dest))
		var uses uint64
		if v != nil {
			uses, _ = binary.Uvarint(v)
		}
		uses++

		newV := make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(newV, uses)
		return sb.Put([]byte(dest), newV[:n])
	})
	if err != nil {
		logging.Error("Failed to increment uses for %s(%q->%q): %v",
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
	return mc.bolt.Update(func(tx *bbolt.Tx) error {
		mb := tx.Bucket([]byte(COLLECTION))
		return mb.DeleteBucket([]byte(tag))
	})
}

type MarkovSource struct {
	*Collection
	tag string
}

func (mc *Collection) Source(tag string) markov.Source {
	return &MarkovSource{mc, tag}
}

func (ms *MarkovSource) GetLinks(source string) (markov.Links, error) {
	bLinks := markov.Links{}
	err := ms.bolt.View(func(tx *bbolt.Tx) error {
		return ms.getLinksTx(tx, []byte(ms.tag), []byte(source), &bLinks)
	})
	if err != nil {
		return nil, fmt.Errorf("markov getlinks(%q, %q): %v", ms.tag, source, err)
	}
	return bLinks, nil
}

func (ms *MarkovSource) getLinksTx(tx *bbolt.Tx, tag, source []byte, blinks *markov.Links) error {
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
