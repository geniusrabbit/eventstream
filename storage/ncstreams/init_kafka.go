// +build kafka allstorage all

package ncstreams

import (
	"github.com/geniusrabbit/eventstream/storage"
)

func init() {
	storage.RegisterConnector(connector, "kafka")
}
