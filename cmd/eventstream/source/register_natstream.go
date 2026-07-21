//go:build natstream || allsource || all
// +build natstream allsource all

package source

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	esource "github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/eventstream/source/natstream"
)

func init() {
	Register("natstream", func(ctx context.Context, p Params) (eventstream.Sourcer, error) {
		var opts []esource.Option
		if p.Format != "" {
			opts = append(opts, esource.WithFormat(p.Format))
		}
		return natstream.Open(ctx, p.Connect, opts...)
	})
}
