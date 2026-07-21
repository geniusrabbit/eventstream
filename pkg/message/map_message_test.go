package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapMessageAccessors(t *testing.T) {
	msg := MapMessage{
		"id":   int64(1),
		"text": "hello",
		"n":    float64(3.14),
	}

	assert.Equal(t, int64(1), msg.Item("id", 0))
	assert.Equal(t, "missing", msg.Item("nope", "missing"))
	assert.Equal(t, "hello", msg.Str("text", ""))
	assert.Equal(t, "3.14", msg.Str("n", ""))
	assert.Contains(t, msg.JSON(), `"text":"hello"`)
	assert.Equal(t, map[string]any(msg), msg.Map())

	v, err := msg.Get("id")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), v)

	_, err = msg.Get("missing")
	assert.ErrorIs(t, err, ErrFieldNotFound)
}
