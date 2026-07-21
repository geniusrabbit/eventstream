// Package condition provides message filters used by stream wrappers
// (expression-based and custom predicates).
package condition

import (
	"context"

	"github.com/geniusrabbit/eventstream/pkg/message"
)

// Condition decides whether a message should be processed by a stream.
type Condition interface {
	Check(ctx context.Context, msg message.Message) bool
}
