// +build kafka allsource all

package kafka

import (
	"github.com/geniusrabbit/eventstream/source"
)

func init() {
	source.RegisterConnector(connector, "kafka")
}
