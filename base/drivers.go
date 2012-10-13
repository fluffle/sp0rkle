package base

import (
	"github.com/fluffle/goevent/event"
)

// Interface for a driver
type Driver interface {
	Name() string
	RegisterHandlers(event.EventRegistry)
}

type PluginManager interface {
	Add(p Plugin)
	Apply(string, *Line) string
}

type PluginProvider interface {
	RegisterPlugins(PluginManager)
}

type HttpProvider interface {
	RegisterHttpHandlers()
}
