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
	nl := &Line{Line: *line.Line.Copy()}
	nl.Addressed = line.Addressed
	return line
}
