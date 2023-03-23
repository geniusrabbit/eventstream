package condition

import (
	"context"
	"fmt"
	"testing"

	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/stretchr/testify/assert"
)

func TestCondition(t *testing.T) {
	cond1, err := NewExpression("a == 1")
	if !assert.NoError(t, err) {
		return
	}
	cond2 := Custom(func(ctx context.Context, msg message.Message) bool {
		return msg.ItemCast("a", message.FieldTypeInt, 0, "") == int64(1)
	})
	conds := []Condition{cond1, cond2}
	ctx := context.Background()

	for _, cond := range conds {
		t.Run(fmt.Sprintf("%T", cond), func(t *testing.T) {
			assert.True(t, cond.Check(ctx, message.MapMessage{"a": 1}))
			assert.False(t, cond.Check(ctx, message.MapMessage{"a": 2}))
		})
	}
}
