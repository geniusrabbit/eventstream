//
// @project geniusrabbit::eventstream 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package nats

import (
	"log"
	"net/url"
	"strings"

	"github.com/nats-io/nats"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/geniusrabbit/eventstream/source"
	ncnats "github.com/geniusrabbit/notificationcenter/nats"
)

var (
	errInvalidStreamObject = errors.New("[nats] invalid stream object")
)

type sourceSubscriber struct {
	debug      bool
	format     converter.Converter
	subscriber *ncnats.Subscriber
	streams    []eventstream.Streamer
}

// Subscribe stream to data processing pipeline
func (s *sourceSubscriber) Subscribe(stream eventstream.Streamer) error {
	if stream == nil {
		return errInvalidStreamObject
	}
	for _, st := range s.streams {
		if st.ID() == stream.ID() {
			return fmt.Errorf("[nats] stream [%s] already registered", st.ID())
		}
	}
	s.streams = append(s.streams, stream)
	return nil
}

// Handle notification message
func (s *sourceSubscriber) Handle(message notificationcenter.Message) error {
	msg, err := eventstream.MessageDecode(message.Data(), s.format)
	if err != nil {
		return err
	}
	for _, stream := range s.streams {
		if !stream.Check(msg) {
			continue
		}
		if err = s.stream.Put(msg); err != nil {
			break
		}
	}
	if err != nil {
		err = message.Ack()
	}
	return err
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
