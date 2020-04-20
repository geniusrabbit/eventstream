package ncstreams

import (
	"context"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/kafka"
	"github.com/geniusrabbit/notificationcenter/nats"
	"github.com/geniusrabbit/notificationcenter/natstream"
)

func connector(config *source.Config) (eventstream.Sourcer, error) {
	return Open(config.Connect, WithDebug(config.Debug), WithFormat(config.Format))
}

// Open new source by URLs
func Open(url string, options ...Option) (eventstream.Sourcer, error) {
	var (
		err        error
		opts       Options
		ctx        = context.Background()
		subscriber nc.Subscriber
	)
	for _, opt := range options {
		opt(&opts)
	}
	switch {
	case strings.HasPrefix(url, `nats://`):
		subscriber, err = nats.NewSubscriber(nats.WithNatsURL(url))
	case strings.HasPrefix(url, `natstream://`):
		subscriber, err = natstream.NewSubscriber(natstream.WithNatsURL(url))
	case strings.HasPrefix(url, `kafka://`):
		subscriber, err = kafka.NewSubscriber(kafka.WithKafkaURL(url))
	}
	if err != nil {
		return nil, err
	}
	subscriberWrapper := &sourceSubscriber{
		debug:      opts.Debug,
		subscriber: subscriber,
		logger:     opts.getLogger(),
		format:     opts.getFormat(),
	}
	if err := subscriber.Subscribe(ctx, subscriberWrapper); err != nil {
		return nil, err
	}
	return subscriberWrapper, nil
}
