//go:build vertica || allsource || all

package driversinit

import (
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/vertica"
)

func init() {
	storage.RegisterConnector("vertica", vertica.Connector)
}
