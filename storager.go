//
// @project geniusrabbit::eventstream 2017, 2019-2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2020
//

package eventstream

import (
	"io"
)

// Storager describe method of interaction with storage.
// Storage creates new stream interfaces to process
// data from sources.
//go:generate mockgen -source $GOFILE -package mocks -destination internal/mocks/storage.go
type Storager interface {
	// Closer extension of interface
	io.Closer

	// Stream returns new stream writer for some specific configs
	Stream(opts ...interface{}) (Streamer, error)
}
