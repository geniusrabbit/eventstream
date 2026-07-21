// Package kafka opens a Kafka publish Storager.
package kafka

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	nckafka "github.com/geniusrabbit/notificationcenter/v2/kafka"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

func newPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return nckafka.NewPublisher(ctx, nckafka.WithKafkaURL(url))
}

// Open connects a publish Storager.
func Open(ctx context.Context, url string) (eventstream.Storager, error) {
	return storage.Open(ctx, url, newPublisher)
}
