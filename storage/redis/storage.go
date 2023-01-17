package redis

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/go-redis/redis/v8"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/internal/patternkey"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/stream"
)

type redisStreamConfig struct {
	Key         string        `json:"key"`
	Expiration  time.Duration `json:"expiration"`
	Incrementor bool          `json:"incrementor"`
}

// Storage for the redis
type Storage struct {
	cli redis.UniversalClient
}

// NewStorage from connection URL
func NewStorage(ctx context.Context, connectUrl string, options ...storage.Option) (*Storage, error) {
	var (
		err  error
		cli  redis.UniversalClient
		conf storage.Config
	)
	for _, o := range options {
		o(&conf)
	}
	if cli, err = connectRedis(ctx, connectUrl); err != nil {
		return nil, err
	}
	return &Storage{cli: cli}, nil
}

// Stream vertica processor
func (st *Storage) Stream(options ...any) (eventstream.Streamer, error) {
	var (
		err        error
		conf       stream.Config
		strmConf   redisStreamConfig
		metricExec metrics.Metricer
	)
	for _, opt := range options {
		switch o := opt.(type) {
		case stream.Option:
			o(&conf)
		case *stream.Config:
			conf = *o
		default:
			stream.WithObjectConfig(o)(&conf)
		}
	}
	if err = conf.Decode(&strmConf); err != nil {
		return nil, err
	}
	if metricExec, err = conf.Metrics.Metric(); err != nil {
		return nil, err
	}
	return eventstream.NewStreamWrapper(&Stream{
		cli:        st.cli,
		id:         conf.Name,
		expiration: strmConf.Expiration,
		incremetor: strmConf.Incrementor,
		key:        patternkey.PatternKeyFromTemplate(strmConf.Key),
	}, conf.Where, metricExec)
}

// Close vertica connection
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
