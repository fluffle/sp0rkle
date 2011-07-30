package factdriver

import (
	"launchpad.net/gobson/bson"
	"lib/db"
	"lib/factoids"
	"log"
	"sp0rkle/base"
)

const driverName string = "factoids"

type factoidDriver struct {
	*factoids.FactoidCollection

	// Keep a reference to the last factoid looked up around
	// for use with 'edit that' and 'delete that' commands.
	lastseen bson.ObjectId

	// A list of text processing plugins to apply to factoid values
	plugins []base.Plugin
}

func FactoidDriver(db *db.Database) *factoidDriver {
	fc, err := factoids.Collection(db)
	if err != nil {
		log.Fatalf("factoid collection failed: %v\n", err)
	}
	return &factoidDriver{
		FactoidCollection: fc,
		lastseen:          "",
		plugins:           make([]base.Plugin, 0),
	}
}

func (fd *factoidDriver) Name() string {
	return driverName
}

func (fd *factoidDriver) AddPlugin(p base.Plugin) {
	fd.plugins = append(fd.plugins, p)
}

func (fd *factoidDriver) ApplyPlugins(val string, line *base.Line) string {
	for _, p := range fd.plugins {
		val = p.Apply(val, line)
	}
	return val
}
