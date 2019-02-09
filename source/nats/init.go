// +build nats allsource all

package nats

import (
	"github.com/geniusrabbit/eventstream/source"
)

func init() {
	source.RegisterConnector(connector, "nats")
}
