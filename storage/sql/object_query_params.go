package sql

import (
	"errors"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	ev "github.com/geniusrabbit/eventstream"
)

var (
	errInvalidValue         = errors.New("[sql.queryObject] object value is not valid")
	errUnsupportedValueType = errors.New("[sql.queryObject] object value type os not supported")
)

var (
	ipType         = reflect.TypeOf(net.IP{})
	timeType       = reflect.TypeOf(time.Time{})
	int32ArrayType = reflect.TypeOf([]int32{})
	int64ArrayType = reflect.TypeOf([]int64{})
)

// MapObjectIntoQueryParams returns object which links object Fields and Data Fields
func MapObjectIntoQueryParams(object interface{}) (values []Value, fields, inserts []string, err error) {
	objectType, err := reflectTargetStruct(reflect.ValueOf(object))
	if err != nil {
		return nil, nil, nil, err
	}
	for i := 0; i < objectType.NumField(); i++ {
		field := objectType.Field(i)
		// defval := metaByTags(field, "", "field_default", "default")
		name := metaByTags(field, field.Name, "field", "json")
		format := metaByTags(field, "", "field_format", "format")
		target := metaByTags(field, "", "field_target", "target")
		defexp := metaByTags(field, "", "field_defexp", "defexp")
		size := metaByTags(field, "", "field_size", "size")
		fltype := metaByTags(field, "", "field_type", "type")
		if target == "" {
			target = name
		}
		fields = append(fields, target)

		if defexp != "" {
			insert := defexp
			if match := paramsSearch.FindAllStringSubmatch(insert, -1); len(match) > 0 {
				for _, v := range match {
					if v[2] == "" {
						v[2] = typeByValue(field.Type, fltype, size).String()
					}
					values = append(values, valueFromArray(v[1], v[1:]))
					insert = strings.Replace(insert, v[0], "?", -1)
				}
			}
			inserts = append(inserts, insert)
			continue
		}

		inserts = append(inserts, "?")
		length, _ := strconv.Atoi(size)
		values = append(values, Value{
			Key:       name,
			TargetKey: target,
			Type:      typeByValue(field.Type, fltype, size),
			Length:    length,
			Format:    format,
		})
	}
	return values, fields, inserts, nil
}

func reflectTargetStruct(val reflect.Value) (reflect.Type, error) {
	for {
		if !val.IsValid() {
			return nil, errInvalidValue
		}
		switch val.Kind() {
		case reflect.Struct:
			return val.Type(), nil
		case reflect.Interface, reflect.Ptr:
			val = val.Elem()
		default:
			return nil, errUnsupportedValueType
		}
	}
}

func metaByTags(field reflect.StructField, def string, tags ...string) (v string) {
	for _, tag := range tags {
		if v = field.Tag.Get(tag); v != "" {
			return v
		}
	}
	return def
}

func typeByValue(tp reflect.Type, fieldType, size string) ev.FieldType {
	switch tp.Kind() {
	case reflect.String:
		switch {
		case fieldType == "fixed" && size != "":
			return ev.FieldTypeFixed
		case fieldType == "uuid":
			return ev.FieldTypeUUID
		default:
			return ev.FieldTypeString
		}
	case reflect.Int, reflect.Int64:
		if fieldType == "unixnano" {
			return ev.FieldTypeUnixnano
		}
		return ev.FieldTypeInt
	case reflect.Int32:
		return ev.FieldTypeInt32
	case reflect.Int8:
		return ev.FieldTypeInt8
	case reflect.Uint, reflect.Uint64:
		return ev.FieldTypeUint
	case reflect.Uint32:
		return ev.FieldTypeUint32
	case reflect.Uint8:
		return ev.FieldTypeUint8
	case reflect.Float32, reflect.Float64:
		return ev.FieldTypeFloat
	case reflect.Bool:
		return ev.FieldTypeBoolean
	default:
		switch tp {
		case ipType:
			return ev.FieldTypeIP
		case timeType:
			if fieldType == "unixnano" {
				return ev.FieldTypeUnixnano
			}
			return ev.FieldTypeDate
		case int32ArrayType:
			return ev.FieldTypeArrayInt32
		case int64ArrayType:
			return ev.FieldTypeArrayInt64
		}
	}
	return ev.FieldTypeString
}
