package netdriver

import (
)

type netService interface {
	LookupResult(string) string
}

func Init() {
}
