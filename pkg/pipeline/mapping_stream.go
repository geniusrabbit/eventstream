package pipeline

import (
	"context"

	"github.com/geniusrabbit/eventstream/pkg/message"
	"github.com/geniusrabbit/eventstream/stream"
)

type mappingStreamer struct {
	id     string
	mapper Mapper
	sink   stream.Streamer
}

func wrapMapper(sink stream.Streamer, mapper Mapper) stream.Streamer {
	if mapper == nil {
		return sink
	}
	return &mappingStreamer{id: sink.ID(), mapper: mapper, sink: sink}
}

func (s *mappingStreamer) ID() string { return s.id }

func (s *mappingStreamer) Put(ctx context.Context, msg message.Message) error {
	outs, err := s.mapper.Map(ctx, msg)
	if err != nil {
		return err
	}
	for _, out := range outs {
		if err := s.sink.Put(ctx, out); err != nil {
			return err
		}
	}
	return nil
}

func (s *mappingStreamer) Check(ctx context.Context, msg message.Message) bool {
	return s.sink.Check(ctx, msg)
}

func (s *mappingStreamer) Run(ctx context.Context) error { return s.sink.Run(ctx) }

func (s *mappingStreamer) Close() error { return s.sink.Close() }
