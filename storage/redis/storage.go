// Package redis provides a Redis key/value Storager (SET with optional pattern keys).
package redis

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream/internal/patternkey"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/stream"
)

// Storage is a Redis connection that opens key/value write streams.
type Storage struct {
	cli redis.UniversalClient
}

// NewStorage connects to Redis using the given connection URL.
func NewStorage(ctx context.Context, connectURL string) (*Storage, error) {
	cli, err := connectRedis(ctx, connectURL)
	if err != nil {
		return nil, err
	}
	return &Storage{cli: cli}, nil
}

// Stream redis processor
func (st *Storage) Stream(options ...any) (stream.Streamer, error) {
	var conf streamOptions
	for _, opt := range options {
		switch o := opt.(type) {
		case StreamOption:
			o(&conf)
		default:
			return nil, errors.Wrapf(storage.ErrInvalidOption, `redis stream: %T`, opt)
		}
	}
	if conf.key == "" {
		return nil, fmt.Errorf("redis: stream key is required")
	}
	return &Stream{
		cli:        st.cli,
		id:         conf.id,
		expiration: conf.expiration,
		incremetor: conf.incrementor,
		key:        patternkey.PatternKeyFromTemplate(conf.key),
	}, nil
}

// Close redis connection
func (st *Storage) Close() error {
	return st.cli.Close()
}

func connectRedis(ctx context.Context, connectURL string) (redis.UniversalClient, error) {
	u, err := url.Parse(connectURL)
	if err != nil {
		return nil, err
	}
	var (
		password, _           = u.User.Password()
		username              = u.User.Username()
		maxConnAge, _         = time.ParseDuration(u.Query().Get("max_idle_conns"))
		poolTimeout, _        = time.ParseDuration(u.Query().Get("pool_timeout"))
		idleTimeout, _        = time.ParseDuration(u.Query().Get("idle_timeout"))
		idleCheckFrequency, _ = time.ParseDuration(u.Query().Get("idle_check_frequency"))
		cli                   = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:              strings.Split(u.Host, ","),
			DB:                 gocast.Number[int](u.Path[1:]),
			Username:           username,
			Password:           password,
			MaxRetries:         defInt(gocast.Number[int](u.Query().Get("retry")), 3),
			PoolSize:           gocast.Number[int](u.Query().Get("pool_size")),
			MinIdleConns:       gocast.Number[int](u.Query().Get("max_idle_conns")),
			MaxConnAge:         maxConnAge,
			PoolTimeout:        poolTimeout,
			IdleTimeout:        idleTimeout,
			IdleCheckFrequency: idleCheckFrequency,
		})
	)
	return cli, nil
}

func defInt(v1, def int) int {
	if v1 == 0 {
		return def
	}
	return v1
}
