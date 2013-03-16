package bot

import (
	"sync"
)

// A Rewriter rewrites lines sent back to the server via
// ReplyN(), Reply() and Do()
type Rewriter interface {
	Rewrite(input string, context *Context) (output string)
}

type RewriteFunc func(string, *Context) string

func (rwf RewriteFunc) Rewrite(in string, ctx *Context) string {
	return rwf(in, ctx)
}

type RewriteSet interface {
	Rewriter
	Add(Rewriter)
}

type rewriteSet struct {
	sync.RWMutex
	set []Rewriter
}

func newRewriteSet() *rewriteSet {
	return &rewriteSet{set: make([]Rewriter, 0, 10)}
}

func (rws *rewriteSet) Add(rw Rewriter) {
	rws.Lock()
	defer rws.Unlock()
	rws.set = append(rws.set, rw)
}

func (rws *rewriteSet) Rewrite(in string, ctx *Context) string {
	rws.RLock()
	defer rws.RUnlock()
	for _, rw := range rws.set {
		in = rw.Rewrite(in, ctx)
	}
	return in
}
