package db

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/util/diff"
)

const COLLECTION = "migrate"

type Migrator interface {
	Migrate() error
}

type Differ interface {
	Diff() (before, after []string, err error)
}

type Checker interface {
	Migrated() bool
}

type checkFunc func() bool

func (cf checkFunc) Migrated() bool {
	return cf()
}

type M struct {
	Checker
	sync.Once
}

func (m *M) Init(mig Migrator, coll string) {
	m.Do(func() {
		m.Checker = addMigrator(mig, coll)
	})
}

type done struct {
	// Public for serialization purposes.
	Migrated bool
}

func migrated(coll string) bool {
	var d done
	if err := ms.db.Get(K{{"collection", coll}}, &d); err != nil {
		logging.Warn("Checking migrated status for %q: %v", coll, err)
	}
	return d.Migrated
}

type migrator struct {
	Migrator
	migrated bool
}

var ms = &struct {
	sync.RWMutex
	migrators map[string]*migrator
	db        C
}{migrators: make(map[string]*migrator)}

func addMigrator(m Migrator, coll string) Checker {
	ms.Lock()
	defer ms.Unlock()
	// Store migration tracking in new db only.
	ms.db.Init(Bolt, COLLECTION, nil)
	checker := checkFunc(func() bool {
		ms.RLock()
		defer ms.RUnlock()
		if m, ok := ms.migrators[coll]; ok {
			return m.migrated
		}
		return false
	})
	if _, ok := ms.migrators[coll]; ok {
		logging.Warn("Second call to MigratorSet.Add for %q.", coll)
		return checker
	}
	ms.migrators[coll] = &migrator{m, migrated(coll)}
	return checker
}

func Migrate() error {
	ms.db.Init(Bolt, COLLECTION, nil)

	// Holding the lock while migrating prevents the Checker returned by
	// addMigrator from checking migration state (and thus locks up the
	// bot) while migration is running in the background.
	migrators := map[string]*migrator{}
	ms.RLock()
	for coll, m := range ms.migrators {
		migrators[coll] = m
	}
	ms.RUnlock()

	failed := []string{}
	for coll, m := range migrators {
		if m.migrated {
			continue
		}
		logging.Debug("Migrating %q.", coll)
		if err := m.Migrate(); err != nil {
			logging.Error("Migrating %q failed: %v", coll, err)
			failed = append(failed, coll)
			continue
		}
		if differ, ok := m.Migrator.(Differ); ok {
			before, after, err := differ.Diff()
			if err != nil {
				logging.Error("Diffing %q failed: %v", coll, err)
				failed = append(failed, coll)
				continue
			}
			sort.Strings(before)
			sort.Strings(after)
			unified, err := diff.Unified(before, after)
			if err != nil {
				logging.Error("Migration diff: %v\n%s", err, strings.Join(unified, "\n"))
				failed = append(failed, coll)
				continue
			}
		}
		// This is probably a little more locking than strictly necessary.
		ms.Lock()
		if err := ms.db.Put(K{{"collection", coll}}, &done{true}); err != nil {
			logging.Warn("Setting migrated status for %q: %v", coll, err)
		}
		m.migrated = true
		ms.Unlock()
	}
	if len(failed) > 0 {
		return fmt.Errorf("migration failed for: \"%s\"",
			strings.Join(failed, "\", \""))
	}
	return nil
}
