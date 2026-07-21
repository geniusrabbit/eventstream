//go:build nats || allsource || all
// +build nats allsource all

package source

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	esource "github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/eventstream/source/nats"
)

func init() {
	Register("nats", func(ctx context.Context, p Params) (eventstream.Sourcer, error) {
		var opts []esource.Option
		if p.Format != "" {
			opts = append(opts, esource.WithFormat(p.Format))
		}
		return nats.Open(ctx, p.Connect, opts...)
	})
}
