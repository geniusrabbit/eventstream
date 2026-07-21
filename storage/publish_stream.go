package storage

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream/pkg/message"
	"github.com/geniusrabbit/eventstream/stream"
	nc "github.com/geniusrabbit/notificationcenter/v2"
)

type nstream struct {
	id        string
	templates []*messageTemplate
	publisher nc.Publisher
}

func newStream(pub nc.Publisher, conf *streamOptions) (stream.Streamer, error) {
	if len(conf.targets) == 0 {
		return nil, errInvalidMetricsItemConfig
	}
	strm := &nstream{
		id:        conf.id,
		publisher: pub,
	}
	for _, target := range conf.targets {
		template, err := newMessageTemplate(target.Fields, target.Where)
		if err != nil {
			return nil, errors.Wrap(errInvalidMetricsItemConfig, err.Error())
		}
		strm.templates = append(strm.templates, template)
	}
	return strm, nil
}

// ID returns unical stream identificator
func (s *nstream) ID() string {
	return s.id
}

// Put message to stream
func (s *nstream) Put(ctx context.Context, msg message.Message) error {
	return s.publisher.Publish(ctx, s.prepareMessages(msg)...)
}

// Check the message
func (s *nstream) Check(ctx context.Context, msg message.Message) bool {
	for _, tmp := range s.templates {
		if tmp.check(msg) {
			return true
		}
	}
	return false
}

// Close implementation
func (s *nstream) Close() error {
	if closer, _ := s.publisher.(io.Closer); closer != nil {
		return closer.Close()
	}
	return nil
}

// Run loop
func (s *nstream) Run(ctx context.Context) error {
	return nil
}

func (s *nstream) prepareMessages(msg message.Message) (result []any) {
	for _, mt := range s.templates {
		if prepared := mt.prepare(msg); prepared != nil {
			result = append(result, prepared)
		}
	}
	return result
}
