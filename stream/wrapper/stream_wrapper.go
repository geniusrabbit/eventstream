//
// @project geniusrabbit::eventstream 2017, 2019 - 2023
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019 - 2023
//

package wrapper

import (
	"context"

	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/geniusrabbit/eventstream/utils/condition"
)

// StreamWrapper extends any stream interface with conditional
// check support to any message
type StreamWrapper struct {
	// Stream pricessor
	stream stream.Streamer

	// WhereCondition of stream
	whereCondition condition.Condition

	// Metrics executor
	metrics metrics.Metricer
}

// NewStreamWrapper with support condition
func NewStreamWrapper(stream stream.Streamer, whereObj condition.Condition, metrics metrics.Metricer) stream.Streamer {
	if whereObj == nil && metrics == nil {
		return stream
	}
	return &StreamWrapper{
		stream:         stream,
		whereCondition: whereObj,
		metrics:        metrics,
	}
}

// ID returns unical stream identificator
func (s *StreamWrapper) ID() string {
	return s.stream.ID()
}

// Put message to the stream to process information
func (s *StreamWrapper) Put(ctx context.Context, msg message.Message) error {
	if s.metrics != nil {
		s.metrics.Exec(ctx, msg)
	}
	return s.stream.Put(ctx, msg)
}

// Check if the message meets the conditions
func (s *StreamWrapper) Check(ctx context.Context, msg message.Message) bool {
	if s.stream == nil {
		return true
	}
	if !s.stream.Check(ctx, msg) {
		return false
	}
	if s.whereCondition != nil {
		return s.whereCondition.Check(ctx, msg)
	}
	return true
}

// Run the stream reading loop
func (s *StreamWrapper) Run(ctx context.Context) error {
	return s.stream.Run(ctx)
}

// Close stream and shut down all process
func (s *StreamWrapper) Close() error {
	return s.stream.Close()
}
