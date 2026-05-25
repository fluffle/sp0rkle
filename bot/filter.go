package bot

import (
	"strings"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/db/conf"
)

type LineFilter interface {
	ShouldProcess(line *client.Line) bool
}

type LineFilterFunc func(line *client.Line) bool

func (lf LineFilterFunc) ShouldProcess(line *client.Line) bool {
	return lf(line)
}

type nickIgnoreFilter struct {
	ns conf.Namespace
}

func (nf nickIgnoreFilter) ShouldProcess(line *client.Line) bool {
	if line.Nick == "" {
		return true
	}
	return nf.ns.String(strings.ToLower(line.Nick)) == ""
}

type FilterPipeline struct {
	filters []LineFilter
}

func (p *FilterPipeline) Add(lf ...LineFilter) {
	p.filters = append(p.filters, lf...)
}

func (p *FilterPipeline) AddFunc(lf ...LineFilterFunc) {
	// Can't type-cast []LineFilterFunc to []LineFilter, rip.
	for _, f := range lf {
		p.filters = append(p.filters, f)
	}
}

func (p *FilterPipeline) ShouldProcess(line *client.Line) bool {
	for _, f := range p.filters {
		if !f.ShouldProcess(line) {
			return false
		}
	}
	return true
}
