package ncstreams

import (
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

func connector(conf *storage.Config) (eventstream.Storager, error) {
	return Open(conf.Connect, WithDebug(conf.Debug))
}
