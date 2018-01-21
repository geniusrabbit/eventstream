//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

import (
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/demdxx/gocast"
)

// StreamWrapper extends
type StreamWrapper struct {
	// Stream pricessor
	Stream SimpleStreamer

	// WhereCondition of stream
	WhereCondition *govaluate.EvaluableExpression
}

// NewStreamWrapper with support condition
func NewStreamWrapper(stream SimpleStreamer, where string) (_ Streamer, err error) {
	var whereObj *govaluate.EvaluableExpression

	if len(strings.TrimSpace(where)) > 0 {
		if whereObj, err = govaluate.NewEvaluableExpression(where); nil != err {
			return
		}
	}

	return &StreamWrapper{
		Stream:         stream,
		WhereCondition: whereObj,
	}, nil
}

// Put message to stream
func (s *StreamWrapper) Put(msg Message) error {
	return s.Stream.Put(msg)
}

// Check message
func (s *StreamWrapper) Check(msg Message) bool {
	if s.WhereCondition != nil {
		r, err := s.WhereCondition.Evaluate(msg.Map())
		return err == nil && gocast.ToBool(r)
	}
	return true
}

// Run loop
func (s *StreamWrapper) Run() error {
	return s.Stream.Run()
}

// Close stream ans shut down all process
func (s *StreamWrapper) Close() error {
	return s.Stream.Close()
}