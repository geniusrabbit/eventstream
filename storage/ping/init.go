//go:build http || ping || allsource || all
// +build http ping allsource all

package ping

import (
	"context"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

type extraConfig struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

// Connector of the driver
func Connector(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
	var (
		extConf extraConfig
		err     = conf.Decode(&extConf)
	)
	if err != nil {
		return nil, err
	}
	return &pinger{
		URL:    gocast.Or(extConf.URL, conf.Connect),
		Method: extConf.Method,
	}, nil
}
