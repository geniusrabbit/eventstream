//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package vertica

import (
	"database/sql"
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	bsql "github.com/geniusrabbit/eventstream/stream/sql"
)

// New vertica stream
func New(store eventstream.Storager, conn *sql.DB, config eventstream.ConfigItem, debug bool) (stream eventstream.SimpleStreamer, err error) {
	if rawItem := config.String("rawitem", ""); rawItem != "" {
		stream, err = bsql.NewStreamSQLByRaw(
			conn,
			int(config.Int("buffer", 0)),
			time.Duration(config.Int("duration", 0)),
			rawItem,
			config.Item("fields", nil),
			debug,
		)
	} else {
		stream, err = newStreamVerticaByTarget(
			conn,
			int(config.Int("buffer", 0)),
			time.Duration(config.Int("duration", 0)),
			config.String("target", ""),
			config.Item("fields", nil),
			debug,
		)
	}
	return
}

// NewStreamVerticaByTarget params
func newStreamVerticaByTarget(conn *sql.DB, blockSize int, duration time.Duration, target string, fields interface{}, debug bool) (eventstream.SimpleStreamer, error) {
	q, err := stream.NewQueryByPattern(
		`COPY {{target}} ({{fields}}) FROM STDIN DELIMITER '\t' NULL 'null'`,
		target, fields,
	)
	if nil != err {
		return nil, err
	}
	return bsql.NewStreamSQL(conn, blockSize, duration, *q, debug)
}
