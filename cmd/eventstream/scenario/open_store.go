package scenario

import (
	"context"
	"fmt"

	"github.com/geniusrabbit/eventstream"
	cmdstorage "github.com/geniusrabbit/eventstream/cmd/eventstream/storage"
)

func openStore(ctx context.Context, name string, item Item, debug bool) (eventstream.Storager, error) {
	driver := resolveDriver(item)
	connect := itemString(item, "connect")
	if driver == "" {
		return nil, fmt.Errorf("store %q: empty driver", name)
	}
	if connect == "" && driver != "ping" && driver != "http" {
		return nil, fmt.Errorf("store %q: empty connect", name)
	}
	return cmdstorage.Open(ctx, driver, cmdstorage.Params{
		Connect:   connect,
		Debug:     debug,
		InitQuery: itemStringSlice(item, "init_query"),
		Method:    itemString(item, "method"),
		URL:       itemString(item, "url"),
	})
}
