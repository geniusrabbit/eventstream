// +build vertica allsource all

package vertica

import (
	"github.com/geniusrabbit/eventstream/storage"
)

func init() {
	storage.RegisterConnector(connector, "vertica")
}
