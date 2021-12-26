//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package eventstream

import (
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/stream/wrapper"
)

// NewStreamWrapper with support condition
func NewStreamWrapper(stream Streamer, where string, metrics metrics.Metricer) (_ Streamer, err error) {
	return wrapper.NewStreamWrapper(stream, where, metrics)
}
