package storage

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
			options: []Option{WithConfig(&Config{Connect: "test1"})},
			config:  Config{Connect: "test1"},
		},
		{
			options: []Option{WithConnect("driver", "test2")},
			config:  Config{Connect: "test2", Driver: "driver"},
		},
		{
			options: []Option{WithConnect("driver", "test3"), WithDebug(true)},
			config:  Config{Connect: "test3", Driver: "driver", Debug: true},
		},
		{
			options: []Option{WithConnect("driver", "test4"), WithBuffer(100)},
			config:  Config{Connect: "test4", Driver: "driver", Buffer: 100},
		},
		{
			options: []Option{WithConnect("driver", "test5"), WithRawConfig(json.RawMessage("raw"))},
			config:  Config{Connect: "test5", Driver: "driver", Raw: json.RawMessage("raw")},
		},
		{
			options: []Option{WithConnect("driver", "test6"), WithObjectConfig(100)},
			config:  Config{Connect: "test6", Driver: "driver", Raw: json.RawMessage("100")},
		},
	}

	for _, test := range tests {
		var cnf Config
		for _, opt := range test.options {
			opt(&cnf)
		}
		if !reflect.DeepEqual(&cnf, &test.config) {
			t.Errorf("[%s] invalid option result", test.config.Connect)
		}
	}
}
