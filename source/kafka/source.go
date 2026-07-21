// Package kafka opens a Kafka Sourcer.
package kafka

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	nckafka "github.com/geniusrabbit/notificationcenter/v2/kafka"
)

func newSubscriber(ctx context.Context, url string) (nc.Subscriber, error) {
	return nckafka.NewSubscriber(nckafka.WithKafkaURL(url))
}

// Open connects a Sourcer.
func Open(ctx context.Context, url string, options ...source.Option) (eventstream.Sourcer, error) {
	return source.Open(ctx, url, newSubscriber, options...)
}
