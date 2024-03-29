package stream

import (
	"encoding/json"

	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/utils/condition"
)

// Option of the stream
type Option func(cnf *Config)

// WithConfig custom
func WithConfig(conf *Config) Option {
	return func(cnf *Config) {
		*cnf = *conf
	}
}

// WithName of the stream
func WithName(name string) Option {
	return func(cnf *Config) {
		cnf.Name = name
	}
}

// WithDebug mode
func WithDebug(debug bool) Option {
	return func(cnf *Config) {
		cnf.Debug = debug
	}
}

// WithWhere condition
func WithWhere(where string) Option {
	return func(cnf *Config) {
		cnf.Where = where
	}
}

// WithCondition condition
func WithCondition(cond condition.Condition) Option {
	return func(cnf *Config) {
		cnf.WhereCond = cond
	}
}

// WithRawConfig storage config
func WithRawConfig(raw json.RawMessage) Option {
	return func(cnf *Config) {
		cnf.Raw = raw
	}
}

// WithMetrics of the stream
func WithMetrics(metrics metrics.MetricList) Option {
	return func(cnf *Config) {
		cnf.Metrics = metrics
	}
}

// WithObjectConfig converts Object to JSON storage config
func WithObjectConfig(obj any) Option {
	return func(cnf *Config) {
		data, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		cnf.Raw = json.RawMessage(data)
	}
}
