package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geniusrabbit/eventstream/pkg/message"
	"github.com/geniusrabbit/eventstream/pkg/pipeline"
)

type inMsg struct {
	ID  string `json:"id"`
	Val int    `json:"val"`
}

type outMsg struct {
	Sign int    `json:"sign"`
	Text string `json:"text"`
}

func TestMapperFuncTypedPtr(t *testing.T) {
	m := pipeline.MapperFunc(func(_ context.Context, in *inMsg) (*outMsg, error) {
		return &outMsg{Sign: in.Val, Text: in.ID}, nil
	})
	outs, err := m.Map(context.Background(), message.MapMessage{"id": "a", "val": 7})
	require.NoError(t, err)
	require.Len(t, outs, 1)
	assert.Equal(t, float64(7), outs[0].Item("sign", 0))
	assert.Equal(t, "a", outs[0].Str("text", ""))
}

func TestMapperFuncMessageToMessage(t *testing.T) {
	m := pipeline.MapperFunc(func(_ context.Context, msg message.Message) (message.Message, error) {
		return message.MapMessage{"text": msg.Str("id", "")}, nil
	})
	outs, err := m.Map(context.Background(), message.MapMessage{"id": "hello"})
	require.NoError(t, err)
	require.Len(t, outs, 1)
	assert.Equal(t, "hello", outs[0].Str("text", ""))
}

func TestMapperFuncMessageToList(t *testing.T) {
	m := pipeline.MapperFunc(func(_ context.Context, msg message.Message) ([]message.Message, error) {
		return []message.Message{
			message.MapMessage{"n": 1},
			message.MapMessage{"n": 2},
		}, nil
	})
	outs, err := m.Map(context.Background(), message.MapMessage{"id": "x"})
	require.NoError(t, err)
	require.Len(t, outs, 2)
}

func TestMapperFuncValueToMap(t *testing.T) {
	m := pipeline.MapperFunc(func(_ context.Context, in inMsg) (map[string]any, error) {
		return map[string]any{"text": in.ID, "sign": in.Val}, nil
	})
	outs, err := m.Map(context.Background(), message.MapMessage{"id": "z", "val": 9})
	require.NoError(t, err)
	require.Len(t, outs, 1)
	assert.Equal(t, "z", outs[0].Str("text", ""))
	assert.Equal(t, 9, outs[0].Item("sign", 0))
}

func TestMapperFuncPtrToMap(t *testing.T) {
	m := pipeline.MapperFunc[*inMsg](func(_ context.Context, in *inMsg) (map[string]any, error) {
		return map[string]any{"text": in.ID, "sign": in.Val}, nil
	})
	outs, err := m.Map(context.Background(), message.MapMessage{"id": "p", "val": 4})
	require.NoError(t, err)
	require.Len(t, outs, 1)
	assert.Equal(t, "p", outs[0].Str("text", ""))
	assert.Equal(t, 4, outs[0].Item("sign", 0))
}

func TestMapperFuncValueToMessageList(t *testing.T) {
	m := pipeline.MapperFunc(func(_ context.Context, in inMsg) ([]message.Message, error) {
		return []message.Message{message.MapMessage{"text": in.ID}}, nil
	})
	outs, err := m.Map(context.Background(), message.MapMessage{"id": "v", "val": 1})
	require.NoError(t, err)
	require.Len(t, outs, 1)
	assert.Equal(t, "v", outs[0].Str("text", ""))
}

func TestMapperListFunc(t *testing.T) {
	m := pipeline.MapperListFunc[inMsg, outMsg](func(_ context.Context, in *inMsg) ([]*outMsg, error) {
		return []*outMsg{
			{Sign: in.Val, Text: in.ID},
			{Sign: in.Val * 2, Text: in.ID + "-2"},
		}, nil
	})
	outs, err := m.Map(context.Background(), message.MapMessage{"id": "x", "val": 3})
	require.NoError(t, err)
	require.Len(t, outs, 2)
}

func TestFieldMap(t *testing.T) {
	m := pipeline.FieldMap[inMsg, outMsg](map[string]string{
		"text": "id",
		"sign": "val",
	})
	outs, err := m.Map(context.Background(), message.MapMessage{"id": "z", "val": 9})
	require.NoError(t, err)
	require.Len(t, outs, 1)
	assert.Equal(t, "z", outs[0].Str("text", ""))
	assert.Equal(t, float64(9), outs[0].Item("sign", 0))
}
