//go:build kafka || allsource || all
// +build kafka allsource all

package source

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	esource "github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/eventstream/source/kafka"
)

func init() {
	Register("kafka", func(ctx context.Context, p Params) (eventstream.Sourcer, error) {
		var opts []esource.Option
		if p.Format != "" {
			opts = append(opts, esource.WithFormat(p.Format))
		}
		return kafka.Open(ctx, p.Connect, opts...)
	})
}
