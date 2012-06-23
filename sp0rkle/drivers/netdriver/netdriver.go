package netdriver

import (
	"github.com/fluffle/golog/logging"
)

const driverName string = "net"

type netService interface {
	LookupResult(string) string
}

type netDriver struct {
	services map[string]netService
	l logging.Logger
}

func NetDriver(l logging.Logger) *netDriver {
	nd := &netDriver{make(map[string]netService), l}
	nd.services["calc"] = IGoogleCalcService(l)
	return nd
}

func (nd *netDriver) Name() string {
	return driverName
}
