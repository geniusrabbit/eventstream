//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package sql

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/eventstream"
)

var (
	paramsSearch = regexp.MustCompile(`\{\{([^}:]+)(?::([^}|]+))?(?:\|([^}]+))?\}\}`)
	paramParser  = regexp.MustCompile(`([^}:]+)(?::([^}|]+))?(?:\|([^}]+))?`)
)

var (
	errInvalidQueryFieldsParam = errors.New(`[stream::query] invalid fields build param`)
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
	query  string
	values []Value
}

// NewQueryByRaw returns query object from raw SQL query
func NewQueryByRaw(query string, fl interface{}) (queryBuilder *Query, err error) {
	if !isEmptyFields(fl) {
		var (
			values          []Value
			fields, inserts []string
		)
		if values, fields, inserts, err = PrepareFields(fl); err != nil {
			return nil, err
		}
		queryBuilder = &Query{
			query: strings.NewReplacer(
				"{{fields}}", strings.Join(fields, ", "),
				"{{values}}", strings.Join(inserts, ", "),
			).Replace(query),
			values: values,
		}
	} else if args := paramsSearch.FindAllStringSubmatch(query, -1); len(args) > 0 {
		queryBuilder = &Query{query: paramsSearch.ReplaceAllString(query, "?")}
		for _, a := range args {
			queryBuilder.values = append(queryBuilder.values, valueFromArray(a[1], a[1:]))
		}
	}
	return queryBuilder, err
}

// NewQueryByPattern returns query object
func NewQueryByPattern(pattern, target string, fl interface{}) (_ *Query, err error) {
	var (
		fields, inserts []string
		values          []Value
	)
	if fl == nil {
		return nil, errInvalidQueryFieldsParam
	}
	if values, fields, inserts, err = PrepareFields(fl); nil != err {
		return nil, err
	}
	return &Query{
		query: strings.NewReplacer(
			"{{target}}", target,
			"{{fields}}", strings.Join(fields, ", "),
			"{{values}}", strings.Join(inserts, ", "),
		).Replace(pattern),
		values: values,
	}, nil
}

// QueryString - returns the SQL query string
func (q *Query) QueryString() string {
	return q.query
}

// ParamsBy by message
func (q *Query) ParamsBy(msg eventstream.Message) (params []interface{}) {
	for _, v := range q.values {
		params = append(params, msg.ItemCast(v.Key, v.Type, v.Length, v.Format))
	}
	return params
}

// StringParamsBy by message
func (q *Query) StringParamsBy(msg eventstream.Message) (params []string) {
	for _, v := range q.values {
		params = append(
			params,
			gocast.ToString(msg.ItemCast(v.Key, v.Type, v.Length, v.Format)),
		)
	}
	return params
}

// StringByMessage prepare
func (q *Query) StringByMessage(msg eventstream.Message) string {
	var (
		params = q.StringParamsBy(msg)
		items  = strings.Split(q.query, "?")
		result bytes.Buffer
	)
	for i, v := range items {
		result.WriteString(v)
		if i < len(params)-1 {
			result.WriteString(params[i])
		}
	}
	return result.String()
}

// Extract message by special fields and types
func (q *Query) Extract(msg eventstream.Message) map[string]interface{} {
	var resp = make(map[string]interface{})
	for _, v := range q.values {
		resp[v.TargetKey] = msg.ItemCast(v.Key, v.Type, v.Length, v.Format)
	}
	return resp
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

// PrepareFields matching
func PrepareFields(fls interface{}) (values []Value, fields, inserts []string, err error) {
	switch fs := fls.(type) {
	case []string:
		values, fields, inserts = PrepareFieldsByArray(fs)
	case []interface{}:
		if len(fs) == 1 {
			switch reflect.TypeOf(fs[0]).Kind() {
			case reflect.Map:
				values, fields, inserts, err = MapIntoQueryParams(fs[0])
			case reflect.Struct, reflect.Ptr:
				values, fields, inserts, err = MapObjectIntoQueryParams(fs[0])
			default:
				values, fields, inserts = PrepareFieldsByArray(gocast.ToStringSlice(fs))
			}
		} else {
			values, fields, inserts = PrepareFieldsByArray(gocast.ToStringSlice(fs))
		}
	case string:
		values, fields, inserts = PrepareFieldsByString(fs)
	default:
		switch reflect.TypeOf(fs).Kind() {
		case reflect.Map:
			values, fields, inserts, err = MapIntoQueryParams(fs)
		default:
			values, fields, inserts, err = MapObjectIntoQueryParams(fs)
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
	if strings.ContainsAny(field, "=") {
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
	} else {
		if !strings.ContainsAny(field, ":|") {
			fieldValue = field
			values = append(values, Value{Key: field, TargetKey: field})
		} else {
			match := paramParser.FindAllStringSubmatch(field, -1)
			fieldValue = match[0][1]
			values = append(values, valueFromArray(match[0][1], match[0][1:]))
		}
	}
	return values, fieldValue, insert
}

func isEmptyFields(fls interface{}) bool {
	switch fs := fls.(type) {
	case []string:
		return len(fs) < 1
	case []interface{}:
		return len(fs) < 1
	case string:
		return strings.TrimSpace(fs) == ""
	case nil:
		return true
	}
	return false
}
