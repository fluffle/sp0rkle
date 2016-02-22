package conf

import (
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/diff"
)

type migrator struct {
	sync.RWMutex
	migrated bool
}

func (m *migrator) Migrate() error {
	m.Lock()
	defer m.Unlock()
	var all []Entry
	mongo.Init(db.Mongo, COLLECTION, mongoIndexes)
	if err := mongo.All(db.K{}, &all); err != nil {
		return err
	}
	for _, e := range all {
		logging.Debug("Migrating entry %s.", e)
		Bolt(e.Ns).Value(e.Key, e.Value)
	}
	if err := m.diff(); err != nil {
		return err
	}
	m.migrated = true
	return nil
}

func (m *migrator) diff() error {
	var mAll, bAll []Entry
	if err := mongo.All(db.K{}, &mAll); err != nil {
		return err
	}
	if err := bolt.All(db.K{}, &bAll); err != nil {
		return err
	}
	mStr, bStr := make([]string, len(mAll)), make([]string, len(bAll))
	for i, e := range mAll {
		mStr[i] = e.String()
	}
	for i, e := range bAll {
		bStr[i] = e.String()
	}
	sort.Strings(mStr)
	sort.Strings(bStr)
	if len(mAll) != len(bAll) || strings.Join(mStr, "\n") != strings.Join(bStr, "\n") {
		logging.Error("Diff: mlen = %d, blen = %d\n%s\n", len(mAll), len(bAll),
			strings.Join(diff.Unified(mStr, bStr), "\n"))
		return diff.ErrDiff
	}
	return nil
}

func (m *migrator) Migrated() bool {
	m.RLock()
	defer m.RUnlock()
	return m.migrated
}

type both struct {
	bolt, mongo *namespace
	*migrator
}

func (b both) All() []Entry {
	mongo := b.mongo.All()
	bolt := b.bolt.All()
	if !reflect.DeepEqual(mongo, bolt) {
		logging.Warn("All() mismatch (%v vs. %v) for ns %q.",
			mongo, bolt, b.mongo.ns)
	}
	if b.Migrated() {
		return bolt
	}
	return mongo
}

func (b both) String(key string, value ...string) string {
	mongo := b.mongo.String(key, value...)
	bolt := b.bolt.String(key, value...)
	if mongo != bolt {
		logging.Warn("String() mismatch (%q vs. %q) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	if b.Migrated() {
		return bolt
	}
	return mongo
}

func (b both) Int(key string, value ...int) int {
	mongo := b.mongo.Int(key, value...)
	bolt := b.bolt.Int(key, value...)
	if mongo != bolt {
		logging.Warn("Int() mismatch (%d vs. %d) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	if b.Migrated() {
		return bolt
	}
	return mongo
}

func (b both) Float(key string, value ...float64) float64 {
	mongo := b.mongo.Float(key, value...)
	bolt := b.bolt.Float(key, value...)
	if mongo != bolt {
		logging.Warn("Float() mismatch (%f vs. %f) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	if b.Migrated() {
		return bolt
	}
	return mongo
}

func (b both) Value(key string, value ...interface{}) interface{} {
	mongo := b.mongo.Value(key, value...)
	bolt := b.bolt.Value(key, value...)
	if !reflect.DeepEqual(mongo, bolt) {
		logging.Warn("Value() mismatch (%v vs. %v) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	if b.Migrated() {
		return bolt
	}
	return mongo
}

func (b both) Delete(key string) {
	b.mongo.Delete(key)
	b.bolt.Delete(key)
}
