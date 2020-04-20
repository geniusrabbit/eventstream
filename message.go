//
// @project geniusrabbit::eventstream 2017, 2019-2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019-2020
//

package eventstream

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/demdxx/gocast"
	"github.com/myesui/uuid"
)

// Errors set
var (
	ErrUndefinedDataType       = errors.New("[eventstream::message] undefined data types")
	ErrInvalidMessageFieldType = errors.New("[eventstream::message] invalid message field type")
)

type unmarshalel interface {
	Unmarshal(data []byte, v interface{}) error
}

// Message object
type Message map[string]interface{}

// MessageDecode from bytes
func MessageDecode(data []byte, converter unmarshalel) (msg Message, err error) {
	err = converter.Unmarshal(data, &msg)
	return msg, err
}

// JSON data string
func (m Message) JSON() string {
	data, _ := json.Marshal(m)
	return string(data)
}

// Item returns the value by key name or default
func (m Message) Item(key string, def interface{}) interface{} {
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
	return gocast.ToString(m.Item(key, def))
}

// ItemCast converts any key value into the field_type
func (m Message) ItemCast(key string, t FieldType, length int, format string) (v interface{}) {
	v = m.Item(key, nil)
	switch t {
	case FieldTypeString:
		return gocast.ToString(v)
	case FieldTypeFixed:
		switch vv := v.(type) {
		case []byte:
			v = bytesSize(vv, length)
		default:
			v = bytesSize([]byte(gocast.ToString(v)), length)
		}
		if format == "escape" {
			return escapeBytes(v.([]byte), 0)
		}
	case FieldTypeUUID:
		res := valueToUUIDBytes(v)
		if res != nil && format == "escape" {
			return escapeBytes(res, 0)
		}
	case FieldTypeInt:
		return gocast.ToInt64(v)
	case FieldTypeInt32:
		return gocast.ToInt32(v)
	case FieldTypeInt8:
		return int8(gocast.ToInt(v))
	case FieldTypeUint:
		return gocast.ToUint64(v)
	case FieldTypeUint32:
		return gocast.ToUint32(v)
	case FieldTypeUint8:
		return uint8(gocast.ToUint(v))
	case FieldTypeFloat:
		return gocast.ToFloat64(v)
	case FieldTypeBoolean:
		return gocast.ToBool(v)
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
			gocast.ToSlice(arr, v, "")
			v = arr
		} else {
			v = []int32{}
		}
	case FieldTypeArrayInt64:
		if v != nil {
			var arr = []int64{}
			gocast.ToSlice(arr, v, "")
			v = arr
		} else {
			v = []int64{}
		}
	}
	return v
}

// Map returns the message as map[string]interface{}
func (m Message) Map() map[string]interface{} {
	return map[string]interface{}(m)
}

func valueToTime(v interface{}) (tm time.Time) {
	switch vl := v.(type) {
	case nil:
	case int64:
		tm = time.Unix(vl, 0)
	case uint64:
		tm = time.Unix(int64(vl), 0)
	case float64:
		tm = time.Unix(int64(vl), 0)
	case string:
		tm, _ = parseTime(gocast.ToString(v))
	default:
		tm, _ = parseTime(gocast.ToString(v))
	}
	return tm
}

func valueUnixNanoToTime(v interface{}) (tm time.Time) {
	switch vl := v.(type) {
	case nil:
	case int64:
		tm = time.Unix(0, vl)
	case uint64:
		tm = time.Unix(0, int64(vl))
	case float64:
		tm = time.Unix(0, int64(vl))
	case string:
		tm, _ = parseTime(gocast.ToString(v))
	default:
		tm, _ = parseTime(gocast.ToString(v))
	}
	return tm
}

func valueToIP(v interface{}) (ip net.IP) {
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
		ip = net.ParseIP(gocast.ToString(v))
	}
	return ip
}

func valueToUUIDBytes(v interface{}) (res []byte) {
	switch vv := v.(type) {
	case []byte:
		if len(vv) > 16 {
			if _uuid, _ := uuid.Parse(string(vv)); _uuid != nil {
				v = _uuid.Bytes()
			}
		} else {
			res = bytesSize(vv, 16)
		}
	default:
		if v == nil {
			res = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		} else {
			if v, _ = uuid.Parse(gocast.ToString(v)); v != nil {
				res = v.(*uuid.UUID).Bytes()
			}
		}
	}
	return res
}
