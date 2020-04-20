// +build vertica allsource all

package vertica

import (
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

func connector(conf *storage.Config) (eventstream.Storager, error) {
	return Open(conf.Connect, WithDebug(conf.Debug))
}

func init() {
	storage.RegisterConnector(connector, "vertica")
}
