package base

import (
	"github.com/fluffle/goirc/client"
	"strings"
)

// Basic types representing the information we want to store about IRC things
type Nick string

func (n Nick) Lower() string {
	return strings.ToLower(string(n))
}

type Chan string

func (c Chan) Lower() string {
	return strings.ToLower(string(c))
}

// Extend goirc's Line with useful extra information
type Line struct {
	*client.Line
	Addressed bool
}

func (line *Line) Copy() *Line {
	return &Line{Line: line.Line.Copy(), Addressed: line.Addressed}
}

func (line *Line) Storable() (Nick, Chan) {
	return Nick(line.Nick), Chan(line.Args[0])
}
