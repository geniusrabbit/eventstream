//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package eventstream

var typeList = []string{
	"string",
	"fix", // String
	"uuid",
	"int",
	"uint",
	"float",
	"bool",
	"ip",
	"date",
	"unixnano",
	"[]int32",
	"[]int64",
}

// Types enum
const (
	FieldTypeString FieldType = iota
	FieldTypeFixed
	FieldTypeUUID
	FieldTypeInt
	FieldTypeUint
	FieldTypeFloat
	FieldTypeBoolean
	FieldTypeIP
	FieldTypeDate
	FieldTypeUnixnano
	FieldTypeArrayInt32
	FieldTypeArrayInt64
)

// FieldType data
type FieldType int

// String implementaion of fmt.Stringer
func (t FieldType) String() string {
	if t > 0 && int(t) < len(typeList) {
		return typeList[t]
	}
	return typeList[0]
}

// TypeByString name
func TypeByString(t string) FieldType {
	for i, s := range typeList {
		if s == t {
			return FieldType(i)
		}
	}
	return FieldTypeString
}
