//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package stream

import (
	"context"
	"io"

	"github.com/geniusrabbit/eventstream/internal/message"
)

// Streamer interface of data processing describes
// basic methods of data pipeline
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
