// Package stream defines the Streamer sink interface used by sources
// and the pipeline Engine.
package stream

import (
	"context"
	"io"

	"github.com/geniusrabbit/eventstream/pkg/message"
)

// Streamer is a named sink that accepts, filters and processes messages.
type Streamer interface {
	// Close extension
	io.Closer

	// ID returns unical stream identificator
	ID() string

	// Put message to the stream to process information
	Put(ctx context.Context, msg message.Message) error

	// Check if message suits for the stream
	Check(ctx context.Context, msg message.Message) bool

	// Run processing loop
	Run(ctx context.Context) error
}
