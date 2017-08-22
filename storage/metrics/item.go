//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

import (
	"regexp"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/notificationcenter/metrics"
)

var (
	paramsSearch = regexp.MustCompile(`\{\{([^}]+)\}\}`)
)

type metricItem struct {
	Name   string
	Type   string
	Tags   map[string]string
	Value  string
	params [][2]string
}

func (item *metricItem) updateParams() {
	item.params = nil
	item.updateParamsByString(item.Name)

	for k, v := range item.Tags {
		item.updateParamsByString(k)
		item.updateParamsByString(v)
	}
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
		return
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
