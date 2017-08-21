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
	ErrUndefinedMetricsEngine = errors.New("Undefined metrics engine")
)

func init() {
	storage.RegisterConnector(connector, "metrics")
}

// Metrics processor
type Metrics struct {
	metrica notificationcenter.Logger
}

func connector(conf eventstream.ConfigItem, debug bool) (_ *Metrics, err error) {
	var (
		logger     notificationcenter.Logger
		connection = conf.String("connection", "")
	)
	switch {
	case strings.HasPrefix(connection, "nats://"):
		logger, err = connectNATS(connection)
	case strings.HasPrefix(connection, "statsd://"):
		logger, err = connectStatsD(connection)
	default:
		return nil, ErrUndefinedMetricsEngine
	}

	if err != nil {
		return nil, err
	}

	return &Metrics{metrica: logger}, nil
}

// Write messages to storage
func (m *Metrics) Write(message eventstream.Message) error {
	return m.metrica.Send(m.prepareMetricsMessage(message)...)
}

// Close vertica connection
func (m *Metrics) Close() error {
	if cl, _ := m.metrica.(io.Closer); cl != nil {
		return cl.Close()
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func (m *Metrics) prepareMetricsMessage(msg eventstream.Message) (result []interface{}) {
	return
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
