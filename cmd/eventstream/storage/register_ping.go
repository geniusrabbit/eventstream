//go:build ping || http || allsource || all
// +build ping http allsource all

package storage

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage/ping"
)

func init() {
	Register("ping", openPing)
	Register("http", openPing)
	RegisterStream("ping", pingStreamOpts)
	RegisterStream("http", pingStreamOpts)
}

func openPing(ctx context.Context, p Params) (eventstream.Storager, error) {
	opts := []ping.Option{}
	if p.Method != "" {
		opts = append(opts, ping.WithMethod(p.Method))
	}
	connect := p.Connect
	if p.URL != "" {
		opts = append(opts, ping.WithURL(p.URL))
		if connect == "" {
			connect = p.URL
		}
	}
	return ping.Open(ctx, connect, opts...)
}

func pingStreamOpts(streamName string, p StreamParams) ([]any, error) {
	opts := []any{ping.WithStreamID(streamName)}
	if p.URL != "" {
		opts = append(opts, ping.WithStreamURL(p.URL))
	}
	if p.Method != "" {
		opts = append(opts, ping.WithStreamMethod(p.Method))
	}
	if p.ContentType != "" {
		opts = append(opts, ping.WithContentType(p.ContentType))
	}
	return opts, nil
}
