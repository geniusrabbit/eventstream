package message

import (
	"encoding/json"
	"strconv"

	"github.com/demdxx/gocast/v2"
	"github.com/pkg/errors"
)

// MapMessage is the default Message implementation.
type MapMessage map[string]any

// Decode from bytes using the provided unmarshaler.
func Decode(data []byte, converter Unmarshaler) (msg MapMessage, err error) {
	err = converter.Unmarshal(data, &msg)
	return msg, err
}

// MapMessageDecode is an alias for Decode.
func MapMessageDecode(data []byte, converter Unmarshaler) (MapMessage, error) {
	return Decode(data, converter)
}

// JSON data string
func (m MapMessage) JSON() string {
	if m == nil {
		return "null"
	}
	data, _ := json.Marshal(m)
	return string(data)
}

// Get returns the value by key name
func (m MapMessage) Get(name string) (any, error) {
	if v, ok := m[name]; ok {
		return v, nil
	}
	return nil, errors.Wrap(ErrFieldNotFound, name)
}

// Item returns the value by key name or default
func (m MapMessage) Item(key string, def any) any {
	if v, ok := m[key]; ok && v != nil {
		return v
	}
	return def
}

// Str returns the string value by key name or default
func (m MapMessage) Str(key, def string) string {
	switch v := m.Item(key, def).(type) {
	case float64:
		return strconv.FormatFloat(v, 'G', 6, 64)
	default:
		return gocast.Str(v)
	}
}

// Map returns the message as map[string]any
func (m MapMessage) Map() map[string]any {
	return map[string]any(m)
}

var _ Message = (MapMessage)(nil)
