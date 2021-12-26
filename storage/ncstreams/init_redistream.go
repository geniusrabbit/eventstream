//go:build redistream || allsource || all
// +build redistream allsource all

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/redis"
)

func redisNewPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return redis.NewPublisher(redis.WithRedisURL(url))
}

// OpenRedis publisher connectior
func OpenRedis(ctx context.Context, url string, options ...storage.Option) (eventstream.Storager, error) {
	return Open(ctx, url, redisNewPublisher, options...)
}

func init() {
	storage.RegisterConnector("redistream", func(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
		return OpenRedis(ctx, conf.Connect, storage.WithDebug(conf.Debug))
	})
}
