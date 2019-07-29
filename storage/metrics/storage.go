//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package metrics

import (
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

// Metrics processor
type Metrics struct {
	debug   bool
	metrica notificationcenter.Streamer
}

// Stream metrics processor
func (m *Metrics) Stream(conf interface{}) (eventstream.Streamer, error) {
	return newStream(m.metrica, conf.(*storage.StreamConfig))
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

func connectNATS(connection string) (notificationcenter.Streamer, error) {
	var arr = strings.Split(connection, "?")
	if len(arr) != 2 {
		return nil, fmt.Errorf("Undefined NATS topics: %s", connection)
	}

	var vals, err = url.ParseQuery(arr[1])
	if err != nil {
		return nil, err
	}

	return nats.NewStream(
		strings.Split(vals.Get("topics"), ","),
		arr[0],
	)
}

func connectStatsD(connection string) (notificationcenter.Streamer, error) {
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
