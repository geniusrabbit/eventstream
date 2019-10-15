//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//
// * DEPRECATED module

package metrics

import (
	"errors"
	"log"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/metrics"
)

var (
	errInvalidMetricsItemConfig = errors.New("[metrics] invalid metrics item config")
)

type config struct {
	Metrics []*metricItem `json:"metrics"`
	Prefix  string        `json:"prefix"`
}

type mstream struct {
	debug   bool
	id      string
	prefix  string
	metrics []*metricItem
	metrica notificationcenter.Streamer
}

func newStream(metrica notificationcenter.Streamer, conf *stream.Config) (eventstream.Streamer, error) {
	var preConfig config

	if err := conf.Decode(&preConfig); err != nil {
		return nil, err
	}

	conf.Where = strings.TrimSpace(conf.Where)
	stream := &mstream{
		debug:   conf.Debug,
		id:      conf.Name,
		prefix:  preConfig.Prefix,
		metrics: preConfig.Metrics,
		metrica: metrica,
	}

	if conf.Where != "" {
		return eventstream.NewStreamWrapper(stream, conf.Where)
	}

	return stream, nil
}

// ID returns unical stream identificator
func (s *mstream) ID() string {
	return s.id
}

// Put message to stream
func (s *mstream) Put(msg eventstream.Message) error {
	messages := s.prepareMetricsMessage(msg)
	if s.debug {
		for _, msg := range messages {
			m := msg.(metrics.Message)
			log.Printf("[metrics] %s> %s\n", s.prefix, m.Name)
		}
	}
	return s.metrica.Send(messages...)
}

// Checl the message
func (s *mstream) Check(msg eventstream.Message) bool {
	return true
}

// Close implementation
func (s *mstream) Close() error {
	return nil
}

// Run loop
func (s *mstream) Run() error {
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (s *mstream) prepareMetricsMessage(msg eventstream.Message) (result []interface{}) {
	for _, mt := range s.metrics {
		var (
			name     = mt.Name
			replacer = mt.replacer(msg)
		)

		if replacer != nil {
			name = replacer.Replace(mt.Name)
		}

		result = append(result, metrics.Message{
			Name:  name,
			Type:  mt.getType(),
			Tags:  mt.getTags(replacer),
			Value: msg.Item(mt.Value, nil),
		})
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func gIndexOfStr(s string, arr [][2]string) int {
	for i, sv := range arr {
		if sv[1] == s {
			return i
		}
	}
	return -1
}
