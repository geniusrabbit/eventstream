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

// New clickhouse stream
func New(connector Connector, conf *stream.Config, pattern string) (st eventstream.Streamer, err error) {
	var config config

	if err = conf.Decode(&config); err != nil {
		return
	}

	if config.RawQuery != "" {
		st, err = NewStreamSQLByRaw(
			connector,
			int(config.BufferSize),
			time.Duration(config.WriteTimeout)*time.Millisecond,
			config.RawQuery,
			config.Fields,
			conf.Debug,
		)
	} else {
		var q *stream.Query
		if q, err = stream.NewQueryByPattern(pattern, config.Target, config.Fields); err == nil {
			st, err = NewStreamSQL(
				connector,
				int(config.BufferSize),
				time.Duration(config.WriteTimeout)*time.Millisecond,
				*q, conf.Debug,
			)
		}
	}

	return
}
