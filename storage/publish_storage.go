package storage

import (
	"context"
	"io"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
)

type getPublisherFnk func(ctx context.Context, url string) (nc.Publisher, error)

// PublishStorage is a Storager that publishes messages through a notificationcenter Publisher.
type PublishStorage struct {
	publisher nc.Publisher
}

// Open connects a publish storage using the given publisher factory.
func Open(ctx context.Context, url string, pubFnk getPublisherFnk) (eventstream.Storager, error) {
	publisher, err := pubFnk(ctx, url)
	if err != nil {
		return nil, err
	}
	return &PublishStorage{publisher: publisher}, nil
}

// Stream returns a publisher stream configured with StreamOption values.
func (m *PublishStorage) Stream(options ...any) (stream.Streamer, error) {
	var conf streamOptions
	for _, opt := range options {
		switch o := opt.(type) {
		case StreamOption:
			o(&conf)
		default:
			return nil, errors.Wrapf(ErrInvalidOption, `publish stream: %T`, opt)
		}
	}
	return newStream(m.publisher, &conf)
}

// Close publisher connection
func (m *PublishStorage) Close() (err error) {
	if cl, _ := m.publisher.(io.Closer); cl != nil {
		err = cl.Close()
	}
	return err
}
