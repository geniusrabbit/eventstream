//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package stream

import (
	"io"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/notificationcenter"
)

// Options stream
type Options struct {
	Connection string
	RawItem    string
	Target     string
	Fields     interface{}
	When       string
	Options    map[string]interface{}
}

// Get value by sub options
func (o Options) Get(key string) (v interface{}) {
	if nil != o.Options {
		v, _ = o.Options[key]
	}
	return
}

// NewConstructor type
type NewConstructor func(opt Options) (Streamer, error)

// Streamer interface accessor
type Streamer interface {
	// Close implementation
	io.Closer

	// Check message value
	Check(msg eventstream.Message) bool

	// Put message to stream
	Put(msg eventstream.Message) error

	// Process loop
	Process()
}

// ExtStreamer ext interface
type ExtStreamer interface {
	// notification center handler
	notificationcenter.Handler

	// Close implementation
	io.Closer

	// Process loop
	Process()
}
