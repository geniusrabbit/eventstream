package eventstream

import (
	"context"
	"io"
)

// Sourcer ingests messages from a queue (Kafka, NATS, Redis streams, …)
// and fans them out to subscribed streamers.
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
