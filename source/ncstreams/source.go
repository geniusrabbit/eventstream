//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream/converter"
	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream"
)

var (
	errInvalidStreamObject     = errors.New("invalid stream object")
	errStreamAlreadyRegistered = errors.New("stream already registered")
)

type sourceSubscriber struct {
	debug      bool
	format     converter.Converter
	subscriber nc.Subscriber
	streams    []eventstream.Streamer
}

// Subscribe new stream object
func (s *sourceSubscriber) Subscribe(ctx context.Context, stream eventstream.Streamer) error {
	if stream == nil {
		return errInvalidStreamObject
	}
	for _, st := range s.streams {
		if st.ID() == stream.ID() {
			return errors.Wrap(errStreamAlreadyRegistered, st.ID())
		}
	}
	s.streams = append(s.streams, stream)
	return nil
}

// Receive notification message
func (s *sourceSubscriber) Receive(message nc.Message) error {
	msg, err := eventstream.MessageDecode(message.Body(), s.format)
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, stream := range s.streams {
		if stream.Check(ctx, msg) {
			if err = stream.Put(ctx, msg); err != nil {
				return err
			}
		}
	}
	return message.Ack()
}

// Start sunscriber listener
func (s *sourceSubscriber) Start(ctx context.Context) error {
	go func() {
		_ = s.subscriber.Listen(ctx)
	}()
	return nil
}

// Close source subscriber
func (s *sourceSubscriber) Close() error {
	return s.subscriber.Close()
}
