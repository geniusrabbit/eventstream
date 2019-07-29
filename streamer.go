//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package eventstream

import (
	"io"
)

// Streamer interface of data processing describes
// basic methods of data pipeline
type Streamer interface {
	// Close extension
	io.Closer

	// ID returns unical stream identificator
	ID() string

	// Put message to the stream to process information
	Put(msg Message) error

	// Check if message suits for the stream
	Check(msg Message) bool

	// Run processing loop
	Run() error
}
