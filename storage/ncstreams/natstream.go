//go:build natstream || allstorage || all
// +build natstream allstorage all

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/natstream"
)

func natstreamNewPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return natstream.NewPublisher(natstream.WithNatsURL(url))
}

// OpenNATStream publisher connectior
func OpenNATStream(ctx context.Context, url string, options ...storage.Option) (eventstream.Storager, error) {
	return Open(ctx, url, natstreamNewPublisher, options...)
}
