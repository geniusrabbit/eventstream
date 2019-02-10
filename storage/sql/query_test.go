package sql

import "testing"

func Test_RawQuery(t *testing.T) {
	tests := []struct {
		query  string
		fields interface{}
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
	}

	for _, test := range tests {
		if _, err := NewQueryByRaw(test.query, test.fields); err != nil {
			t.Error(err)
		}
	}
}

func Test_PatternQuery(t *testing.T) {
	tests := []struct {
		pattern string
		target  string
		fields  interface{}
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
	}

	for _, test := range tests {
		if _, err := NewQueryByPattern(test.pattern, test.target, test.fields); err != nil {
			t.Error(err)
		}
	}
}
