// +build kafka allsource all

package kafka

import (
	"github.com/geniusrabbit/eventstream/source"
)

func init() {
	source.RegisterConnector(connector, "kafka")
}

func connector(config *source.Config) (eventstream.Sourcer, error) {
	var (
		url, err   = url.Parse(config.Connect)
		subscriber *kafka.Subscriber
	)

	if err != nil {
		return nil, err
	}

	subscriber, err = kafka.NewSubscriber(
		strings.Split(url.Host, ","),
		url.Path[1:],
		strings.Split(url.Query().Get("topics"), ","),
	)

	if err != nil {
		return nil, err
	}

	if config.Format == "" {
		config.Format = "raw"
	}

	subscriber := &sourceSubscriber{
		debug:      config.Debug,
		subscriber: subscriber,
		format:     converter.ByName(config.Format),
	}

	if err := subscriber.Subscribe(subscriber); err != nil {
		return nil, err
	}
	return subscriber
}
