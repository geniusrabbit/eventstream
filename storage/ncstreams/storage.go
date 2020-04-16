//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package ncstreams

import (
	"io"

	nc "github.com/geniusrabbit/notificationcenter"

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
