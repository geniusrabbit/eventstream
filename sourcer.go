//
// @project geniusrabbit::eventstream 2017, 2019-2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2020
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
//
//go:generate mockgen -source $GOFILE -package mocks -destination internal/mocks/source.go
type Sourcer interface {
	// Close extension
	io.Closer

	// Subscribe new stream to data processing.
	// For all subscribed streams sends the same data messages
	Subscribe(ctx context.Context, streams ...Streamer) error

	// Start runs observing for data writing into subscribed streams
	Start(ctx context.Context) error
}
