package base

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/lib/db"
)

// Extend goirc's Line with useful extra information
type Line struct {
	client.Line
	Addressed bool
}

func (line *Line) Copy() *Line {
	return &Line{Line: *line.Line.Copy(), Addressed: line.Addressed}
}

func (line *Line) Storable() (db.StorableNick, db.StorableChan) {
	return db.StorableNick{line.Nick, line.Ident, line.Host},
		db.StorableChan{line.Args[0]}
}
