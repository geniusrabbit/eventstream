//
// @project geniusrabbit::eventstream 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package clickhouse

import (
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/geniusrabbit/eventstream/stream/sql"
)

// New clickhouse stream
func New(connector sql.Connector, config eventstream.ConfigItem, debug bool) (st eventstream.SimpleStreamer, err error) {
	if rawItem := config.String("rawitem", ""); rawItem != "" {
		st, err = sql.NewStreamSQLByRaw(
			connector,
			int(config.Int("buffer", 0)),
			time.Duration(config.Int("duration", 0)),
			rawItem,
			config.Item("fields", nil),
			debug,
		)
	} else {
		var (
			q     *stream.Query
			query = `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`
		)

		if q, err = stream.NewQueryByPattern(query, config.String("target", ""), config.Item("fields", nil)); err == nil {
			st, err = sql.NewStreamSQL(
				connector,
				int(config.Int("buffer", 0)),
				time.Duration(config.Int("duration", 0)),
				*q, debug,
			)
		}
	}
	return
}
