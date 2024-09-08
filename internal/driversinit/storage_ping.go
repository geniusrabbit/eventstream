//go:build ping || http || allsource || all

package driversinit

import (
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/ping"
)

func init() {
	storage.RegisterConnector("ping", ping.Connector)
}
