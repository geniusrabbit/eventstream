//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package eventstream

import (
	"context"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/demdxx/gocast"
)

// StreamWrapper extends any stream interface with conditional
// check support to any message
type StreamWrapper struct {
	// Stream pricessor
	stream Streamer

	// WhereCondition of stream
	whereCondition *govaluate.EvaluableExpression
}

// NewStreamWrapper with support condition
func NewStreamWrapper(stream Streamer, where string) (_ Streamer, err error) {
	var whereObj *govaluate.EvaluableExpression

	if len(strings.TrimSpace(where)) > 0 {
		if whereObj, err = govaluate.NewEvaluableExpression(where); err != nil {
			return
		}
	}

	if whereObj == nil {
		return stream, nil
	}

	return &StreamWrapper{
		stream:         stream,
		whereCondition: whereObj,
	}, nil
}

// ID returns unical stream identificator
func (s *StreamWrapper) ID() string {
	return s.stream.ID()
}

// Put message to the stream to process information
func (s *StreamWrapper) Put(ctx context.Context, msg Message) error {
	return s.stream.Put(ctx, msg)
}

// Check if the message meets the conditions
func (s *StreamWrapper) Check(ctx context.Context, msg Message) bool {
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
