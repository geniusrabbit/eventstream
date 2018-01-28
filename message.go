//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/twinj/uuid"
)

// Errors set
var (
	ErrUndefinedDataType       = errors.New("Undefined data types")
	ErrInvalidMessageFieldType = errors.New("Invalid message field type")
)

// Message object
type Message map[string]interface{}

// MessageDecode from bytes
func MessageDecode(item interface{}, converter converter.Converter) (msg Message, err error) {
	var data []byte
	switch v := item.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil, ErrUndefinedDataType // Undefined type
	}

	err = converter.Unmarshal(data, &msg)
	return msg, err
}

// JSON data string
func (m Message) JSON() string {
	data, _ := json.Marshal(m)
	return string(data)
}

// Item by key name
func (m Message) Item(key string, def interface{}) interface{} {
	if v, ok := m[key]; ok && nil != v {
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

// ItemCast value
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
		switch vv := v.(type) {
		case []byte:
			if len(vv) > 16 {
				v, _ = uuid.Parse(string(vv))
			} else {
				v = bytesSize(vv, 16)
			}
		default:
			if v == nil {
				v = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
			} else {
				if v, _ = uuid.Parse(gocast.ToString(v)); v != nil {
					v = v.(uuid.Uuid).Bytes()
				}
			}
		}
		if v != nil && format == "escape" {
			return escapeBytes(v.([]byte), 0)
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
		var ip = net.ParseIP(gocast.ToString(v))
		switch format {
		case "binarystring":
			v = ip2EscapeString(ip)
		default:
			v = ip
		}
	case FieldTypeDate:
		var tm time.Time
		if v != nil {
			switch v.(type) {
			case int64, uint64, float64:
				tm = time.Unix(gocast.ToInt64(v), 0)
			default:
				tm, _ = parseTime(gocast.ToString(v))
			}
		}

		if "" != format {
			return tm.Format(format)
		}
		v = tm
	case FieldTypeUnixnano:
		var tm time.Time
		if v != nil { 
			tm = time.Unix(0, gocast.ToInt64(v))
		}

		if "" != format {
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

// Map value
func (m Message) Map() map[string]interface{} {
	return map[string]interface{}(m)
}
