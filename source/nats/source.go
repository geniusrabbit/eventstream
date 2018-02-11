//
// @project geniusrabbit::eventstream 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package nats

import (
	"log"
	"net/url"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/geniusrabbit/eventstream/source"
	ncnats "github.com/geniusrabbit/notificationcenter/nats"
	"github.com/nats-io/nats"
)

func init() {
	source.RegisterConnector(connector, "nats")
}

type sourceSubscriber struct {
	debug      bool
	format     converter.Converter
	subscriber *ncnats.Subscriber
}

func connector(config eventstream.ConfigItem, debug bool) (eventstream.Sourcer, error) {
	var (
		url, err   = url.Parse(config.String("connect", ""))
		subscriber *ncnats.Subscriber
	)

	if err != nil {
		return nil, err
	}

	subObject := &sourceSubscriber{
		debug:  debug,
		format: converter.ByName(config.String("format", "raw")),
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
	return subObject, nil
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

func (s *sourceSubscriber) eventDisconnect(conn *nats.Conn) {
	if s.debug {
		log.Println("Event [disconnect]",
			"closed:", conn.IsClosed(),
			"reconnectiong:", conn.IsReconnecting(),
			conn.Reconnects,
		)
	}
}

func (s *sourceSubscriber) eventReconnect(conn *nats.Conn) {
	if s.debug {
		log.Println("Event [reconnect]",
			"closed:", conn.IsClosed(),
			"reconnectiong:", conn.IsReconnecting(),
			conn.Reconnects,
		)
	}
}

func (s *sourceSubscriber) eventClose(conn *nats.Conn) {
	if s.debug {
		log.Println("Event [close]",
			"closed:", conn.IsClosed(),
			"reconnectiong:", conn.IsReconnecting(),
			conn.Reconnects,
		)
	}
}
