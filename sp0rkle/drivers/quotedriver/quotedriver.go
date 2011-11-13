package quotedriver

import (
	"github.com/fluffle/golog/logging"
//	"launchpad.net/gobson/bson"
	"lib/db"
	"lib/quotes"
//	"lib/util"
//	"sp0rkle/base"
//	"strings"
)

const driverName string = "quotes"

type quoteDriver struct {
	*quotes.QuoteCollection

	// logging object
	l logging.Logger
}

func QuoteDriver(db *db.Database, l logging.Logger) *quoteDriver {
	qc := quotes.Collection(db, l)
	return &quoteDriver{
		QuoteCollection: qc,
		l:               l,
	}
}

func (qd *quoteDriver) Name() string {
	return driverName
}
