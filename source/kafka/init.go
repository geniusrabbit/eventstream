// +build kafka allsource all

package kafka

import (
	"net/url"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/notificationcenter/kafka"
)

func init() {
	source.RegisterConnector(connector, "kafka")
}

func connector(config *source.Config) (eventstream.Sourcer, error) {
	var (
		url, err = url.Parse(config.Connect)
		kafkaSub *kafka.Subscriber
	)

	if err != nil {
		return nil, err
	}

	kafkaSub, err = kafka.NewSubscriber(
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
		subscriber: kafkaSub,
		format:     converter.ByName(config.Format),
	}

	if err := kafkaSub.Subscribe(subscriber); err != nil {
		return nil, err
	}
	return subscriber, nil
}
