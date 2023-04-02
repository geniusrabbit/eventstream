//go:build clickhouse || allsource || all

package driversinit

import (
	_ "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/clickhouse"
)

func init() {
	storage.RegisterConnector("clickhouse", clickhouse.Connector)
}
