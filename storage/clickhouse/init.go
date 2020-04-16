// build clickhouse allsource all

package clickhouse

import (
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

func connector(conf *storage.Config) (eventstream.Storager, error) {
	return Open(conf.Connect, storage.WithDebug(conf.Debug))
}

func init() {
	storage.RegisterConnector(connector, "clickhouse")
}
