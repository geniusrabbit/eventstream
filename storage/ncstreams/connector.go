package ncstreams

import (
	"context"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/kafka"
	"github.com/geniusrabbit/notificationcenter/nats"
	"github.com/geniusrabbit/notificationcenter/natstream"
)

func connect(ctx context.Context, connection string) (nc.Publisher, error) {
	switch {
	case strings.HasPrefix(connection, "nats://"):
		return nats.NewPublisher(nats.WithNatsURL(connection))
	case strings.HasPrefix(connection, "natstream://"):
		return natstream.NewPublisher(natstream.WithNatsURL(connection))
	case strings.HasPrefix(connection, "kafka://"):
		return kafka.NewPublisher(ctx, kafka.WithKafkaURL(connection))
	}
	return nil, nc.ErrUndefinedPublisherInterface
}

func connector(conf *storage.Config) (eventstream.Storager, error) {
	ctx := context.Background()
	publisher, err := connect(ctx, conf.Connect)
	if err != nil {
		return nil, err
	}
	return &PublishStorage{publisher: publisher, debug: conf.Debug}, nil
}
