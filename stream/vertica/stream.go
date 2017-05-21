//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package vertica

import (
	"database/sql"
	"time"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/stream"
	bsql "github.com/geniusrabbit/eventstream/stream/sql"
)

// New vertica stream
func New(opt stream.Options) (stream.Streamer, error) {
	if "" != opt.RawItem {
		return NewStreamVerticaByRaw(
			storage.Get(opt.Connection).(*sql.DB),
			gocast.ToInt(opt.Get("buffer")),
			time.Duration(gocast.ToInt(opt.Get("duration"))),
			opt.RawItem,
			opt.Fields,
		)
	}
	return NewStreamVerticaByTarget(
		storage.Get(opt.Connection).(*sql.DB),
		gocast.ToInt(opt.Get("buffer")),
		time.Duration(gocast.ToInt(opt.Get("duration"))),
		opt.Target,
		opt.Fields,
	)
}

// NewStreamVerticaByRaw query
func NewStreamVerticaByRaw(conn *sql.DB, blockSize int, duration time.Duration, query string, fields interface{}) (stream.Streamer, error) {
	return bsql.NewStreamSQLByRaw(conn, blockSize, duration, query, fields)
}

// NewStreamVerticaByTarget params
func NewStreamVerticaByTarget(conn *sql.DB, blockSize int, duration time.Duration, target string, fields interface{}) (stream.Streamer, error) {
	q, err := stream.NewQueryByPattern(`COPY ${target} (${fields}) FROM STDIN DELIMITER '\t' NULL 'null'`, target, fields)
	if nil != err {
		return nil, err
	}
	return bsql.NewStreamSQL(conn, blockSize, duration, *q), nil
}
