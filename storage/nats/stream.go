//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package nats

import (
	"errors"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/geniusrabbit/notificationcenter"
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
	stream           notificationcenter.Streamer
}

func newStream(natsStream notificationcenter.Streamer, conf *stream.Config) (eventstream.Streamer, error) {
	var preConfig config
	if err := conf.Decode(&preConfig); err != nil {
		return nil, err
	}
	stream := &nstream{
		debug:  conf.Debug,
		id:     conf.Name,
		stream: natsStream,
	}
	return stream, nil
}

// ID returns unical stream identificator
func (s *nstream) ID() string {
	return s.id
}

// Put message to stream
func (s *nstream) Put(msg eventstream.Message) error {
	messages := s.prepareMessages(msg)
	return s.stream.Send(messages...)
}

// Check the message
func (s *nstream) Check(msg eventstream.Message) bool {
	return true
}

// Close implementation
func (s *nstream) Close() error {
	return nil
}

// Run loop
func (s *nstream) Run() error {
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
