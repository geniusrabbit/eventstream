//
// @project geniusrabbit::eventstream 2017, 2019-2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2020
//

package message

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/google/uuid"
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
	}
	return gocast.Str(m.Item(key, def))
}

// ItemCast converts any key value into the field_type
func (m Message) ItemCast(key string, t FieldType, length int, format string) (v any) {
	v = m.Item(key, nil)
	switch t {
	case FieldTypeString:
		return gocast.Str(v)
	case FieldTypeFixed:
		var res []byte
		switch vv := v.(type) {
		case []byte:
			res = bytesSize(vv, length)
		default:
			res = bytesSize([]byte(gocast.Str(v)), length)
		}
		if format == "escape" {
			return escapeBytes(res, 0)
		}
		v = res
	case FieldTypeUUID:
		res := valueToUUIDBytes(v)
		if res != nil && format == "escape" {
			return escapeBytes(res, 0)
		}
		v = res
	case FieldTypeInt:
		return gocast.Number[int64](v)
	case FieldTypeInt32:
		return gocast.Number[int32](v)
	case FieldTypeInt8:
		return gocast.Number[int8](v)
	case FieldTypeUint:
		return gocast.Number[uint64](v)
	case FieldTypeUint32:
		return gocast.Number[uint32](v)
	case FieldTypeUint8:
		return gocast.Number[uint8](v)
	case FieldTypeFloat:
		return gocast.Number[float64](v)
	case FieldTypeBoolean:
		return gocast.Bool(v)
	case FieldTypeIP:
		ip := valueToIP(v)
		switch format {
		case "binarystring":
			v = ip2EscapeString(ip)
		case "fix":
			v = bytesSize(ip, 16)
		default:
			v = ip
		}
	case FieldTypeDate:
		var tm = valueToTime(v)
		if format != "" {
			return tm.Format(format)
		}
		v = tm
	case FieldTypeUnixnano:
		var tm = valueUnixNanoToTime(v)
		if format != "" {
			return tm.Format(format)
		}
		v = tm
	case FieldTypeArrayInt32:
		if v != nil {
			var arr = []int32{}
			//lint:ignore SA1019 deprecation
			_ = gocast.ToSlice(arr, v, "")
			v = arr
		} else {
			v = []int32{}
		}
	case FieldTypeArrayInt64:
		if v != nil {
			var arr = []int64{}
			//lint:ignore SA1019 deprecation
			_ = gocast.ToSlice(arr, v, "")
			v = arr
		} else {
			v = []int64{}
		}
	}
	return v
}

// Map returns the message as map[string]any
func (m Message) Map() map[string]any {
	return map[string]any(m)
}

func valueToTime(v any) (tm time.Time) {
	switch vl := v.(type) {
	case nil:
	case int64:
		tm = time.Unix(vl, 0)
	case uint64:
		tm = time.Unix(int64(vl), 0)
	case float64:
		tm = time.Unix(int64(vl), 0)
	case string:
		tm, _ = parseTime(gocast.Str(v))
	default:
		tm, _ = parseTime(gocast.Str(v))
	}
	return tm
}

func valueUnixNanoToTime(v any) (tm time.Time) {
	switch vl := v.(type) {
	case nil:
	case int64:
		tm = time.Unix(0, vl)
	case uint64:
		tm = time.Unix(0, int64(vl))
	case float64:
		tm = time.Unix(0, int64(vl))
	case string:
		tm, _ = parseTime(gocast.Str(v))
	default:
		tm, _ = parseTime(gocast.Str(v))
	}
	return tm
}

func valueToIP(v any) (ip net.IP) {
	switch vl := v.(type) {
	case net.IP:
		ip = vl
	case uint:
		ip = make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, uint32(vl))
	case uint32:
		ip = make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, vl)
	case int:
		ip = make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, uint32(vl))
	case int32:
		ip = make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, uint32(vl))
	default:
		ip = net.ParseIP(gocast.Str(v))
	}
	return ip
}

func valueToUUIDBytes(v any) (res []byte) {
	switch vv := v.(type) {
	case []byte:
		if len(vv) > 16 {
			if _uuid, err := uuid.Parse(string(vv)); err == nil {
				res = _uuid[:]
			}
		} else {
			res = bytesSize(vv, 16)
		}
	case string:
		if _uuid, err := uuid.Parse(vv); err == nil {
			res = _uuid[:]
		}
	default:
		if v == nil {
			res = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		} else {
			if _uuid, err := uuid.Parse(gocast.Str(v)); err == nil {
				res = _uuid[:]
			}
		}
	}
	return res
}
