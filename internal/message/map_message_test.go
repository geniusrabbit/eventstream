package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessageItemCast(t *testing.T) {
	now := time.Now()
	msg := MapMessage{
		"id":        int64(1),
		"date":      now.Format("2006/01/02 15:04:05"),
		"timestamp": now.UnixNano(),
		"text":      []byte("text"),
		"ip":        "127.0.0.2",
		"active":    "1",
		"props":     []int{1, 2, 3, 4},
	}

	assert.Equal(t, int64(1), msg.ItemCast("id", FieldTypeInt, 0, ""))
	assert.Equal(t, int8(1), msg.ItemCast("id", FieldTypeInt8, 0, ""))
	assert.Equal(t, int32(1), msg.ItemCast("id", FieldTypeInt32, 0, ""))
	assert.Equal(t, int64(1), msg.ItemCast("id", FieldTypeInt64, 0, ""))
	assert.Equal(t, uint64(1), msg.ItemCast("id", FieldTypeUint, 0, ""))
	assert.Equal(t, uint8(1), msg.ItemCast("id", FieldTypeUint8, 0, ""))
	assert.Equal(t, uint32(1), msg.ItemCast("id", FieldTypeUint32, 0, ""))
	assert.Equal(t, uint64(1), msg.ItemCast("id", FieldTypeUint64, 0, ""))

	assert.Equal(t, now.Format("2006-01-02 15:04:05"),
		msg.ItemCast("date", FieldTypeDate, 0, "2006-01-02 15:04:05"))
	assert.Equal(t, normTimeHours(now.Unix()),
		normTimeHours(msg.ItemCast("timestamp", FieldTypeUnixnano, 0, "").(time.Time).Unix()))

	assert.Equal(t, "text", msg.ItemCast("text", FieldTypeString, 0, ""))
	assert.Equal(t, "tex", msg.ItemCast("text", FieldTypeString, 3, ""))
	assert.Equal(t, "text ", msg.ItemCast("text", FieldTypeString, 5, ""))

	assert.Equal(t, []byte("tex"), msg.ItemCast("text", FieldTypeFixed, 3, ""))
	assert.Equal(t, []byte{'t', 'e', 'x', 't', 0}, msg.ItemCast("text", FieldTypeFixed, 5, ""))

	assert.Equal(t, "\\000\\000\\000\\000\\000\\000\\000\\000\\000\\000\\377\\377\\177\\000\\000\\002",
		msg.ItemCast("ip", FieldTypeIP, 0, "binarystring"))

	assert.Equal(t, true, msg.ItemCast("active", FieldTypeBoolean, 0, ""))

	assert.Equal(t, []int32{1, 2, 3, 4}, msg.ItemCast("props", FieldTypeArrayInt32, 0, ""))
	assert.Equal(t, []int64{1, 2, 3, 4}, msg.ItemCast("props", FieldTypeArrayInt64, 0, ""))
}

func normTimeHours(t int64) int64 {
	return t - t%3600
}
