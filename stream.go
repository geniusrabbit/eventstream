//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

import "io"

// SimpleStreamer short interface accessor for wrapper use
type SimpleStreamer interface {
	// Close extension
	io.Closer

	// Put message to stream
	Put(msg Message) error

	// Run loop
	Run() error
}

// Streamer basic interface
type Streamer interface {
	SimpleStreamer

	// Check message value
	Check(msg Message) bool
}
