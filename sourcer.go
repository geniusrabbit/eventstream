//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package eventstream

import (
	"context"
	"io"
)

// Sourcer describes the input stream interface.
// The source accepts messages from some queue popeline
// like Kafka, NATS, RabbitMQ and etc and send this data
// one by one into the stream processor.
type Sourcer interface {
	// Close extension
	io.Closer

	// Subscribe new stream to data processing.
	// For all subscribed streams sends the same data messages
	Subscribe(ctx context.Context, stream Streamer) error

	// Start runs observing for data writing into subscribed streams
	Start(ctx context.Context) error
}
