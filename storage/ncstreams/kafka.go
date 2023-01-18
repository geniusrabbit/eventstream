//go:build kafka || allstorage || all
// +build kafka allstorage all

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/kafka"
)

func kafkaNewPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return kafka.NewPublisher(ctx, kafka.WithKafkaURL(url))
}

// OpenKafka publisher connectior
func OpenKafka(ctx context.Context, url string, options ...storage.Option) (eventstream.Storager, error) {
	return Open(ctx, url, kafkaNewPublisher, options...)
}
