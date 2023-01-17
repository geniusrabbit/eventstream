package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessageItemCast(t *testing.T) {
	now := time.Now()
	msg := Message{
		"id":   int64(1),
		"date": now.Format("2006/01/02 15:04:05"),
		"text": []byte("text"),
		"ip":   "127.0.0.2",
	}

	assert.Equal(t, int32(1), msg.ItemCast("id", FieldTypeInt32, 0, ""))
	assert.Equal(t, now.Format("2006-01-02 15:04:05"),
		msg.ItemCast("date", FieldTypeDate, 0, "2006-01-02 15:04:05"))
	assert.Equal(t, "text", msg.ItemCast("text", FieldTypeString, 0, ""))
	assert.Equal(t, "tex", msg.ItemCast("text", FieldTypeString, 3, ""))
	assert.Equal(t, "text ", msg.ItemCast("text", FieldTypeString, 5, ""))
	assert.Equal(t, []byte("tex"), msg.ItemCast("text", FieldTypeFixed, 3, ""))
	assert.Equal(t, []byte{'t', 'e', 'x', 't', 0}, msg.ItemCast("text", FieldTypeFixed, 5, ""))
	assert.Equal(t, "\\000\\000\\000\\000\\000\\000\\000\\000\\000\\000\\377\\377\\177\\000\\000\\002",
		msg.ItemCast("ip", FieldTypeIP, 0, "binarystring"))
}
