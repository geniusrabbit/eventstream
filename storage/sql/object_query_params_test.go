package sql

import (
	"reflect"
	"testing"
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
