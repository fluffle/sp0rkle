package main

import (
	"github.com/fluffle/goirc/client"
	"strings"
	"time"
)

type FactoidPlugin func(string, *client.Line) string

type PluginManager interface {
	AddPlugin(p FactoidPlugin)
	ApplyPlugins(string, *client.Line) string
}

type PluginProvider interface {
	RegisterPlugins(PluginManager)
}

func (fd *factoidDriver) AddPlugin(p FactoidPlugin) {
	fd.plugins = append(fd.plugins, p)
}

func (fd *factoidDriver) ApplyPlugins(val string, line *client.Line) string {
	for _, p := range fd.plugins {
		val = p(val, line)
	}
	return val
}

func (fd *factoidDriver) RegisterPlugins(pm PluginManager) {
	// pm == fd in this case, but meh.
	pm.AddPlugin(FactoidPlugin(plug_identifiers))
}

// Replicate perlfu's $<stuff> identifiers
func plug_identifiers(val string, line *client.Line) string {
	ts := time.LocalTime()
	val = strings.Replace(val, "$nick", line.Nick, -1)
	val = strings.Replace(val, "$chan", line.Args[0], -1)
	val = strings.Replace(val, "$username", line.Ident, -1)
	val = strings.Replace(val, "$user", line.Ident, -1)
	val = strings.Replace(val, "$host", line.Host, -1)
	val = strings.Replace(val, "$date", ts.Format(time.ANSIC), -1)
	val = strings.Replace(val, "$time", ts.Format("15:04:05"), -1)
	return val
}
