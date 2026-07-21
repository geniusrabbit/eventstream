package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/geniusrabbit/eventstream/pkg/message"
)

// Mapper transforms a message into zero or more outgoing messages.
type Mapper interface {
	Map(ctx context.Context, msg message.Message) ([]message.Message, error)
}

type mapperFunc struct {
	fn func(ctx context.Context, msg message.Message) ([]message.Message, error)
}

func (m mapperFunc) Map(ctx context.Context, msg message.Message) ([]message.Message, error) {
	return m.fn(ctx, msg)
}

var (
	messageType    = reflect.TypeFor[message.Message]()
	mapAnyType     = reflect.TypeFor[map[string]any]()
	mapMessageType = reflect.TypeFor[message.MapMessage]()
)

// MapperFunc wraps a mapping function. Supported shapes (In / Out inferred):
//
//	func(ctx, message.Message) ([]message.Message, error)
//	func(ctx, message.Message) (message.Message, error)
//	func(ctx, T) ([]message.Message, error)
//	func(ctx, T) (map[string]any, error)
//	func(ctx, *T) (map[string]any, error)
//	func(ctx, *T) (*U, error)
//	func(ctx, *T) ([]*U, error)
func MapperFunc[In, Out any](fn func(ctx context.Context, in In) (Out, error)) Mapper {
	return mapperFunc{fn: func(ctx context.Context, msg message.Message) ([]message.Message, error) {
		in, err := decodeIn[In](msg)
		if err != nil {
			return nil, err
		}
		out, err := fn(ctx, in)
		if err != nil {
			return nil, err
		}
		return encodeOut(out)
	}}
}

// MapperListFunc wraps a typed 1→N mapping function.
func MapperListFunc[In, Out any](fn func(ctx context.Context, in *In) ([]*Out, error)) Mapper {
	return MapperFunc(fn)
}

// FieldMap renames fields from In JSON/field tags into Out using mapping
// where key is output field name and value is input field name.
func FieldMap[In, Out any](mapping map[string]string) Mapper {
	return MapperFunc(func(_ context.Context, in *In) (*Out, error) {
		inMap, err := structToMap(in)
		if err != nil {
			return nil, err
		}
		outMap := make(map[string]any, len(mapping))
		for outKey, inKey := range mapping {
			if v, ok := inMap[inKey]; ok {
				outMap[outKey] = v
			}
		}
		return mapToStruct[Out](outMap)
	})
}

func decodeIn[In any](msg message.Message) (In, error) {
	var zero In
	inType := reflect.TypeFor[In]()

	switch {
	case inType == messageType:
		return any(msg).(In), nil
	case inType == mapAnyType:
		return any(msg.Map()).(In), nil
	case inType == mapMessageType:
		if mm, ok := msg.(message.MapMessage); ok {
			return any(mm).(In), nil
		}
		return any(message.MapMessage(msg.Map())).(In), nil
	case inType.Kind() == reflect.Pointer:
		elem := reflect.New(inType.Elem())
		if err := json.Unmarshal([]byte(msg.JSON()), elem.Interface()); err != nil {
			return zero, fmt.Errorf("mapper decode: %w", err)
		}
		return elem.Interface().(In), nil
	default:
		var v In
		if err := json.Unmarshal([]byte(msg.JSON()), &v); err != nil {
			return zero, fmt.Errorf("mapper decode: %w", err)
		}
		return v, nil
	}
}

func encodeOut[Out any](out Out) ([]message.Message, error) {
	v := any(out)
	if isNil(v) {
		return nil, nil
	}

	switch o := v.(type) {
	case message.Message:
		return []message.Message{o}, nil
	case []message.Message:
		return o, nil
	case map[string]any:
		return []message.Message{message.MapMessage(o)}, nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice {
		result := make([]message.Message, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i).Interface()
			if isNil(elem) {
				continue
			}
			if m, ok := elem.(message.Message); ok {
				result = append(result, m)
				continue
			}
			encoded, err := encodeValue(elem)
			if err != nil {
				return nil, err
			}
			result = append(result, encoded)
		}
		return result, nil
	}

	encoded, err := encodeValue(v)
	if err != nil {
		return nil, err
	}
	return []message.Message{encoded}, nil
}

func encodeValue(v any) (message.Message, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("mapper encode: %w", err)
	}
	var m message.MapMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("mapper encode: %w", err)
	}
	return m, nil
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
		return rv.IsNil()
	default:
		return false
	}
}

func structToMap(v any) (map[string]any, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func mapToStruct[T any](m map[string]any) (*T, error) {
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
