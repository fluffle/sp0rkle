package base

import (
	"github.com/fluffle/goirc/client"
)

// Extend goirc's Line with useful extra information
type Line struct {
	client.Line
	Addressed bool
}

func (line *Line) Copy() *Line {
	return &Line{Line: *line.Line.Copy(), Addressed: line.Addressed}
}
