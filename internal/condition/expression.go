package condition

import (
	"context"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/eventstream/internal/message"
	"github.com/tonalfitness/govaluate/v3"
)

// Expression condition checker
type Expression struct {
	whereCondition *govaluate.EvaluableExpression
}

// NewExpression creates new condition expression
func NewExpression(where string) (_ *Expression, err error) {
	var whereObj *govaluate.EvaluableExpression
	if len(where) > 0 {
		if whereObj, err = govaluate.NewEvaluableExpression(where); err != nil {
			return nil, err
		}
	}
	return &Expression{
		whereCondition: whereObj,
	}, nil
}

// Check message by condition
func (e *Expression) Check(ctx context.Context, msg message.Message) bool {
	if e.whereCondition == nil {
		return true
	}
	res, err := e.whereCondition.Eval(govaluate.MapParameters(msg.Map()))
	return err == nil && gocast.Bool(res)
}
