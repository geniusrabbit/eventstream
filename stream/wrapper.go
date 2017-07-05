//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package stream

import (
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
)

// Wrapper implementation of basic streamer
type Wrapper struct {
	converter converter.Converter
	stream    Streamer
}

// NewWrapper handler
func NewWrapper(stream Streamer, converter converter.Converter) ExtStreamer {
	return &Wrapper{
		stream:    stream,
		converter: converter,
	}
}

// Handle item processin
func (l *Wrapper) Handle(item interface{}) error {
	msg, err := eventstream.MessageDecode(item, l.converter)
	if err == nil && l.stream.Check(msg) {
		err = l.stream.Put(msg)
	}
	return err
}

// Process loop
func (l *Wrapper) Process() {
	l.stream.Process()
}

// Close implementation
func (l *Wrapper) Close() error {
	return l.stream.Close()
}
