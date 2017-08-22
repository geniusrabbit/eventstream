//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	statsdbase "gopkg.in/alexcesaro/statsd.v2"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/metrics"
	"github.com/geniusrabbit/notificationcenter/nats"
	"github.com/geniusrabbit/notificationcenter/statsd"
)

// Errors set
var (
	ErrUndefinedMetricsEngine = errors.New(`Undefined metrics engine or wrong "connect"`)
)

func init() {
	storage.RegisterConnector(connector, "metrics")
}

// Metrics processor
type Metrics struct {
	metrica notificationcenter.Logger
}

func connector(conf eventstream.ConfigItem, debug bool) (_ eventstream.Storager, err error) {
	var (
		logger  notificationcenter.Logger
		connect = conf.String("connect", "")
	)
	switch {
	case strings.HasPrefix(connect, "nats://"):
		logger, err = connectNATS(connect)
	case strings.HasPrefix(connect, "statsd://"):
		logger, err = connectStatsD(connect)
	default:
		return nil, ErrUndefinedMetricsEngine
	}

	if err != nil {
		return nil, err
	}

	return &Metrics{metrica: logger}, nil
}

// Stream metrics processor
func (m *Metrics) Stream(conf eventstream.ConfigItem) (eventstream.Streamer, error) {
	stream, err := newStream(m.metrica, conf)
	if err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(stream, conf.String("where", ""))
}

// Close vertica connection
func (m *Metrics) Close() error {
	if cl, _ := m.metrica.(io.Closer); cl != nil {
		return cl.Close()
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Connection helpers
///////////////////////////////////////////////////////////////////////////////

func connectNATS(connection string) (notificationcenter.Logger, error) {
	var arr = strings.Split(connection, "?")
	if len(arr) != 2 {
		return nil, fmt.Errorf("Undefined NATS topics: %s", connection)
	}

	var vals, err = url.ParseQuery(arr[1])
	if err != nil {
		return nil, err
	}

	return nats.NewLogger(
		strings.Split(vals.Get("topics"), ","),
		arr[0],
	)
}

func connectStatsD(connection string) (notificationcenter.Logger, error) {
	var (
		url, err = url.Parse(connection)
		tags     []string
	)

	if err != nil {
		return nil, err
	}

	for k, v := range url.Query() {
		if len(v) > 0 {
			tags = append(tags, k, v[0])
		}
	}

	return statsd.NewUDP(
		url.Host,
		metrics.InfluxFormat,
		statsdbase.TagsFormat(statsdbase.InfluxDB),
		statsdbase.Tags(tags...),
	)
}
