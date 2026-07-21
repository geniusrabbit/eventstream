package storage

import (
	"testing"

	"github.com/geniusrabbit/eventstream/pkg/message"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageTemplate(t *testing.T) {
	testMsg1 := message.MapMessage{"type": "error", "msg": "test"}
	testMsg2 := message.MapMessage{"type": "notify", "msg": "test"}
	targetMsg := map[string]any{"category": "processed", "type": "error", "message": "test"}
	tpl, err := newMessageTemplate(map[string]any{
		"category": "processed",
		"type":     "{{type}}",
		"message":  "{{msg}}",
	}, `type=="error"`)

	assert.NoError(t, err)
	assert.NotNil(t, tpl)
	assert.True(t, tpl.check(testMsg1))
	assert.False(t, tpl.check(testMsg2))
	assert.Equal(t, targetMsg, tpl.prepare(testMsg1))
}
