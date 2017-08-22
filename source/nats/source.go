//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package nats

import (
	"net/url"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/notificationcenter/nats"
)

func init() {
	source.RegisterConnector(connector, "nats")
}

type sourceSubscriber struct {
	format     converter.Converter
	subscriber *nats.Subscriber
}

func connector(config eventstream.ConfigItem, debug bool) (eventstream.Sourcer, error) {
	var (
		url, err   = url.Parse(config.String("connect", ""))
		subscriber *nats.Subscriber
	)

	if err != nil {
		return nil, err
	}

	subscriber, err = nats.NewSubscriber(
		"nats://"+url.Host,
		url.Path[1:],
		strings.Split(url.Query().Get("topics"), ","),
	)

	if err != nil {
		return nil, err
	}

	return &sourceSubscriber{
		subscriber: subscriber,
		format:     converter.ByName(config.String("format", "raw")),
	}, nil
}

// Subscribe stream object
func (s *sourceSubscriber) Subscribe(stream eventstream.Streamer) error {
	return s.subscriber.Subscribe(&subs{
		format: s.format,
		stream: stream,
	})
}

// Start sunscriber listener
func (s *sourceSubscriber) Start() error {
	go s.subscriber.Listen()
	return nil
}

// Close source subscriber
func (s *sourceSubscriber) Close() error {
	return s.subscriber.Close()
}
