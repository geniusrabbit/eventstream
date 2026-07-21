// Package message defines the transport bag used across the pipeline,
// conditions and storage sinks.
package message

import (
	"github.com/pkg/errors"
)

// ErrFieldNotFound is returned when a requested field is missing.
var (
	ErrFieldNotFound = errors.New("field not found")
)

// Unmarshaler decodes bytes into a value (json/bson/raw converters).
type Unmarshaler interface {
	Unmarshal(data []byte, v any) error
}

// Message is a transport bag used across pipeline, conditions and sinks.
type Message interface {
	Item(key string, def any) any
	Str(key, def string) string
	JSON() string
	Map() map[string]any
}

type dummyMessage struct{}

func (m *dummyMessage) Item(key string, def any) any { return def }
func (m *dummyMessage) Str(key, def string) string   { return def }
func (m *dummyMessage) JSON() string                 { return "null" }
func (m *dummyMessage) Map() map[string]any          { return nil }

// Empty is a nil-safe empty message.
var Empty Message = (*dummyMessage)(nil)
