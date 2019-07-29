package sql

import (
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

type config struct {
	RawQuery     string      `json:"raw_query"`
	Target       string      `json:"target"`
	BufferSize   uint        `json:"buffer_size"`
	WriteTimeout uint        `json:"write_timeout"`
	Fields       interface{} `json:"fields"`
}

// New stream for SQL type integrations
func New(connector Connector, conf *storage.StreamConfig, pattern string) (stream eventstream.Streamer, err error) {
	var config config

	if err = conf.Decode(&config); err != nil {
		return
	}

	if config.RawQuery != "" {
		stream, err = NewStreamSQLByRaw(
			conf.Name,
			connector,
			config.RawQuery,
			config.Fields,
			WithBlockSize(int(config.BufferSize)),
			WithFlushIntervals(time.Duration(config.WriteTimeout)*time.Millisecond),
			WithDebug(conf.Debug),
		)
	} else {
		var query *Query
		if query, err = NewQueryByPattern(pattern, config.Target, config.Fields); err == nil {
			stream, err = NewStreamSQL(
				conf.Name,
				connector,
				*query,
				WithBlockSize(int(config.BufferSize)),
				WithFlushIntervals(time.Duration(config.WriteTimeout)*time.Millisecond),
				WithDebug(conf.Debug),
			)
		}
	}
	return
}
