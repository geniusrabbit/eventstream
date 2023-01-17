package ncstreams

import (
	"strings"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/eventstream"
	"github.com/tonalfitness/govaluate/v3"
)

type messageTemplate struct {
	// Fields additional data field -> value for target message
	fields map[string]any

	// Mapping from one message name to other
	mapping map[string]string

	// WhereCondition of stream
	condition *govaluate.EvaluableExpression
}

func newMessageTemplate(fields map[string]any, where string) (tpl *messageTemplate, err error) {
	tpl = &messageTemplate{}
	if where != `` {
		if tpl.condition, err = govaluate.NewEvaluableExpression(where); err != nil {
			return nil, err
		}
	}
	if fields == nil {
		return tpl, err
	}
	tpl.fields = map[string]any{}
	tpl.mapping = map[string]string{}
	for key, value := range fields {
		s := gocast.Str(value)
		if strings.HasPrefix(s, `{{`) && strings.HasSuffix(s, `}}`) {
			tpl.mapping[key] = s[2 : len(s)-2]
		} else {
			tpl.fields[key] = s
		}
	}
	if len(tpl.fields) == 0 {
		tpl.fields = nil
	}
	if len(tpl.mapping) == 0 {
		tpl.mapping = nil
	}
	return tpl, nil
}

func (t *messageTemplate) check(msg eventstream.Message) bool {
	if t.condition == nil {
		return true
	}
	res, _ := t.condition.Evaluate(msg.Map())
	return gocast.Bool(res)
}

func (t *messageTemplate) prepare(msg eventstream.Message) map[string]any {
	data := msg.Map()
	if t.mapping == nil && t.fields == nil {
		return data
	}
	newData := map[string]any{}
	if t.mapping != nil {
		for key, target := range t.mapping {
			newData[key] = data[target]
		}
	}
	if t.fields != nil {
		for key, value := range t.fields {
			newData[key] = value
		}
	}
	return newData
}
