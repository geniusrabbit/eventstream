package sql

import (
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
)

type config struct {
	RawQuery     string      `json:"raw_query"`
	Target       string      `json:"target"`
	BufferSize   uint        `json:"buffer_size"`
	WriteTimeout uint        `json:"write_timeout"`
	Fields       interface{} `json:"fields"`
}

// New stream for SQL type integrations
func New(connector Connector, pattern string, conf *stream.Config, options ...Option) (stream eventstream.Streamer, err error) {
	var (
		config      config
		queryOption Option
	)
	if err = conf.Decode(&config); err != nil {
		return
	}
	if config.RawQuery != "" {
		queryOption = WithQueryRawFields(config.RawQuery, config.Fields)
	} else {
		queryOption = WithQueryByPattern(pattern, config.Target, config.Fields)
	}
	return NewStreamSQL(
		conf.Name,
		connector,
		append(
			options,
			queryOption,
			WithBlockSize(int(config.BufferSize)),
			WithFlushIntervals(time.Duration(config.WriteTimeout)*time.Millisecond),
			WithDebug(conf.Debug),
		)...,
	)
}
