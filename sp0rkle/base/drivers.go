package base

import (
	"github.com/fluffle/goirc/event"
)

// Interface for a driver
type Driver interface {
	Name() string
	RegisterHandlers(event.EventRegistry)
}

type Plugin interface {
	Apply(string, *Line) string
}

type PluginManager interface {
	AddPlugin(p Plugin)
	ApplyPlugins(string, *Line) string
}

type PluginProvider interface {
	RegisterPlugins(PluginManager)
}
