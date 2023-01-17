package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareValue(t *testing.T) {
	var tests = []struct {
		val    string
		target string
		vars   []map[string]string
	}{
		{val: "{{ @env:USER}}", target: os.Getenv("USER")},
		{val: "{{@env:HOME }}", target: os.Getenv("HOME")},
		{val: "dir:{{ @env:HOME }}", target: "dir:" + os.Getenv("HOME")},
		{
			val:    "connect:{{host}}:{{ port }}/{{db}}?user={{ @env:USER }}",
			target: "connect:localhost:1234/users?user=" + os.Getenv("USER"),
			vars:   []map[string]string{{"host": "localhost", "port": "1234", "db": "users"}},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.target, PrepareValue(test.val, test.vars...))
	}
}
