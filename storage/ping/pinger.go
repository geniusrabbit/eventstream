package ping

import (
	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
	"github.com/geniusrabbit/eventstream/stream"
)

type pinger struct {
	URL    string
	Method string
}

// Close implements eventstream.Storager.
func (p *pinger) Close() error {
	return nil
}

// Stream implements eventstream.Storager.
func (p *pinger) Stream(options ...any) (stream.Streamer, error) {
	var (
		err        error
		conf       stream.Config
		strmConf   pingStreamConfig
		metricExec metrics.Metricer
	)
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
	if err = conf.Decode(&strmConf); err != nil {
		return nil, err
	}
	if metricExec, err = conf.Metrics.Metric(); err != nil {
		return nil, err
	}
	strm := &pingStream{
		url:         patternkey.PatternKeyFromTemplate(gocast.Or(strmConf.URL, p.URL)),
		method:      gocast.Or(strmConf.Method, p.Method),
		contentType: gocast.Or(strmConf.ContentType, "application/json"),
	}
	cond, err := conf.Condition()
	if err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(strm, cond, metricExec), nil
}

var _ eventstream.Storager = (*pinger)(nil)
