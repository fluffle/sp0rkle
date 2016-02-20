package conf

import (
	"reflect"

	"github.com/fluffle/golog/logging"
)

type both struct {
	bolt, mongo *namespace
}

func (b both) All() []Entry {
	mongo := b.mongo.All()
	bolt := b.bolt.All()
	if !reflect.DeepEqual(mongo, bolt) {
		logging.Warn("All() mismatch (%v vs. %v) for ns %q.",
			mongo, bolt, b.mongo.ns)
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
	return mongo
}

func (b both) Int(key string, value ...int) int {
	mongo := b.mongo.Int(key, value...)
	bolt := b.bolt.Int(key, value...)
	if mongo != bolt {
		logging.Warn("Int() mismatch (%d vs. %d) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
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
	return mongo
}

func (b both) Value(key string, value ...interface{}) interface{} {
	mongo := b.mongo.Value(key, value...)
	bolt := b.bolt.Value(key, value...)
	if !reflect.DeepEqual(mongo, bolt) {
		logging.Warn("Value() mismatch (%v vs. %v) for ns %q, key %q.",
			mongo, bolt, b.mongo.ns, key)
	}
	return mongo
}

func (b both) Delete(key string) {
	b.mongo.Delete(key)
	b.bolt.Delete(key)
}
