// +build kafka allsource all

package ncstreams

import (
	"github.com/geniusrabbit/eventstream/source"
)

func init() {
	source.RegisterConnector(connector, "kafka")
}
