//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package kafka

import (
	"errors"
	"fmt"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/kafka"
)

var (
	errInvalidStreamObject = errors.New("[kafka] invalid stream object")
)

type sourceSubscriber struct {
	debug      bool
	format     converter.Converter
	subscriber *kafka.Subscriber
	streams    []eventstream.Streamer
}

// Subscribe new stream object
func (s *sourceSubscriber) Subscribe(stream eventstream.Streamer) error {
	if stream == nil {
		return errInvalidStreamObject
	}
	for _, st := range s.streams {
		if st.ID() == stream.ID() {
			return fmt.Errorf("[kafka] stream [%s] already registered", st.ID())
		}
	}
	s.streams = append(s.streams, stream)
	return nil
}

// Handle notification message
func (s *sourceSubscriber) Handle(message notificationcenter.Message) error {
	msg, err := eventstream.MessageDecode(message.Body(), s.format)
	if err != nil {
		return err
	}
	for _, stream := range s.streams {
		if !stream.Check(msg) {
			continue
		}
		if err = stream.Put(msg); err != nil {
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
