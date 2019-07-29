//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package sql

import (
	"bytes"
	"errors"
	"fmt"
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
	errInvalidQueryFields = errors.New(`[stream::query] invalid fields build param`)
)

// Value item
type Value struct {
	Key       string
	TargetKey string
	Type      eventstream.FieldType
	Length    int
	Format    string
}

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
	Q      string
	Values []Value
}

// NewQueryByRaw returns query object from raw SQL query
func NewQueryByRaw(query string, fl interface{}) (q *Query, err error) {
	if !isEmptyFields(fl) {
		var (
			values          []Value
			fields, inserts []string
		)
		if values, fields, inserts, err = PrepareFields(fl); nil != err {
			return nil, err
		}
		q = &Query{
			Q: strings.NewReplacer(
				"{{fields}}", strings.Join(fields, ", "),
				"{{values}}", strings.Join(inserts, ", "),
			).Replace(query),
			Values: values,
		}
	} else if args := paramsSearch.FindAllStringSubmatch(query, -1); len(args) > 0 {
		q = &Query{Q: paramsSearch.ReplaceAllString(query, "?")}
		for _, a := range args {
			q.Values = append(q.Values, valueFromArray(a[1], a[1:]))
		}
	}
	return
}

// NewQueryByPattern returns query object
func NewQueryByPattern(pattern, target string, fl interface{}) (_ *Query, err error) {
	var (
		fields, inserts []string
		values          []Value
	)

	if fl == nil {
		return nil, errInvalidQueryFields
	}

	if values, fields, inserts, err = PrepareFields(fl); nil != err {
		return nil, err
	}

	return &Query{
		Q: strings.NewReplacer(
			"{{target}}", target,
			"{{fields}}", strings.Join(fields, ", "),
			"{{values}}", strings.Join(inserts, ", "),
		).Replace(pattern),
		Values: values,
	}, nil
}

// ParamsBy by message
func (q *Query) ParamsBy(msg eventstream.Message) (params []interface{}) {
	for _, v := range q.Values {
		params = append(params, msg.ItemCast(v.Key, v.Type, v.Length, v.Format))
	}
	return
}

// StringParamsBy by message
func (q *Query) StringParamsBy(msg eventstream.Message) (params []string) {
	for _, v := range q.Values {
		params = append(
			params,
			gocast.ToString(msg.ItemCast(v.Key, v.Type, v.Length, v.Format)),
		)
	}
	return
}

// StringByMessage prepare
func (q *Query) StringByMessage(msg eventstream.Message) string {
	var (
		params = q.StringParamsBy(msg)
		items  = strings.Split(q.Q, "?")
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
	for _, v := range q.Values {
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
		values, fields, inserts = PrepareFieldsByArray(gocast.ToStringSlice(fs))
	case string:
		values, fields, inserts = PrepareFieldsByString(fs)
	default:
		fmt.Println("YYY", fls)
		err = errInvalidQueryFields
	}

	if len(inserts) < 1 || len(fields) > len(values) {
		fmt.Println("ZZZ", fls)
		err = errInvalidQueryFields
	}
	return
}

// PrepareFieldsByArray matching and returns raw fields for insert
// Example: [service=srv:int, name:string]
// Result: [srv:int, name:string], [service,name], [?,?]
func PrepareFieldsByArray(fls []string) (values []Value, fields, inserts []string) {
	for _, fl := range fls {
		if "" == fl {
			continue
		}

		if strings.ContainsAny(fl, "=") {
			if field := strings.SplitN(fl, "=", 2); '@' == field[1][0] {
				fields = append(fields, field[0])
				insert := field[1][1:]
				if match := paramsSearch.FindAllStringSubmatch(insert, -1); len(match) > 0 {
					for _, v := range match {
						values = append(values, valueFromArray(v[1], v[1:]))
						insert = strings.Replace(insert, v[0], "?", -1)
					}
				}
				inserts = append(inserts, insert)
			} else {
				fields = append(fields, field[0])
				if !strings.ContainsAny(field[1], ":|") {
					values = append(values, Value{Key: field[1], TargetKey: field[0]})
				} else {
					match := paramParser.FindAllStringSubmatch(field[1], -1)
					values = append(values, valueFromArray(field[0], match[0][1:]))
				}
				inserts = append(inserts, "?")
			}
		} else {
			if !strings.ContainsAny(fl, ":|") {
				fields = append(fields, fl)
				values = append(values, Value{Key: fl, TargetKey: fl})
			} else {
				match := paramParser.FindAllStringSubmatch(fl, -1)
				fields = append(fields, match[0][1])
				values = append(values, valueFromArray(match[0][1], match[0][1:]))
			}
			inserts = append(inserts, "?")
		}
	}
	return
}

// PrepareFieldsByString matching and returns raw fields for insert
func PrepareFieldsByString(fls string) (values []Value, fields, inserts []string) {
	return PrepareFieldsByArray(strings.Split(fls, ","))
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
	default:
	}
	return true
}
