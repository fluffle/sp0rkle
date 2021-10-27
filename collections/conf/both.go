package conf

import (
	"reflect"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/util/diff"
)

type migrator struct{}

func (migrator) MigrateTo(newState db.MigrationState) error {
	if newState != db.MONGO_PRIMARY {
		return nil
	}
	var all []Entry
	mongo.Init(db.Mongo, COLLECTION, mongoIndexes)
	bolt.Init(db.Bolt.Keyed(), COLLECTION, nil)
	if err := mongo.All(db.K{}, &all); err != nil {
		return err
	}
	if err := bolt.BatchPut(all); err != nil {
		logging.Error("Migrating conf entries: %v.", err)
		return err
	}
	logging.Debug("Migrated %d conf entries.", len(all))
	return nil
}

func (migrator) Diff() ([]string, []string, error) {
	mongo.Init(db.Mongo, COLLECTION, mongoIndexes)
	bolt.Init(db.Bolt.Keyed(), COLLECTION, nil)
	var mAll, bAll Entries
	if err := mongo.All(db.K{}, &mAll); err != nil {
		return nil, nil, err
	}
	if err := bolt.All(db.K{}, &bAll); err != nil {
		return nil, nil, err
	}
	return mAll.Strings(), bAll.Strings(), nil
}

type both struct {
	bolt, mongo *namespace
	db.Checker
}

func (b both) All() Entries {
	switch b.Check() {
	case db.MONGO_ONLY:
		return b.mongo.All()
	case db.BOLT_ONLY:
		return b.bolt.All()
	}
	mAll := b.mongo.All()
	bAll := b.bolt.All()
	unified, err := diff.SortDiff(mAll, bAll)
	if err == diff.ErrDiff {
		logging.Warn("All() mismatch for ns %q (-mongo, +bolt): %v\n%s",
			b.mongo.ns, err, unified)
	}
	if b.Check() >= db.BOLT_PRIMARY {
		return bAll
	}
	return mAll
}

func (b both) String(key string, value ...string) string {
	switch b.Check() {
	case db.MONGO_ONLY:
		return b.mongo.String(key, value...)
	case db.BOLT_ONLY:
		return b.bolt.String(key, value...)
	}
	mongo := b.mongo.String(key, value...)
	bolt := b.bolt.String(key, value...)
	if mongo != bolt {
		logging.Warn("String() mismatch (%q vs. %q) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	if b.Check() >= db.BOLT_PRIMARY {
		return bolt
	}
	return mongo
}

func (b both) Int(key string, value ...int) int {
	switch b.Check() {
	case db.MONGO_ONLY:
		return b.mongo.Int(key, value...)
	case db.BOLT_ONLY:
		return b.bolt.Int(key, value...)
	}
	mongo := b.mongo.Int(key, value...)
	bolt := b.bolt.Int(key, value...)
	if mongo != bolt {
		logging.Warn("Int() mismatch (%d vs. %d) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	if b.Check() >= db.BOLT_PRIMARY {
		return bolt
	}
	return mongo
}

func (b both) Float(key string, value ...float64) float64 {
	switch b.Check() {
	case db.MONGO_ONLY:
		return b.mongo.Float(key, value...)
	case db.BOLT_ONLY:
		return b.bolt.Float(key, value...)
	}
	mongo := b.mongo.Float(key, value...)
	bolt := b.bolt.Float(key, value...)
	if mongo != bolt {
		logging.Warn("Float() mismatch (%f vs. %f) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	if b.Check() >= db.BOLT_PRIMARY {
		return bolt
	}
	return mongo
}

func (b both) Value(key string, value ...interface{}) interface{} {
	switch b.Check() {
	case db.MONGO_ONLY:
		return b.mongo.Value(key, value...)
	case db.BOLT_ONLY:
		return b.bolt.Value(key, value...)
	}
	mongo := b.mongo.Value(key, value...)
	bolt := b.bolt.Value(key, value...)
	if !reflect.DeepEqual(mongo, bolt) {
		logging.Warn("Value() mismatch (%v vs. %v) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	if b.Check() >= db.BOLT_PRIMARY {
		return bolt
	}
	return mongo
}

func (b both) Delete(key string) {
	b.mongo.Delete(key)
	b.bolt.Delete(key)
}
