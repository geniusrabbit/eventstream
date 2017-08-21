//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

import "io"

// Sourcer interface
type Sourcer interface {
	// Close extension
	io.Closer

	// Subscribe stream object
	Subscribe(stream Streamer) error

	// Start listeners
	Start() error
}
