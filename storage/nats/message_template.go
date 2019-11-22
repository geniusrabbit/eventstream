package nats

import (
	"github.com/Knetic/govaluate"
	"github.com/geniusrabbit/eventstream"
)

type messageTemplate struct {

	// WhereCondition of stream
	whereCondition *govaluate.EvaluableExpression
}

func (t *messageTemplate) prepare(msg eventstream.Message) interface{} {
	return msg.Map()
}
