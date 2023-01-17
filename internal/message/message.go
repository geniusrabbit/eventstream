//
// @project geniusrabbit::eventstream 2017, 2019-2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2020
//

package message

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
)

// Errors set
var (
	ErrUndefinedDataType       = errors.New("undefined message data types")
	ErrInvalidMessageFieldType = errors.New("invalid message field type")
)

type unmarshalel interface {
	Unmarshal(data []byte, v any) error
}

// Message object
type Message map[string]any

// MessageDecode from bytes
func MessageDecode(data []byte, converter unmarshalel) (msg Message, err error) {
	err = converter.Unmarshal(data, &msg)
	return msg, err
}

// JSON data string
func (m Message) JSON() string {
	if m == nil {
		return "null"
	}
	data, _ := json.Marshal(m)
	return string(data)
}

// Item returns the value by key name or default
func (m Message) Item(key string, def any) any {
	if v, ok := m[key]; ok && v != nil {
		return v
	}
	return def
}

// String item value
func (m Message) String(key, def string) string {
	switch v := m.Item(key, def).(type) {
	case float64:
		return strconv.FormatFloat(v, 'G', 6, 64)
	default:
		return gocast.Str(v)
	}
}

// Map returns the message as map[string]any
func (m Message) Map() map[string]any {
	return map[string]any(m)
}

// ItemCast converts any key value into the field_type
func (m Message) ItemCast(key string, t FieldType, length int, format string) any {
	newVal := t.Cast(m.Item(key, nil))
	if newVal == nil {
		return newVal
	}
	switch t {
	case FieldTypeString:
		if length > 0 {
			s := newVal.(string)
			if length < len(s) {
				return s[:length]
			} else if length != len(s) {
				return s + strings.Repeat(" ", length-len(s))
			}
		}
	case FieldTypeFixed:
		res := bytesSize(newVal.([]byte), length)
		if format == "escape" {
			return escapeBytes(res, 0)
		}
		newVal = res
	case FieldTypeUUID:
		if format == "escape" {
			return escapeBytes(newVal.([]byte), 0)
		}
	case FieldTypeIP:
		ip := newVal.(net.IP)
		switch format {
		case "binarystring":
			newVal = ip2EscapeString(ip)
		case "fix":
			newVal = bytesSize(ip, 16)
		}
	case FieldTypeDate, FieldTypeUnixnano:
		if format != "" {
			return newVal.(time.Time).Format(format)
		}
	}
	return newVal
}
