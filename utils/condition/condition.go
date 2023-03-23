package condition

import (
	"context"

	"github.com/geniusrabbit/eventstream/internal/message"
)

type Condition interface {
	Check(ctx context.Context, msg message.Message) bool
}
