package calcdriver

import (
	"github.com/fluffle/golog/logging"
)

const driverName string = "calc"

type calcDriver struct {
	l logging.Logger
}

func CalcDriver(l logging.Logger) *calcDriver {
	return &calcDriver{l}
}

func (cd *calcDriver) Name() string {
	return driverName
}
