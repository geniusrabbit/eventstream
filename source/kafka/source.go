//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package kafka

import (
	"net/url"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/notificationcenter/kafka"
)

type sourceSubscriber struct {
	debug      bool
	format     converter.Converter
	subscriber *kafka.Subscriber
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

	return &sourceSubscriber{
		debug:      config.Debug,
		subscriber: subscriber,
		format:     converter.ByName(config.Format),
	}, nil
}

// Subscribe stream object
func (s *sourceSubscriber) Subscribe(stream eventstream.Streamer) error {
	return s.subscriber.Subscribe(&subs{
		debug:  s.debug,
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
