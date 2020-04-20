// +build nats allsource all

package ncstreams

import (
	"github.com/geniusrabbit/eventstream/source"
)

func init() {
	source.RegisterConnector(connector, "nats")
}
