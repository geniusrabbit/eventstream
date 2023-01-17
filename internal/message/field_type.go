//
// @project geniusrabbit::eventstream 2017 - 2019, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019, 2022
//

package message

import (
	"github.com/demdxx/gocast/v2"
)

var typeList = []string{
	"string",
	"fix", // String
	"uuid",
	"int",
	"int8",
	"int32",
	"uint",
	"uint8",
	"uint32",
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
	FieldTypeInt8
	FieldTypeInt32
	FieldTypeUint
	FieldTypeUint8
	FieldTypeUint32
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
func (t FieldType) Cast(v any) any {
	switch t {
	case FieldTypeString:
		return gocast.Str(v)
	case FieldTypeFixed:
		switch vv := v.(type) {
		case []byte:
			v = vv
		default:
			v = []byte(gocast.Str(v))
		}
	case FieldTypeUUID:
		v = valueToUUIDBytes(v)
	case FieldTypeInt:
		return gocast.Number[int64](v)
	case FieldTypeInt8:
		return gocast.Number[int8](v)
	case FieldTypeInt32:
		return gocast.Number[int32](v)
	case FieldTypeUint:
		return gocast.Number[uint64](v)
	case FieldTypeUint8:
		return gocast.Number[uint8](v)
	case FieldTypeUint32:
		return gocast.Number[uint32](v)
	case FieldTypeFloat:
		return gocast.Number[float64](v)
	case FieldTypeBoolean:
		return gocast.Bool(v)
	case FieldTypeIP:
		v = valueToIP(v)
	case FieldTypeDate:
		v = valueToTime(v)
	case FieldTypeUnixnano:
		v = valueUnixNanoToTime(v)
	case FieldTypeArrayInt32:
		v = gocast.Cast[[]int32](v)
	case FieldTypeArrayInt64:
		v = gocast.Cast[[]int64](v)
	}
	return v
}
