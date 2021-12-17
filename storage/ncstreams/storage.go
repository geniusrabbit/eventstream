//
// @project geniusrabbit::eventstream 2017, 2019, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019, 2020
//

package ncstreams

import (
	"context"
	"io"

	nc "github.com/geniusrabbit/notificationcenter"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
)

type getPublisherFnk func(ctx context.Context, url string) (nc.Publisher, error)

// PublishStorage processor
type PublishStorage struct {
	// Debug mode of the storage
	debug bool

	// Stream interface
	publisher nc.Publisher
}

// Open new storage connection
func Open(ctx context.Context, url string, pubFnk getPublisherFnk, options ...Option) (eventstream.Storager, error) {
	var (
		opts           Options
		publisher, err = pubFnk(ctx, url)
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
