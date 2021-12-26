//go:build nats || allstorage || all
// +build nats allstorage all

// Package nats contains ints stream implementation
//
// @project geniusrabbit::eventstream 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//
package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/nats"
)

func natsNewPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return nats.NewPublisher(nats.WithNatsURL(url))
}

// OpenNATS publisher connectior
func OpenNATS(ctx context.Context, url string, options ...storage.Option) (eventstream.Storager, error) {
	return Open(ctx, url, natsNewPublisher, options...)
}

func init() {
	storage.RegisterConnector("nats", func(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
		return OpenNATS(ctx, conf.Connect, storage.WithDebug(conf.Debug))
	})
}
