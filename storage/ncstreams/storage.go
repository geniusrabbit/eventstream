//
// @project geniusrabbit::eventstream 2017, 2019, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019, 2020
//

package ncstreams

import (
	"context"
	"io"
	"strings"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/kafka"
	"github.com/geniusrabbit/notificationcenter/nats"
	"github.com/geniusrabbit/notificationcenter/natstream"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
)

// PublishStorage processor
type PublishStorage struct {
	// Debug mode of the storage
	debug bool

	// Stream interface
	publisher nc.Publisher
}

// Open new storage connection
func Open(url string, options ...Option) (eventstream.Storager, error) {
	var (
		opts           Options
		ctx            = context.Background()
		publisher, err = connect(ctx, url)
	)
	if err != nil {
		return nil, err
	}
	for _, opt := range options {
		opt(&opts)
	}
	return &PublishStorage{publisher: publisher, debug: opts.Debug}, nil
}

// Stream metrics processor
func (m *PublishStorage) Stream(options ...interface{}) (streamObj eventstream.Streamer, err error) {
	var conf stream.Config
	for _, opt := range options {
		switch o := opt.(type) {
		case stream.Option:
			o(&conf)
		case *stream.Config:
			conf = *o
		default:
			stream.WithObjectConfig(o)(&conf)
		}
	}
	if streamObj, err = newStream(m.publisher, &conf); err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(streamObj, conf.Where)
}

// Close vertica connection
func (m *PublishStorage) Close() (err error) {
	if cl, _ := m.publisher.(io.Closer); cl != nil {
		err = cl.Close()
	}
	return err
}

func connect(ctx context.Context, connection string) (nc.Publisher, error) {
	switch {
	case strings.HasPrefix(connection, "nats://"):
		return nats.NewPublisher(nats.WithNatsURL(connection))
	case strings.HasPrefix(connection, "natstream://"):
		return natstream.NewPublisher(natstream.WithNatsURL(connection))
	case strings.HasPrefix(connection, "kafka://"):
		return kafka.NewPublisher(ctx, kafka.WithKafkaURL(connection))
	}
	return nil, nc.ErrUndefinedPublisherInterface
}
