//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package nats

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/nats"
	"github.com/geniusrabbit/notificationcenter/natstream"
)

// NATS processor
type NATS struct {
	// Debug mode of the storage
	debug bool

	// Stream interface
	stream nc.Streamer
}

// Stream metrics processor
func (m *NATS) Stream(options ...interface{}) (streamObj eventstream.Streamer, err error) {
	var conf stream.Config
	for _, opt := range options {
		switch o := opt.(type) {
		case stream.Option:
			o(&conf)
		case *stream.Config:
			conf = *o
		default:
			stream.WithObjectConfig(o)(&conf)
		}
	}
	if streamObj, err = newStream(m.stream, &conf); err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(streamObj, conf.Where)
}

// Close vertica connection
func (m *NATS) Close() error {
	if cl, _ := m.stream.(io.Closer); cl != nil {
		return cl.Close()
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Connection helpers
///////////////////////////////////////////////////////////////////////////////

func connectNATS(connection string) (nc.Streamer, error) {
	var arr = strings.Split(connection, "?")
	if len(arr) != 2 {
		return nil, fmt.Errorf("Undefined NATS topics: %s", connection)
	}
	var vals, err = url.ParseQuery(arr[1])
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(connection, "nats://") {
		return nats.NewStream(strings.Split(vals.Get("topics"), ","), arr[0])
	}
	return natstream.NewStream(
		arr[0],
		vals.Get("cluster_id"),
		vals.Get("client_id"),
		vals.Get("topic"),
	)
}
