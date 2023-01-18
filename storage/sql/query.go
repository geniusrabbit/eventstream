//
// @project geniusrabbit::eventstream 2017, 2019 - 2023
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019 - 2023
//

package sql

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/eventstream"
)

var (
	paramsSearch = regexp.MustCompile(`\{\{([^}:]+)(?::([^}|]+))?(?:\|([^}]+))?\}\}`)
	paramParser  = regexp.MustCompile(`([^}:]+)(?::([^}|]+))?(?:\|([^}]+))?`)
)

var (
	errInvalidQueryFieldsValue = errors.New(`[stream::query] invalid fields build value`)
)

// Value item
type Value struct {
	Key       string // Field key in the object
	TargetKey string // Target Key in the database
	Type      eventstream.FieldType
	Length    int
	Format    string
}

// vector of options contains of: {key, typeName, [format]}
func valueFromArray(target string, a []string) Value {
	var (
		length   int64
		typeName string
	)

	if len(a) > 1 {
		vals := strings.Split(a[1], "*")
		if len(vals) == 2 {
			typeName = vals[0]
			length, _ = strconv.ParseInt(vals[1], 10, 64)
		} else {
			typeName = a[1]
		}
	}

	if len(a) > 2 {
		return Value{
			Key:       a[0],
			TargetKey: target,
			Type:      eventstream.TypeByString(typeName),
			Length:    int(length),
			Format:    a[2],
		}
	} else if len(a) > 1 {
		return Value{
			Key:       a[0],
			TargetKey: target,
			Type:      eventstream.TypeByString(typeName),
			Length:    int(length),
		}
	}
	return Value{Key: a[0], TargetKey: target}
}

// Query extractor
type Query struct {
	query     string
	iterateBy string
	values    []Value
}

// NewQuery returns query object from SQL query
func NewQuery(query string, opts ...QueryOption) (*Query, error) {
	var (
		err             error
		queryConfig     QueryConfig
		values          []Value
		fields, inserts []string
	)
	for _, opt := range opts {
		opt(&queryConfig)
	}
	query = strings.ReplaceAll(query, "{{target}}", queryConfig.Target)

	if isEmptyFields(queryConfig.FieldObject) {
		if args := paramsSearch.FindAllStringSubmatch(query, -1); len(args) > 0 {
			query = paramsSearch.ReplaceAllString(query, "?")
			for _, arg := range args {
				values = append(values, valueFromArray(arg[1], arg[1:]))
			}
		}
	} else if values, fields, inserts, err = PrepareFields(queryConfig.FieldObject); err != nil {
		return nil, err
	}

	return &Query{
		query: strings.NewReplacer(
			"{{fields}}", strings.Join(fields, ", "),
			"{{values}}", strings.Join(inserts, ", "),
		).Replace(query),
		iterateBy: queryConfig.IterateBy,
		values:    values,
	}, nil
}

// QueryString - returns the SQL query string
func (q *Query) QueryString() string {
	return q.query
}

// ParamsBy by message
func (q *Query) ParamsBy(msg eventstream.Message) paramsResult {
	var (
		iterated      = q.iterateBy != ``
		iteratorSlice = gocast.Cast[[]any](msg.Item(q.iterateBy, nil))
		countSlice    = tMax(len(iteratorSlice), 1)
		params        = acquireResultParams()
	)
	for i := 0; i < countSlice; i++ {
		subParams := make([]any, 0, len(q.values))
		for _, v := range q.values {
			var val any
			if iterated && strings.HasPrefix(v.Key, "$iter") {
				subItem := iteratorSlice
				if strings.HasPrefix(v.Key, "$iter.") {
					subItem = gocast.Cast[[]any](msg.Item(v.Key[6:], nil))
				}
				val = v.Type.CastExt(arrValue(i, subItem), v.Length, v.Format)
			} else {
				val = msg.ItemCast(v.Key, v.Type, v.Length, v.Format)
			}
			subParams = append(subParams, val)
		}
		params = append(params, subParams)
	}
	return params
}

// Extract message by special fields and types
func (q *Query) Extract(msg eventstream.Message) []map[string]any {
	var (
		iterated      = q.iterateBy != ``
		iteratorSlice = gocast.Cast[[]any](msg.Item(q.iterateBy, nil))
		countSlice    = tMax(len(iteratorSlice), 1)
		resp          = make([]map[string]any, 0, countSlice)
	)
	for i := 0; i < countSlice; i++ {
		subResp := make(map[string]any, len(q.values))
		for _, v := range q.values {
			var val any
			if iterated && strings.HasPrefix(v.Key, "$iter") {
				subItem := iteratorSlice
				if strings.HasPrefix(v.Key, "$iter.") {
					subItem = gocast.Cast[[]any](msg.Item(v.Key[6:], nil))
				}
				val = v.Type.CastExt(arrValue(i, subItem), v.Length, v.Format)
			} else {
				val = msg.ItemCast(v.Key, v.Type, v.Length, v.Format)
			}
			subResp[v.TargetKey] = val
		}
		resp = append(resp, subResp)
	}
	return resp
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

// PrepareFields matching
func PrepareFields(fls any) (values []Value, fields, inserts []string, err error) {
	switch fs := fls.(type) {
	case []string:
		values, fields, inserts = PrepareFieldsByArray(fs)
	case []any:
		if len(fs) == 1 {
			switch reflect.TypeOf(fs[0]).Kind() {
			case reflect.Map:
				values, fields, inserts, err = MapIntoQueryParams(fs[0])
			case reflect.Struct, reflect.Ptr:
				values, fields, inserts, err = ConvertObjectIntoQueryParams(fs[0])
			default:
				values, fields, inserts = PrepareFieldsByArray(gocast.Slice[string](fs))
			}
		} else {
			values, fields, inserts = PrepareFieldsByArray(gocast.Slice[string](fs))
		}
	case string:
		values, fields, inserts = PrepareFieldsByString(fs)
	default:
		switch reflect.TypeOf(fs).Kind() {
		case reflect.Map:
			values, fields, inserts, err = MapIntoQueryParams(fs)
		default:
			values, fields, inserts, err = ConvertObjectIntoQueryParams(fs)
		}
	}
	if err == errInvalidValue || (err == nil && (len(inserts) < 1 || len(fields) > len(values))) {
		err = errInvalidQueryFieldsValue
	}
	return values, fields, inserts, err
}

// PrepareFieldsByArray matching and returns raw fields for insert
// Example: [service=srv:int, name:string]
// Result: [srv:int, name:string], [service,name], [?,?]
func PrepareFieldsByArray(fields []string) (values []Value, fieldNames, inserts []string) {
	for _, field := range fields {
		if len(field) == 0 {
			continue
		}
		fieldValues, fieldValue, insertValue := prepareOneFieldByString(field)
		if len(fieldValues) != 0 {
			values = append(values, fieldValues...)
		}
		if len(fieldValue) != 0 {
			fieldNames = append(fieldNames, fieldValue)
		}
		if len(insertValue) != 0 {
			inserts = append(inserts, insertValue)
		}
	}
	return values, fieldNames, inserts
}

// PrepareFieldsByString matching and returns raw fields for insert
func PrepareFieldsByString(fls string) (values []Value, fields, inserts []string) {
	return PrepareFieldsByArray(strings.Split(fls, ","))
}

func prepareOneFieldByString(field string) (values []Value, fieldValue, insert string) {
	insert = "?"
	switch {
	case strings.ContainsAny(field, "="):
		if field := strings.SplitN(field, "=", 2); field[1][0] == '@' {
			fieldValue = field[0]
			insert = field[1][1:]
			if match := paramsSearch.FindAllStringSubmatch(insert, -1); len(match) > 0 {
				for _, v := range match {
					values = append(values, valueFromArray(v[1], v[1:]))
					insert = strings.Replace(insert, v[0], "?", -1)
				}
			}
		} else {
			fieldValue = field[0]
			if !strings.ContainsAny(field[1], ":|") {
				values = append(values, Value{Key: field[1], TargetKey: field[0]})
			} else {
				match := paramParser.FindAllStringSubmatch(field[1], -1)
				values = append(values, valueFromArray(field[0], match[0][1:]))
			}
		}
	case !strings.ContainsAny(field, ":|"):
		fieldValue = field
		values = append(values, Value{Key: field, TargetKey: field})
	default:
		match := paramParser.FindAllStringSubmatch(field, -1)
		fieldValue = match[0][1]
		values = append(values, valueFromArray(match[0][1], match[0][1:]))
	}
	return values, fieldValue, insert
}

func isEmptyFields(fls any) bool {
	switch fs := fls.(type) {
	case []string:
		return len(fs) == 0
	case []any:
		return len(fs) == 0
	case string:
		return strings.TrimSpace(fs) == ""
	case nil:
		return true
	}
	return gocast.IsEmpty(fls)
}

func tMax[T gocast.Numeric](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func arrValue[T any](idx int, arr []T) any {
	if idx < len(arr) {
		return arr[idx]
	}
	return nil
}
