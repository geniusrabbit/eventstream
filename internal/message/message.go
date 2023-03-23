package message

import (
	"github.com/pkg/errors"
)

// Errors set
var (
	ErrUndefinedDataType       = errors.New("undefined message data types")
	ErrInvalidMessageFieldType = errors.New("invalid message field type")
	ErrFieldNotFound           = errors.New("field not found")
)

type unmarshalel interface {
	Unmarshal(data []byte, v any) error
}

// Message interface
// TODO: refactor this interface
type Message interface {
	Get(name string) (any, error)
	Item(key string, def any) any
	ItemCast(key string, fieldType FieldType, length int, format string) any
	Str(key, def string) string
	JSON() string
	Map() map[string]any
}

type dummyMessage struct{}

func (m *dummyMessage) Get(name string) (any, error) { return nil, errors.Wrap(ErrFieldNotFound, name) }
func (m *dummyMessage) Item(key string, def any) any { return def }
func (m *dummyMessage) ItemCast(key string, fieldType FieldType, length int, format string) any {
	return fieldType.CastExt(m.Item(key, nil), length, format)
}
func (m *dummyMessage) Str(key, def string) string { return def }
func (m *dummyMessage) JSON() string               { return "null" }
func (m *dummyMessage) Map() map[string]any        { return nil }

var EmptyMessage Message = (*dummyMessage)(nil)
