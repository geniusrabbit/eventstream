//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package ncstreams

import (
	"context"
	"errors"
	"io"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	nc "github.com/geniusrabbit/notificationcenter"
)

var (
	errInvalidMetricsItemConfig = errors.New("[metrics] invalid metrics item config")
)

type configItem struct {
	Fields map[string]string `json:"fields"`
	Where  string            `json:"where"`
}

type config struct {
	Target []configItem `json:"targets"`
}

type nstream struct {
	debug            bool
	id               string
	messageTemplates []*messageTemplate
	stream           nc.Publisher
}

func newStream(pub nc.Publisher, conf *stream.Config) (eventstream.Streamer, error) {
	var preConfig config
	if err := conf.Decode(&preConfig); err != nil {
		return nil, err
	}
	stream := &nstream{
		debug:  conf.Debug,
		id:     conf.Name,
		stream: pub,
	}
	return stream, nil
}

// ID returns unical stream identificator
func (s *nstream) ID() string {
	return s.id
}

// Put message to stream
func (s *nstream) Put(ctx context.Context, msg eventstream.Message) error {
	messages := s.prepareMessages(msg)
	return s.stream.Publish(ctx, messages...)
}

// Check the message
func (s *nstream) Check(ctx context.Context, msg eventstream.Message) bool {
	return true
}

// Close implementation
func (s *nstream) Close() error {
	if closer, _ := s.stream.(io.Closer); closer != nil {
		return closer.Close()
	}
	return nil
}

// Run loop
func (s *nstream) Run(ctx context.Context) error {
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (s *nstream) prepareMessages(msg eventstream.Message) (result []interface{}) {
	for _, mt := range s.messageTemplates {
		if msg := mt.prepare(msg); msg != nil {
			result = append(result, msg)
		}
	}
	return result
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func gIndexOfStr(s string, arr [][2]string) int {
	for i, sv := range arr {
		if sv[1] == s {
			return i
		}
	}
	return -1
}
