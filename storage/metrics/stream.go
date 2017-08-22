//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

import (
	"errors"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/metrics"
)

var (
	errInvalidMetricsItemConfig = errors.New("Invalid metrics item config")
)

type stream struct {
	prefix  string
	metrics []*metricItem
	metrica notificationcenter.Logger
}

func newStream(metrica notificationcenter.Logger, conf eventstream.ConfigItem) (*stream, error) {
	var metrics []*metricItem

	switch mts := conf.Item("metrics", nil).(type) {
	case []interface{}:
		for _, mit := range mts {
			switch mp := mit.(type) {
			case map[string]interface{}:
				if name, ok := mp["name"]; ok {
					var (
						tp, _   = mp["type"]
						tags, _ = mp["tags"]
						vl, _   = mp["value"]
						mtags   map[string]string
					)

					if tags != nil {
						switch ntags := tags.(type) {
						case []interface{}:
							if len(ntags) > 0 {
								mtags, _ = gocast.ToStringMap(ntags[0], "", false)
							}
						case []map[string]interface{}:
							if len(ntags) > 0 {
								mtags, _ = gocast.ToStringMap(ntags[0], "", false)
							}
						default:
							mtags, _ = gocast.ToStringMap(ntags, "", false)
						}
					}

					item := &metricItem{
						Name:  gocast.ToString(name),
						Type:  gocast.ToString(tp),
						Tags:  mtags,
						Value: gocast.ToString(vl),
					}
					item.updateParams()
					metrics = append(metrics, item)
				}
			default:
				return nil, errInvalidMetricsItemConfig
			}
		}
	default:
		return nil, errInvalidMetricsItemConfig
	}

	return &stream{
		prefix:  conf.String("prefix", ""),
		metrics: metrics,
		metrica: metrica,
	}, nil
}

// Put message to stream
func (s *stream) Put(msg eventstream.Message) error {
	return s.metrica.Send(s.prepareMetricsMessage(msg)...)
}

// Close implementation
func (s *stream) Close() error {
	return nil
}

// Run loop
func (s *stream) Run() error {
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (s *stream) prepareMetricsMessage(msg eventstream.Message) (result []interface{}) {
	for _, mt := range s.metrics {
		if replacer := mt.replacer(msg); replacer == nil {
			result = append(result, metrics.Message{
				Name:  mt.Name,
				Type:  mt.getType(),
				Tags:  mt.Tags,
				Value: msg.Item(mt.Value, nil),
			})
		} else {
			result = append(result, metrics.Message{
				Name:  replacer.Replace(mt.Name),
				Type:  mt.getType(),
				Tags:  mt.getTags(replacer),
				Value: msg.Item(mt.Value, nil),
			})
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func gIndexOfStr(s string, arr [][2]string) int {
	for i, sv := range arr {
		if sv[1] == s {
			return i
		}
	}
	return -1
}
