package pipeline

import (
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/pkg/condition"
	"github.com/geniusrabbit/eventstream/stream"
)

type streamConfig struct {
	sources      []string
	storageName  string
	storageOpts  []any
	streamer     stream.Streamer
	where        string
	whereCond    condition.Condition
	metrics      metrics.MetricList
	mapper       Mapper
	haveStorage  bool
	haveStreamer bool
}

// StreamOption configures a pipeline stream binding.
type StreamOption func(*streamConfig)

// WithSources selects one or more registered sources to subscribe.
func WithSources(names ...string) StreamOption {
	return func(c *streamConfig) {
		c.sources = append(c.sources, names...)
	}
}

// WithStorage builds the sink via a registered storage.Stream(opts...).
func WithStorage(name string, streamOpts ...any) StreamOption {
	return func(c *streamConfig) {
		c.storageName = name
		c.storageOpts = streamOpts
		c.haveStorage = true
	}
}

// WithStreamer uses a pre-built streamer as the sink.
func WithStreamer(s stream.Streamer) StreamOption {
	return func(c *streamConfig) {
		c.streamer = s
		c.haveStreamer = true
	}
}

// WithWhere sets a where-expression condition.
func WithWhere(expr string) StreamOption {
	return func(c *streamConfig) { c.where = expr }
}

// WithCondition sets a pre-built condition.
func WithCondition(cond condition.Condition) StreamOption {
	return func(c *streamConfig) { c.whereCond = cond }
}

// WithMetrics attaches metric executors to the stream wrapper.
func WithMetrics(list metrics.MetricList) StreamOption {
	return func(c *streamConfig) { c.metrics = list }
}

// WithMapper wraps the sink with a message mapper.
func WithMapper(m Mapper) StreamOption {
	return func(c *streamConfig) { c.mapper = m }
}
