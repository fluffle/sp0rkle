package decisiondriver

// A simple driver to implement decisions based on random numbers. No, not 4.

import (
	"github.com/fluffle/golog/logging"
)

const driverName string = "decisions"

type decisionDriver struct {
	l *logging.Logger
}

func DecisionDriver(l *logging.Logger) *decisionDriver {
	return &decisionDriver{l}
}

func (dd *decisionDriver) Name() string {
	return driverName
}
