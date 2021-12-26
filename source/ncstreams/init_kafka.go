//go:build kafka || allsource || all
// +build kafka allsource all

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/kafka"
)

func kafkaNewSubscriber(ctx context.Context, url string) (nc.Subscriber, error) {
	return kafka.NewSubscriber(kafka.WithKafkaURL(url))
}

// OpenKafka new source by URLs
func OpenKafka(ctx context.Context, url string, options ...Option) (eventstream.Sourcer, error) {
	return Open(ctx, url, kafkaNewSubscriber, options...)
}

func init() {
	source.RegisterConnector("kafka", func(ctx context.Context, config *source.Config) (eventstream.Sourcer, error) {
		return OpenKafka(ctx, config.Connect, WithDebug(config.Debug), WithFormat(config.Format))
	})
}
