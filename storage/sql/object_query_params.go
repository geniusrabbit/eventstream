package sql

import (
	"errors"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/google/uuid"
)

var (
	errInvalidValue         = errors.New("[sql.queryObject] value is not valid")
	errInvalidObjectValue   = errors.New("[sql.queryObject] object value is not valid")
	errInvalidMapValue      = errors.New("[sql.queryObject] map value is not valid")
	errUnsupportedValueType = errors.New("[sql.queryObject] object value type is not supported")
)

var (
	ipType          = reflect.TypeFor[net.IP]()
	timeType        = reflect.TypeFor[time.Time]()
	uuidType        = reflect.TypeFor[uuid.UUID]()
	int32ArrayType  = reflect.TypeFor[[]int32]()
	int64ArrayType  = reflect.TypeFor[[]int64]()
	stringArrayType = reflect.TypeFor[[]string]()
)

// ConvertObjectIntoQueryParams returns object which links object Fields and Data Fields
func ConvertObjectIntoQueryParams(object any) (values []Value, fields, inserts []string, err error) {
	objectValue, err := reflectTargetStruct(reflect.ValueOf(object))
	if err != nil {
		return nil, nil, nil, err
	}
	if objectValue.Kind() != reflect.Struct {
		return nil, nil, nil, errInvalidObjectValue
	}
	objectType := objectValue.Type()
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
					insert = strings.ReplaceAll(insert, v[0], "?")
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

// MapIntoQueryParams returns object which links map fields
func MapIntoQueryParams(object any) (values []Value, fieldNames, inserts []string, err error) {
	objectValue, err := reflectTargetStruct(reflect.ValueOf(object))
	if err != nil {
		return nil, nil, nil, err
	}
	if objectValue.Kind() != reflect.Map {
		return nil, nil, nil, errInvalidMapValue
	}
	data, err := gocast.TryMap[string, string](object)
	if err != nil {
		return nil, nil, nil, err
	}
	orderedData := make([][2]string, 0, len(data))
	for target, value := range data {
		orderedData = append(orderedData, [2]string{target, value})
	}
	sort.Slice(orderedData, func(i, j int) bool { return orderedData[i][0] < orderedData[j][0] })
	for _, item := range orderedData {
		target, value := item[0], item[1]
		fieldValues, fieldValue, insertValue := prepareOneFieldByString(target + `=` + value)
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
	return values, fieldNames, inserts, nil
}

func reflectTargetStruct(val reflect.Value) (reflect.Value, error) {
	for {
		if !val.IsValid() {
			return reflect.Value{}, errInvalidValue
		}
		switch val.Kind() {
		case reflect.Struct:
			return val, nil
		case reflect.Interface, reflect.Pointer:
			val = val.Elem()
		case reflect.Map:
			return val, nil
		default:
			return reflect.Value{}, errUnsupportedValueType
		}
	}
}

func metaByTags(field reflect.StructField, def string, tags ...string) string {
	for _, tag := range tags {
		if v := field.Tag.Get(tag); v != "" {
			return v
		}
	}
	return def
}

func typeByValue(tp reflect.Type, fieldType, size string) FieldType {
	switch tp.Kind() {
	case reflect.String:
		switch {
		case fieldType == "fixed" && size != "":
			return FieldTypeFixed
		case fieldType == "uuid":
			return FieldTypeUUID
		default:
			return FieldTypeString
		}
	case reflect.Int, reflect.Int64:
		if fieldType == "unixnano" {
			return FieldTypeUnixnano
		}
		return FieldTypeInt
	case reflect.Int32:
		return FieldTypeInt32
	case reflect.Int8:
		return FieldTypeInt8
	case reflect.Uint, reflect.Uint64:
		return FieldTypeUint
	case reflect.Uint32:
		return FieldTypeUint32
	case reflect.Uint8:
		return FieldTypeUint8
	case reflect.Float32, reflect.Float64:
		return FieldTypeFloat
	case reflect.Bool:
		return FieldTypeBoolean
	default:
		switch tp {
		case ipType:
			return FieldTypeIP
		case timeType:
			if fieldType == "unixnano" {
				return FieldTypeUnixnano
			}
			return FieldTypeDate
		case uuidType:
			return FieldTypeUUID
		case int32ArrayType:
			return FieldTypeArrayInt32
		case int64ArrayType:
			return FieldTypeArrayInt64
		case stringArrayType:
			return FieldTypeArrayString
		}
	}
	return FieldTypeString
}
