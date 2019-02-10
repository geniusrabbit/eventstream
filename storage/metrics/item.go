//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package metrics

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/notificationcenter/metrics"
)

var (
	paramsSearch = regexp.MustCompile(`\{\{([^}]+)\}\}`)
)

type metricItem struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Tags   map[string]string `json:"tags"`
	Value  string            `json:"value"`
	params [][2]string
}

func (item *metricItem) updateParams(tags interface{}) (err error) {
	item.params = nil
	item.Tags = map[string]string{}
	item.updateParamsByString(item.Name)

	switch v := tags.(type) {
	case nil:
	case map[string]string:
		item.Tags = v
	case map[string]interface{}:
		item.Tags, err = gocast.ToStringMap(v, "", false)
	case []map[string]interface{}:
		for _, it := range v {
			if it == nil {
				continue
			}

			var tags map[string]string
			if tags, err = gocast.ToStringMap(it, "", false); err != nil {
				break
			}

			if tags != nil {
				for k, v := range tags {
					item.Tags[k] = v
				}
			}
		}
	case []interface{}:
		for _, it := range v {
			if it == nil {
				continue
			}

			var tags map[string]string
			if tags, err = gocast.ToStringMap(it, "", false); err != nil {
				break
			}

			if tags != nil {
				for k, v := range tags {
					item.Tags[k] = v
				}
			}
		}
	case []map[string]string:
		for _, tags := range v {
			if tags == nil {
				continue
			}
			for k, v := range tags {
				item.Tags[k] = v
			}
		}
	default:
		return fmt.Errorf("[metrics] unsupported tags type %T", item.Tags)
	}

	if item.Tags != nil {
		for k, v := range item.Tags {
			item.updateParamsByString(k)
			item.updateParamsByString(v)
		}
	}

	return
}

// UnmarshalJSON data
func (item *metricItem) UnmarshalJSON(data []byte) (err error) {
	var itemVal struct {
		Name  string      `json:"name"`
		Type  string      `json:"type"`
		Tags  interface{} `json:"tags"`
		Value string      `json:"value"`
	}

	if err = json.Unmarshal(data, &itemVal); err == nil {
		item.Name = itemVal.Name
		item.Type = itemVal.Type
		item.Value = itemVal.Value
		err = item.updateParams(itemVal.Tags)
	}

	return err
}

func (item *metricItem) replacer(msg eventstream.Message) *strings.Replacer {
	if len(item.params) < 1 {
		return nil
	}

	var params []string
	for _, param := range item.params {
		params = append(params, param[0], msg.String(param[1], ""))
	}
	return strings.NewReplacer(params...)
}

func (item *metricItem) updateParamsByString(str string) {
	args := paramsSearch.FindAllStringSubmatch(str, -1)
	for _, it := range args {
		s := strings.TrimSpace(it[1])
		if len(s) > 0 && gIndexOfStr(s, item.params) == -1 {
			item.params = append(item.params, [2]string{it[0], s})
		}
	}
}

func (item *metricItem) getTags(replacer *strings.Replacer) (resp map[string]string) {
	if replacer == nil || item.Tags == nil || len(item.Tags) < 1 {
		return item.Tags
	}

	resp = map[string]string{}
	for k, v := range item.Tags {
		resp[replacer.Replace(k)] = replacer.Replace(v)
	}
	return
}

func (item *metricItem) getType() int {
	switch item.Type {
	case "counter", "increment":
		return metrics.MessageTypeIncrement
	case "gauge":
		return metrics.MessageTypeGauge
	case "timing":
		return metrics.MessageTypeTiming
	case "count":
		return metrics.MessageTypeCount
	case "unique":
		return metrics.MessageTypeUnique
	}
	return metrics.MessageTypeUndefined
}
