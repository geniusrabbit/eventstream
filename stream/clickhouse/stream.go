//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package clickhouse

import (
	"database/sql"
	"time"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/stream"
	bsql "github.com/geniusrabbit/eventstream/stream/sql"
)

// New clickhouse stream
func New(opt stream.Options) (stream.Streamer, error) {
	if "" == opt.RawItem {
		return NewStreamClickhouseByTarget(
			storage.Get(opt.Connection).(*sql.DB),
			gocast.ToInt(opt.Get("buffer")),
			time.Duration(gocast.ToInt(opt.Get("duration"))),
			opt.When,
			opt.Target,
			opt.Fields,
		)
	}
	return bsql.NewStreamSQLByRaw(
		storage.Get(opt.Connection).(*sql.DB),
		gocast.ToInt(opt.Get("buffer")),
		time.Duration(gocast.ToInt(opt.Get("duration"))),
		opt.When,
		opt.RawItem,
		opt.Fields,
	)
}

// NewStreamClickhouseByTarget params
func NewStreamClickhouseByTarget(
	conn *sql.DB,
	blockSize int,
	duration time.Duration,
	when,
	target string,
	fields interface{},
) (stream.Streamer, error) {
	q, err := stream.NewQueryByPattern(`INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`, target, fields)
	if nil != err {
		return nil, err
	}
	return bsql.NewStreamSQL(conn, blockSize, duration, when, *q)
}
