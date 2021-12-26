//
// @project geniusrabbit::eventstream 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package message

import "github.com/demdxx/gocast"

var typeList = []string{
	"string",
	"fix", // String
	"uuid",
	"int",
	"int32",
	"int8",
	"uint",
	"uint32",
	"uint8",
	"float",
	"bool",
	"ip",
	"date",
	"unixnano",
	"[]int32",
	"[]int64",
}

// Field scalar types enum
const (
	FieldTypeString FieldType = iota
	FieldTypeFixed
	FieldTypeUUID
	FieldTypeInt
	FieldTypeInt32
	FieldTypeInt8
	FieldTypeUint
	FieldTypeUint32
	FieldTypeUint8
	FieldTypeFloat
	FieldTypeBoolean
	FieldTypeIP
	FieldTypeDate
	FieldTypeUnixnano
	FieldTypeArrayInt32
	FieldTypeArrayInt64
)

// FieldType of data represents scalar types supported
// by eventstream message processing
type FieldType int

// TypeByString name
func TypeByString(t string) FieldType {
	for i, s := range typeList {
		if s == t {
			return FieldType(i)
		}
	}
	return FieldTypeString
}

// String implementaion of fmt.Stringer
func (t FieldType) String() string {
	if t > 0 && int(t) < len(typeList) {
		return typeList[t]
	}
	return typeList[0]
}

// Cast value into the fieldType
func (t FieldType) Cast(v interface{}) interface{} {
	switch t {
	case FieldTypeString:
		return gocast.ToString(v)
	case FieldTypeFixed:
		switch vv := v.(type) {
		case []byte:
			v = vv
		default:
			v = []byte(gocast.ToString(v))
		}
	case FieldTypeUUID:
		v = valueToUUIDBytes(v)
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
		v = valueToIP(v)
	case FieldTypeDate:
		v = valueToTime(v)
	case FieldTypeUnixnano:
		v = valueUnixNanoToTime(v)
	case FieldTypeArrayInt32:
		if v != nil {
			var arr = []int32{}
			_ = gocast.ToSlice(arr, v, "")
			v = arr
		} else {
			v = []int32{}
		}
	case FieldTypeArrayInt64:
		if v != nil {
			var arr = []int64{}
			_ = gocast.ToSlice(arr, v, "")
			v = arr
		} else {
			v = []int64{}
		}
	}
	return v
}
