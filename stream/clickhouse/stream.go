//
// @project geniusrabbit::eventstream 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package clickhouse

import (
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/geniusrabbit/eventstream/stream/sql"
)

// New clickhouse stream
func New(connector sql.Connector, conf *stream.Config) (st eventstream.Streamer, err error) {
	return sql.New(connector, conf, `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`)
}
