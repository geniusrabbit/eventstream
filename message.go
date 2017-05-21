//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

import (
	"errors"
	"time"
	"unicode"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/eventstream/converter"
)

var (
	timeFormats = []string{
		"2006-01-02",
		"01-02-2006",
		time.RFC1123Z,
		time.RFC3339Nano,
		time.UnixDate,
		time.RubyDate,
		time.RFC1123,
		time.RFC3339,
		time.RFC822,
		time.RFC850,
		time.RFC822Z,
	}
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

// Item by key name
func (m Message) Item(key string, def interface{}) interface{} {
	if v, ok := m[key]; ok && nil != v {
		return v
	}
	return def
}

// ItemCast value
func (m Message) ItemCast(key string, t FieldType, format string) (v interface{}) {
	v = m.Item(key, nil)
	switch t {
	case FieldTypeString:
		return gocast.ToString(v)
	case FieldTypeInt:
		return gocast.ToInt64(v)
	case FieldTypeFloat:
		return gocast.ToFloat64(v)
	case FieldTypeBoolean:
		return gocast.ToBool(v)
	case FieldTypeDate:
		var tm time.Time
		if nil != v {
			s := gocast.ToString(v)
			if isInt(s) {
				tm = time.Unix(gocast.ToInt64(v), 0)
			} else {
				tm, _ = parseTime(s)
			}
		}

		if "" != format {
			return tm.Format(format)
		}
		v = tm
	}
	return v
}

///////////////////////////////////////////////////////////////////////////////
/// Field type
///////////////////////////////////////////////////////////////////////////////

var typeList = []string{
	"string",
	"int",
	"float",
	"bool",
	"date",
}

// FieldType data
type FieldType int

// String implementaion of fmt.Stringer
func (t FieldType) String() string {
	if t > 0 && int(t) < len(typeList) {
		return typeList[t]
	}
	return typeList[0]
}

// Types enum
const (
	FieldTypeString FieldType = iota
	FieldTypeInt
	FieldTypeFloat
	FieldTypeBoolean
	FieldTypeDate
)

// TypeByString name
func TypeByString(t string) FieldType {
	for i, s := range typeList {
		if s == t {
			return FieldType(i)
		}
	}
	return FieldTypeString
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

// ParseTime from string
func parseTime(tm string) (t time.Time, err error) {
	for _, f := range timeFormats {
		if t, err = time.Parse(f, tm); nil == err {
			break
		}
	}
	return
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
