package eventstream

import "github.com/geniusrabbit/eventstream/internal/message"

type unmarshalel interface {
	Unmarshal(data []byte, v interface{}) error
}

type (
	Formater  = message.Formater
	FieldType = message.FieldType
	Message   = message.Message
)

// Field scalar types enum
const (
	FieldTypeString     = message.FieldTypeString
	FieldTypeFixed      = message.FieldTypeFixed
	FieldTypeUUID       = message.FieldTypeUUID
	FieldTypeInt        = message.FieldTypeInt
	FieldTypeInt32      = message.FieldTypeInt32
	FieldTypeInt8       = message.FieldTypeInt8
	FieldTypeUint       = message.FieldTypeUint
	FieldTypeUint32     = message.FieldTypeUint32
	FieldTypeUint8      = message.FieldTypeUint8
	FieldTypeFloat      = message.FieldTypeFloat
	FieldTypeBoolean    = message.FieldTypeBoolean
	FieldTypeIP         = message.FieldTypeIP
	FieldTypeDate       = message.FieldTypeDate
	FieldTypeUnixnano   = message.FieldTypeUnixnano
	FieldTypeArrayInt32 = message.FieldTypeArrayInt32
	FieldTypeArrayInt64 = message.FieldTypeArrayInt64
)

// TypeByString name
func TypeByString(t string) FieldType {
	return message.TypeByString(t)
}

// MessageDecode from bytes
func MessageDecode(data []byte, converter unmarshalel) (msg Message, err error) {
	return message.MessageDecode(data, converter)
}
