//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package eventstream

import (
	"context"
	"io"

	"github.com/geniusrabbit/eventstream/stream"
)

// StreamConfig of the stream
type StreamConfig = stream.Config

// Streamer interface of data processing describes
// basic methods of data pipeline
//go:generate mockgen -source $GOFILE -package mocks -destination internal/mocks/stream.go
type Streamer interface {
	// Close extension
	io.Closer

	// ID returns unical stream identificator
	ID() string

	// Put message to the stream to process information
	Put(ctx context.Context, msg Message) error

	// Check if message suits for the stream
	Check(ctx context.Context, msg Message) bool

	// Run processing loop
	Run(ctx context.Context) error
}
