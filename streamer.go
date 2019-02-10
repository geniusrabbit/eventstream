//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package eventstream

import (
	"io"
)

// Streamer interface of data processing
type Streamer interface {
	// Close extension
	io.Closer

	// Put message to stream
	Put(msg Message) error

	// Check message value
	Check(msg Message) bool

	// Run loop
	Run() error
}
