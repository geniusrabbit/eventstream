// +build clickhouse allsource all

package clickhouse

import (
	"github.com/geniusrabbit/eventstream/storage"
)

func init() {
	storage.RegisterConnector(connector, "clickhouse")
}
