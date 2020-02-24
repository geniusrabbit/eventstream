package sql

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Reflection(t *testing.T) {
	var tests = []struct {
		ref   interface{}
		found bool
	}{
		{
			ref:   nil,
			found: false,
		},
		{
			ref:   &[]string{""},
			found: false,
		},
		{
			ref:   &[]int{1},
			found: false,
		},
		{
			ref:   struct{}{},
			found: true,
		},
		{
			ref:   &struct{}{},
			found: true,
		},
	}

	for _, test := range tests {
		if tp, _ := reflectTargetStruct(reflect.ValueOf(test.ref)); (tp != nil) != test.found {
			t.Errorf("invalid ref[%T] result %v", test.ref, test.found)
		}
	}
}

func Test_MapObjectIntoQueryParams(t *testing.T) {
	var tests = []struct {
		ref     interface{}
		values  []Value
		fields  []string
		inserts []string
	}{
		{
			ref: struct {
				ID          uint64    `field:"id"`
				Title       string    `field:"title" field_type:"string" field_size:"64"`
				Description string    `field:"desc" field_target:"description"`
				SubID       uint32    `field:"sub_id" field_defexp:"{{sub_id}}+1"`
				CreatedAt   time.Time `field:"created_at" field_target:"year" field_format:"2006"`
			}{},
			values: []Value{
				valueFromArray("id", []string{"id", "uint"}),
				valueFromArray("title", []string{"title", "string*64"}),
				valueFromArray("description", []string{"desc", "string"}),
				valueFromArray("sub_id", []string{"sub_id", "uint32"}),
				valueFromArray("year", []string{"created_at", "date", "2006"}),
			},
			fields:  []string{"id", "title", "description", "sub_id", "year"},
			inserts: []string{`?`, `?`, `?`, `?+1`, `?`},
		},
	}

	for _, test := range tests {
		values, fields, inserts, err := MapObjectIntoQueryParams(test.ref)
		assert.NoError(t, err)
		assert.ElementsMatch(t, test.values, values)
		assert.ElementsMatch(t, test.fields, fields)
		assert.ElementsMatch(t, test.inserts, inserts)
	}
}
