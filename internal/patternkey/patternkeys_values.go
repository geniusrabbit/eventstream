package patternkey

import (
	"regexp"
	"strings"

	"github.com/geniusrabbit/eventstream/internal/message"
)

var metricNames = regexp.MustCompile(`\{\{([^\}]*)\}\}`)

// PatternKey contains single j=key/value template
type PatternKey struct {
	value string
	vars  []string
}

// PatternKeyFromTemplate returns key with variables
func PatternKeyFromTemplate(template string) *PatternKey {
	subNames := metricNames.FindAllStringSubmatch(template, -1)
	names := make([]string, 0, len(subNames))
	for _, name := range subNames {
		names = append(names, name[1])
	}
	return &PatternKey{value: template, vars: names}
}

// Prepare values by message
func (key *PatternKey) Prepare(msg message.Message) string {
	val := key.value
	for _, vr := range key.vars {
		val = strings.ReplaceAll(val, "{{"+vr+"}}", msg.Str(vr, ""))
	}
	return val
}

// PatterKeys data processing
type PatterKeys struct {
	keys []*PatternKey
}

// PatternKeysFrom returns list of pattern keys
func PatternKeysFrom(values ...string) *PatterKeys {
	vals := &PatterKeys{keys: make([]*PatternKey, 0, len(values))}
	for _, val := range values {
		valTmp := PatternKeyFromTemplate(val)
		vals.keys = append(vals.keys, valTmp)
	}
	return vals
}

// Prepare values by message
func (tags *PatterKeys) Prepare(msg message.Message) []string {
	if tags == nil || len(tags.keys) == 0 {
		return nil
	}
	vals := make([]string, 0, len(tags.keys))
	for _, tag := range tags.keys {
		vals = append(vals, tag.Prepare(msg))
	}
	return vals
}
