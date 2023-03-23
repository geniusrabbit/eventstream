//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package eventstream

import (
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/stream/wrapper"
	"github.com/geniusrabbit/eventstream/utils/condition"
)

// NewStreamWrapper with support condition
func NewStreamWrapper(stream Streamer, where condition.Condition, metrics metrics.Metricer) Streamer {
	return wrapper.NewStreamWrapper(stream, where, metrics)
}
