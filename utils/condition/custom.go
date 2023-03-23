package condition

import (
	"context"

	"github.com/geniusrabbit/eventstream/internal/message"
)

// Custom condition checker
type Custom func(ctx context.Context, msg message.Message) bool

// Check message by condition
func (c Custom) Check(ctx context.Context, msg message.Message) bool {
	return c(ctx, msg)
}
