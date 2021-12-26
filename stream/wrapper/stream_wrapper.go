//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package wrapper

import (
	"context"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/demdxx/gocast"

	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/stream"
)

// StreamWrapper extends any stream interface with conditional
// check support to any message
type StreamWrapper struct {
	// Stream pricessor
	stream stream.Streamer

	// WhereCondition of stream
	whereCondition *govaluate.EvaluableExpression

	// Metrics executor
	metrics metrics.Metricer
}

// NewStreamWrapper with support condition
func NewStreamWrapper(stream stream.Streamer, where string, metrics metrics.Metricer) (_ stream.Streamer, err error) {
	var whereObj *govaluate.EvaluableExpression

	if len(strings.TrimSpace(where)) > 0 {
		if whereObj, err = govaluate.NewEvaluableExpression(where); err != nil {
			return nil, err
		}
	}

	if whereObj == nil && metrics == nil {
		return stream, nil
	}

	return &StreamWrapper{
		stream:         stream,
		whereCondition: whereObj,
		metrics:        metrics,
	}, nil
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
		r, err := s.whereCondition.Evaluate(msg.Map())
		return err == nil && gocast.ToBool(r)
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
