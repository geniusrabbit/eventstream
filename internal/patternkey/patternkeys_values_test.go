package patternkey

import (
	"testing"

	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/stretchr/testify/assert"
)

func TestTagValues(t *testing.T) {
	tagVals := PatternKeysFrom("{{name}}_{{id}}", "os_{{os}}")
	assert.NotNil(t, tagVals)
	res := tagVals.Prepare(message.Message{
		"name": "testname",
		"id":   1,
		"os":   "mac",
	})
	assert.ElementsMatch(t, []string{"testname_1", "os_mac"}, res)
}
