//go:build redistream || allsource || all
// +build redistream allsource all

package source

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	esource "github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/eventstream/source/redistream"
)

func init() {
	open := func(ctx context.Context, p Params) (eventstream.Sourcer, error) {
		var opts []esource.Option
		if p.Format != "" {
			opts = append(opts, esource.WithFormat(p.Format))
		}
		return redistream.Open(ctx, p.Connect, opts...)
	}
	Register("redistream", open)
	Register("redis", open)
}
