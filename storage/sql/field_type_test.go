package sql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFieldTypeCastExt(t *testing.T) {
	now := time.Now()
	values := map[string]any{
		"id":        int64(1),
		"date":      now.Format("2006/01/02 15:04:05"),
		"timestamp": now.UnixNano(),
		"text":      []byte("text"),
		"ip":        "127.0.0.2",
		"active":    "1",
		"props":     []int{1, 2, 3, 4},
	}

	assert.Equal(t, int64(1), FieldTypeInt.CastExt(values["id"], 0, ""))
	assert.Equal(t, int8(1), FieldTypeInt8.CastExt(values["id"], 0, ""))
	assert.Equal(t, int32(1), FieldTypeInt32.CastExt(values["id"], 0, ""))
	assert.Equal(t, int64(1), FieldTypeInt64.CastExt(values["id"], 0, ""))
	assert.Equal(t, uint64(1), FieldTypeUint.CastExt(values["id"], 0, ""))
	assert.Equal(t, uint8(1), FieldTypeUint8.CastExt(values["id"], 0, ""))
	assert.Equal(t, uint32(1), FieldTypeUint32.CastExt(values["id"], 0, ""))
	assert.Equal(t, uint64(1), FieldTypeUint64.CastExt(values["id"], 0, ""))

	assert.Equal(t, now.Format("2006-01-02 15:04:05"),
		FieldTypeDate.CastExt(values["date"], 0, "2006-01-02 15:04:05"))
	assert.Equal(t, normTimeHours(now.Unix()),
		normTimeHours(FieldTypeUnixnano.CastExt(values["timestamp"], 0, "").(time.Time).Unix()))

	assert.Equal(t, "text", FieldTypeString.CastExt(values["text"], 0, ""))
	assert.Equal(t, "tex", FieldTypeString.CastExt(values["text"], 3, ""))
	assert.Equal(t, "text ", FieldTypeString.CastExt(values["text"], 5, ""))

	assert.Equal(t, []byte("tex"), FieldTypeFixed.CastExt(values["text"], 3, ""))
	assert.Equal(t, []byte{'t', 'e', 'x', 't', 0}, FieldTypeFixed.CastExt(values["text"], 5, ""))

	assert.Equal(t, "\\000\\000\\000\\000\\000\\000\\000\\000\\000\\000\\377\\377\\177\\000\\000\\002",
		FieldTypeIP.CastExt(values["ip"], 0, "binarystring"))

	assert.Equal(t, true, FieldTypeBoolean.CastExt(values["active"], 0, ""))

	assert.Equal(t, []int32{1, 2, 3, 4}, FieldTypeArrayInt32.CastExt(values["props"], 0, ""))
	assert.Equal(t, []int64{1, 2, 3, 4}, FieldTypeArrayInt64.CastExt(values["props"], 0, ""))
}

func TestFieldTypeArrayString(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, FieldTypeArrayString.CastExt([]any{"a", "b"}, 0, ""))
	assert.Equal(t, []string{"a", "b"}, FieldTypeArrayString.CastExt([]string{"a", "b"}, 0, ""))
	assert.Equal(t, []string{}, FieldTypeArrayString.CastExt(nil, 0, ""))
	assert.NotEqual(t, []string{"sss", "ddd"}, FieldTypeArrayString.Cast("sss,ddd"))

	values, fields, inserts := PrepareFieldsByArray([]string{"tags=t:[]string"})
	assert.Equal(t, []Value{{
		Key:       "t",
		TargetKey: "tags",
		Type:      FieldTypeArrayString,
	}}, values)
	assert.Equal(t, []string{"tags"}, fields)
	assert.Equal(t, []string{"?"}, inserts)

	got, _, _, err := ConvertObjectIntoQueryParams(struct {
		Tags []string `field:"tags"`
	}{})
	assert.NoError(t, err)
	assert.Equal(t, FieldTypeArrayString, got[0].Type)
	assert.Equal(t, "tags", got[0].Key)
}

func normTimeHours(t int64) int64 {
	return t - t%3600
}
