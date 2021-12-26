package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
)

// Stream for the redis key value
type Stream struct {
	id         string
	key        *patternkey.PatternKey
	expiration time.Duration
	incremetor bool
	cli        redis.UniversalClient
}

// ID returns unical stream identificator
func (s *Stream) ID() string { return s.id }

// Put message to the stream to process information
func (s *Stream) Put(ctx context.Context, msg message.Message) error {
	if s.incremetor {
		return s.cli.Incr(ctx, s.key.Prepare(msg)).Err()
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.cli.Set(ctx, s.key.Prepare(msg), data, s.expiration).Err()
}

// Check if message suits for the stream
func (s *Stream) Check(ctx context.Context, msg message.Message) bool { return true }

// Run processing loop
func (s *Stream) Run(ctx context.Context) error { return nil }

// Close stream
func (s *Stream) Close() error { return nil }
