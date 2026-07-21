package eventstream

import (
	"io"
)

// Storager opens Streamer sinks that write processed events
// to a concrete backend (SQL, queue, HTTP, …).
//
//go:generate mockgen -source $GOFILE -package mocks -destination internal/mocks/storage.go
type Storager interface {
	// Closer extension of interface
	io.Closer

	// Stream returns new stream writer for some specific configs
	Stream(opts ...any) (Streamer, error)
}
