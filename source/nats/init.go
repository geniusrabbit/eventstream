// +build nats allsource all

package nats

import (
	"github.com/geniusrabbit/eventstream/source"
)

func init() {
	source.RegisterConnector(connector, "nats")
}

func connector(config *source.Config) (eventstream.Sourcer, error) {
	var (
		url, err   = url.Parse(config.Connect)
		subscriber *ncnats.Subscriber
	)

	if err != nil {
		return nil, err
	}

	if config.Format == "" {
		config.Format = "raw"
	}

	subObject := &sourceSubscriber{
		debug:  config.Debug,
		format: converter.ByName(config.Format),
	}

	subscriber, err = ncnats.NewSubscriber(
		"nats://"+url.Host,
		url.Path[1:],
		strings.Split(url.Query().Get("topics"), ","),
		nats.DisconnectHandler(subObject.eventDisconnect),
		nats.ReconnectHandler(subObject.eventReconnect),
		nats.ClosedHandler(subObject.eventClose),
	)

	if err != nil {
		return nil, err
	}

	subObject.subscriber = subscriber
	if err := subscriber.Subscribe(subObject); err != nil {
		return nil, err
	}
	return subObject, nil
}
