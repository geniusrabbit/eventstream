//
// @project geniusrabbit::eventstream 2017, 2019-2023
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2023
//

package message

import (
	"encoding/json"
	"strconv"

	"github.com/demdxx/gocast/v2"
	"github.com/pkg/errors"
)

// Map object
type MapMessage map[string]any

// MessageDecode from bytes
func MapMessageDecode(data []byte, converter unmarshalel) (msg MapMessage, err error) {
	err = converter.Unmarshal(data, &msg)
	return msg, err
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

// ItemCast converts any key value into the field_type
func (m MapMessage) ItemCast(key string, fieldType FieldType, length int, format string) any {
	return fieldType.CastExt(m.Item(key, nil), length, format)
}

var _ Message = (MapMessage)(nil)
