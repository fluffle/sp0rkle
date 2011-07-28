package main

import (
	"github.com/fluffle/goirc/event"
)

// Interface for a driver
type Driver interface {
	RegisterHandlers(event.EventRegistry)
}
