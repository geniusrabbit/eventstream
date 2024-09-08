//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package ncstreams

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
	nc "github.com/geniusrabbit/notificationcenter/v2"
)

type nstream struct {
	debug     bool
	id        string
	templates []*messageTemplate
	stream    nc.Publisher
}

func newStream(pub nc.Publisher, conf *stream.Config) (eventstream.Streamer, error) {
	var preConfig config
	if err := conf.Decode(&preConfig); err != nil {
		return nil, err
	}
	if len(preConfig.Targets) == 0 {
		return nil, errInvalidMetricsItemConfig
	}
	stream := &nstream{
		debug:  conf.Debug,
		id:     conf.Name,
		stream: pub,
	}
	for _, target := range preConfig.Targets {
		var fields map[string]any
		if len(target.Fields) > 0 {
			fields = target.Fields
		}
		template, err := newMessageTemplate(fields, target.Where)
		if err != nil {
			return nil, errors.Wrap(errInvalidMetricsItemConfig, err.Error())
		}
		stream.templates = append(stream.templates, template)
	}
	return stream, nil
}

// ID returns unical stream identificator
func (s *nstream) ID() string {
	return s.id
}

// Put message to stream
func (s *nstream) Put(ctx context.Context, msg eventstream.Message) error {
	return s.stream.Publish(ctx, s.prepareMessages(msg)...)
}

// Check the message
func (s *nstream) Check(ctx context.Context, msg eventstream.Message) bool {
	for _, tmp := range s.templates {
		if tmp.check(msg) {
			return true
		}
	}
	return false
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

func (s *nstream) prepareMessages(msg eventstream.Message) (result []any) {
	for _, mt := range s.templates {
		if msg := mt.prepare(msg); msg != nil {
			result = append(result, msg)
		}
	}
	return result
}
