// Package source provides a shared notificationcenter-based subscriber
// used by concrete source drivers (source/nats, source/kafka, …).
package source

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/pkg/converter"
	"github.com/geniusrabbit/eventstream/stream"
)

var (
	errInvalidStreamObject     = errors.New("invalid stream object")
	errStreamAlreadyRegistered = errors.New("stream already registered")
)

type getSubscriberFnk func(ctx context.Context, url string) (nc.Subscriber, error)

type sourceSubscriber struct {
	format     converter.Converter
	logger     *zap.Logger
	subscriber nc.Subscriber
	streams    []stream.Streamer
}

// Open creates a Sourcer that subscribes via the given notificationcenter factory.
func Open(ctx context.Context, url string, subFnk getSubscriberFnk, options ...Option) (eventstream.Sourcer, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	subscriber, err := subFnk(ctx, url)
	if err != nil {
		return nil, err
	}
	subscriberWrapper := &sourceSubscriber{
		subscriber: subscriber,
		logger:     opts.getLogger(),
		format:     opts.getFormat(),
	}
	if err := subscriber.Subscribe(ctx, subscriberWrapper); err != nil {
		return nil, err
	}
	return subscriberWrapper, nil
}

// Subscribe new stream object
func (s *sourceSubscriber) Subscribe(ctx context.Context, streams ...stream.Streamer) error {
	if len(streams) == 0 {
		return errInvalidStreamObject
	}
	for _, strm := range streams {
		for _, st := range s.streams {
			if st.ID() == strm.ID() {
				return errors.Wrap(errStreamAlreadyRegistered, st.ID())
			}
		}
	}
	s.streams = append(s.streams, streams...)
	return nil
}

// Receive notification message
func (s *sourceSubscriber) Receive(msg nc.Message) error {
	decoded, err := eventstream.MessageDecode(msg.Body(), s.format)
	if err != nil {
		s.logger.Error(`messgage decode`, zap.Error(err))
		return err
	}
	ctx := context.Background()
	for _, strm := range s.streams {
		if strm.Check(ctx, decoded) {
			if err = strm.Put(ctx, decoded); err != nil {
				return err
			}
		}
	}
	return msg.Ack()
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
