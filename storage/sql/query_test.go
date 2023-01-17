package sql

import (
	"testing"
	"time"
)

func TestRawQuery(t *testing.T) {
	tests := []struct {
		query  string
		fields any
	}{
		{
			query:  `INSERT INTO my_table ({{fields}}) VALUES({{values}})`,
			fields: "service=srv,msg,error=err,timestamp=@toTimestamp({{timestamp:date}})",
		},
		{
			query: `COPY test ({{fields}}) FROM STDIN DELIMITER '\t' NULL 'null'`,
			fields: []string{
				"service=srv",
				"msg",
				"error=err:string",
				"timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
			},
		},
		{
			query: `INSERT INTO testlog (service, msg, error, timestamp)` +
				`VALUES({{srv}}, {{msg}}, {{err}}, toTimestamp({{timestamp:date}}))`,
		},
		{
			query: `INSERT INTO my_table ({{fields}}) VALUES({{values}})`,
			fields: struct {
				Service   string    `field:"srv" target:"service"`
				Message   string    `field:"msg"`
				Error     string    `field:"err" target:"error"`
				Timestamp time.Time `field:"timestamp" format:"2006-01-02"`
			}{
				Service:   "test",
				Message:   "msg",
				Error:     "error",
				Timestamp: time.Now(),
			},
		},
		{
			query: `INSERT INTO my_table ({{fields}}) VALUES({{values}})`,
			fields: map[string]any{
				"service":   "srv",
				"msg":       "msg",
				"error":     "err",
				"timestamp": "timestamp:date|2006-01-02",
			},
		},
		{
			query: `INSERT INTO my_table ({{fields}}) VALUES({{values}})`,
			fields: []any{map[string]any{
				"service":   "srv",
				"msg":       "msg",
				"error":     "err",
				"timestamp": "timestamp:date|2006-01-02",
			}},
		},
	}

	for _, test := range tests {
		if _, err := NewQueryByRaw(test.query, test.fields); err != nil {
			t.Error(err)
		}
	}
}

func TestPatternQuery(t *testing.T) {
	tests := []struct {
		pattern string
		target  string
		fields  any
	}{
		{
			pattern: `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
			target:  "test_target",
			fields:  "service=srv,msg,error=err,timestamp=@toTimestamp({{timestamp:date}})",
		},
		{
			pattern: `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
			target:  "test_target",
			fields: []string{
				"service=srv",
				"msg",
				"error=err:string",
				"timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
			},
		},
		{
			pattern: `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
			target:  "test_target",
			fields: map[string]any{
				"service":   "srv",
				"msg":       "msg",
				"error":     "err",
				"timestamp": "timestamp:date|2006-01-02",
			},
		},
		{
			pattern: `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
			target:  "test_target",
			fields: []any{map[string]any{
				"service":   "srv",
				"msg":       "msg",
				"error":     "err",
				"timestamp": "timestamp:date|2006-01-02",
			}},
		},
	}

	for _, test := range tests {
		if _, err := NewQueryByPattern(test.pattern, test.target, test.fields); err != nil {
			t.Error(err)
		}
	}
}
