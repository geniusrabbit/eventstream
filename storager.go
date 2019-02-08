//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package eventstream

import (
	"io"
)

// Storager describe method of interaction with storage
type Storager interface {
	// Closer extension of interface
	io.Closer

	// Stream new processor
	Stream(conf interface{}) (Streamer, error)
}
