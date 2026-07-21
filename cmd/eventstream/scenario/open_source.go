package scenario

import (
	"context"
	"fmt"

	"github.com/geniusrabbit/eventstream"
	cmdsource "github.com/geniusrabbit/eventstream/cmd/eventstream/source"
)

func openSource(ctx context.Context, name string, item Item) (eventstream.Sourcer, error) {
	driver := resolveDriver(item)
	connect := itemString(item, "connect")
	if driver == "" {
		return nil, fmt.Errorf("source %q: empty driver", name)
	}
	if connect == "" {
		return nil, fmt.Errorf("source %q: empty connect", name)
	}
	return cmdsource.Open(ctx, driver, cmdsource.Params{
		Connect: connect,
		Format:  itemString(item, "format"),
	})
}
