package eventstream

import "github.com/geniusrabbit/eventstream/pkg/message"

type unmarshaler interface {
	Unmarshal(data []byte, v any) error
}

// EmptyMessage is a nil-safe empty message.
var EmptyMessage = message.Empty

// Message and MapMessage re-export pkg/message types for the root API.
type (
	Message    = message.Message
	MapMessage = message.MapMessage
)

// MessageDecode decodes bytes into a Message using the given unmarshaler.
func MessageDecode(data []byte, converter unmarshaler) (msg Message, err error) {
	return message.Decode(data, converter)
}
