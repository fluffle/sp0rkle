package db

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/fluffle/goirc/logging"
	"github.com/fluffle/sp0rkle/util/diff"
)

const COLLECTION = "migrate"

type MigrationState int

const (
	INVALID_STATE MigrationState = iota - 1
	MONGO_ONLY
	MONGO_PRIMARY
	BOLT_PRIMARY
	BOLT_ONLY
)

var migrationStates = []string{
	"MONGO_ONLY",
	"MONGO_PRIMARY",
	"BOLT_PRIMARY",
	"BOLT_ONLY",
}

var ErrInvalidState = errors.New("invalid migration state")

func StateForName(s string) MigrationState {
	for state, name := range migrationStates {
		if name == s {
			return MigrationState(state)
		}
	}
	return INVALID_STATE
}

func (s MigrationState) Valid() bool {
	return s >= MONGO_ONLY && s <= BOLT_ONLY
}

func (s MigrationState) String() string {
	if !s.Valid() {
		return fmt.Sprintf("invalid state %d", s)
	}
	return migrationStates[int(s)]
}

type Migrator interface {
	MigrateTo(MigrationState) error
}

type Differ interface {
	Diff() (before, after []string, err error)
}

type Checker interface {
	Check() MigrationState
}

type checkFunc func() MigrationState

func (cf checkFunc) Check() MigrationState {
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
	collection string
	// Public for serialization purposes.
	State MigrationState
}

func (d *done) K() Key {
	return K{S{"collection", d.collection}}
}

func getMigrationState(coll string) MigrationState {
	d := &done{collection: coll}
	if err := ms.db.Get(d.K(), d); err != nil {
		logging.Warn("Checking migrated status for %q: %v", coll, err)
	}
	return d.State
}

type migrator struct {
	Migrator
	state MigrationState
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
	ms.db.Init(Bolt.Keyed(), COLLECTION, nil)
	checker := checkFunc(func() MigrationState {
		ms.RLock()
		defer ms.RUnlock()
		if m, ok := ms.migrators[coll]; ok {
			return m.state
		}
		return MONGO_ONLY
	})
	if _, ok := ms.migrators[coll]; ok {
		logging.Warn("Second call to MigratorSet.Add for %q.", coll)
		return checker
	}
	state := getMigrationState(coll)
	ms.migrators[coll] = &migrator{m, state}
	logging.Debug("Added migrator for %s, current state == %s.", coll, state)
	return checker
}

func MigrateTo(newState MigrationState) error {
	if !newState.Valid() {
		return ErrInvalidState
	}

	ms.db.Init(Bolt.Keyed(), COLLECTION, nil)

	// Holding the lock while migrating prevents the Checker returned by
	// addMigrator from checking migration state (and thus locks up the
	// bot) while migration is running in the background.
	migrators := map[string]*migrator{}
	ms.RLock()
	for coll, m := range ms.migrators {
		migrators[coll] = m
	}
	logging.Debug("Migrating %d collections to %s.", len(migrators), newState)
	ms.RUnlock()

	failed := []string{}
	for coll, m := range migrators {
		if m.state >= newState {
			logging.Debug("Skipping %s as it is in %s already.", coll, m.state)
			continue
		}
		logging.Debug("Migrating %q to state %s.", coll, newState)
		if err := m.MigrateTo(newState); err != nil {
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
		if err := ms.db.Put(&done{collection: coll, State: newState}); err != nil {
			logging.Warn("Setting migrated status for %q: %v", coll, err)
		}
		m.state = newState
		ms.Unlock()
	}
	if len(failed) > 0 {
		return fmt.Errorf("migration failed for: \"%s\"",
			strings.Join(failed, "\", \""))
	}
	return nil
}
