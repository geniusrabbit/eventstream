//
// @project geniusrabbit::eventstream 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package vertica

import (
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/geniusrabbit/eventstream/stream/sql"
)

// New vertica stream
func New(connector sql.Connector, conf *stream.Config) (st eventstream.Streamer, err error) {
	return sql.New(connector, conf, `COPY {{target}} ({{fields}}) FROM STDIN DELIMITER '\t' NULL 'null'`)
}
