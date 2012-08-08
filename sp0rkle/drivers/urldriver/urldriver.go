package urldriver

import (
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/lib/urls"
)

const driverName string = "urls"

type urlDriver struct {
	*urls.UrlCollection
	l logging.Logger
}

func UrlDriver(db *db.Database, l logging.Logger) *urlDriver {
	return &urlDriver{urls.Collection(db, l), l}
}

func (ud *urlDriver) Name() string {
	return driverName
}
