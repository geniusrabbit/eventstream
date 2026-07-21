package condition

import (
	"context"

	"github.com/demdxx/gocast/v2"
	"github.com/tonalfitness/govaluate/v3"

	"github.com/geniusrabbit/eventstream/pkg/message"
)

// Expression evaluates a govaluate where-clause against a message map.
type Expression struct {
	whereCondition *govaluate.EvaluableExpression
}

// NewExpression compiles a where expression. An empty string matches all messages.
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
	res, err := e.whereCondition.Eval(&messageWrapper{Message: msg})
	return err == nil && gocast.Bool(res)
}

type messageWrapper struct {
	message.Message
}

func (m *messageWrapper) Get(name string) (any, error) {
	return m.Item(name, nil), nil
}
