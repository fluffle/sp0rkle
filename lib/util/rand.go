package util

import (
	"math/rand"
	"sync"
	"time"
)

// A random number generator seeded with something less predictable than 1
var RNG *rand.Rand = NewRand(time.Now())

// Gratuitously stolen from pkg/rand, cos they aren't usable externally.
type lockedSource struct {
	sync.Mutex
	rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.Lock()
	defer r.Unlock()
	return r.Source.Int63()
}

func (r *lockedSource) Seed(seed int64) {
	r.Lock()
	r.Source.Seed(seed)
	r.Unlock()
}

func NewRand(seed int64) *rand.Rand {
	return rand.New(&lockedSource{Source: rand.NewSource(seed)})

}
