//
// @project geniusrabbit::eventstream 2017, 2020 - 2023
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020 - 2023
//

package ncstreams

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/utils/converter"
)

var (
	errInvalidStreamObject     = errors.New("invalid stream object")
	errStreamAlreadyRegistered = errors.New("stream already registered")
)

type getSubscriberFnk func(ctx context.Context, url string) (nc.Subscriber, error)

type sourceSubscriber struct {
	debug      bool
	format     converter.Converter
	logger     *zap.Logger
	subscriber nc.Subscriber
	streams    []eventstream.Streamer
}

// Open new source by URLs
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
		debug:      opts.Debug,
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
func (s *sourceSubscriber) Subscribe(ctx context.Context, streams ...eventstream.Streamer) error {
	if len(streams) == 0 {
		return errInvalidStreamObject
	}
	for _, stream := range streams {
		for _, st := range s.streams {
			if st.ID() == stream.ID() {
				return errors.Wrap(errStreamAlreadyRegistered, st.ID())
			}
		}
	}
	s.streams = append(s.streams, streams...)
	return nil
}

// Receive notification message
func (s *sourceSubscriber) Receive(message nc.Message) error {
	msg, err := eventstream.MessageDecode(message.Body(), s.format)
	if err != nil {
		s.logger.Error(`messgage decode`, zap.Error(err))
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
