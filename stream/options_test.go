package stream

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestOptins(t *testing.T) {
	var tests = []struct {
		options []Option
		config  Config
	}{
		{
			options: []Option{WithConfig(&Config{Name: "test1"})},
			config:  Config{Name: "test1"},
		},
		{
			options: []Option{WithName("test2")},
			config:  Config{Name: "test2"},
		},
		{
			options: []Option{WithName("test3"), WithDebug(true)},
			config:  Config{Name: "test3", Debug: true},
		},
		{
			options: []Option{WithName("test4"), WithWhere("id = 100")},
			config:  Config{Name: "test4", Where: "id = 100"},
		},
		{
			options: []Option{WithName("test5"), WithRawConfig(json.RawMessage("raw"))},
			config:  Config{Name: "test5", Raw: json.RawMessage("raw")},
		},
		{
			options: []Option{WithName("test5"), WithObjectConfig(100)},
			config:  Config{Name: "test5", Raw: json.RawMessage("100")},
		},
	}

	for _, test := range tests {
		var cnf Config
		for _, opt := range test.options {
			opt(&cnf)
		}
		if !reflect.DeepEqual(&cnf, &test.config) {
			t.Errorf("[%s] invalid option result", test.config.Name)
		}
	}
}
