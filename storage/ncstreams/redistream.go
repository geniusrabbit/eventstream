//go:build redistream || allsource || all
// +build redistream allsource all

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/redis"
)

func redisNewPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return redis.NewPublisher(redis.WithRedisURL(url))
}

// OpenRedis publisher connectior
func OpenRedis(ctx context.Context, url string, options ...storage.Option) (eventstream.Storager, error) {
	return Open(ctx, url, redisNewPublisher, options...)
}
