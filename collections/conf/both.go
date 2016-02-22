package conf

import (
	"reflect"

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db"
)

type migrator struct{}

func (migrator) Migrate() error {
	var all []Entry
	mongo.Init(db.Mongo, COLLECTION, mongoIndexes)
	if err := mongo.All(db.K{}, &all); err != nil {
		return err
	}
	for _, e := range all {
		logging.Debug("Migrating conf entry %s.", e)
		Bolt(e.Ns).Value(e.Key, e.Value)
	}
	return nil
}

func (migrator) Diff() ([]string, []string, error) {
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
