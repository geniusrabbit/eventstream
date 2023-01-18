package sql

import (
	"testing"
	"time"

	"github.com/geniusrabbit/eventstream"
	"github.com/stretchr/testify/assert"
)

func TestParamExtractionsQuery(t *testing.T) {
	var (
		msg1 = eventstream.Message{
			"srv":       "info",
			"msg":       "test",
			"err":       "err",
			"timestamp": time.Now().Format(time.RFC3339),
			"ext":       []any{1, 2},
		}
		tests = []struct {
			queryTarget string
			query       string
			target      string
			iterateBy   string
			fields      any
			countTarget int
			msg         eventstream.Message
		}{
			// Raw Queries
			{
				queryTarget: `INSERT INTO testlog (service, msg, error, timestamp) VALUES(?, ?, ?, toTimestamp(?))`,
				query: `INSERT INTO testlog (service, msg, error, timestamp)` +
					` VALUES({{srv}}, {{msg}}, {{err}}, toTimestamp({{timestamp:date}}))`,
				msg: msg1,
			},
			{
				queryTarget: `COPY test (service, msg, error, timestamp) FROM STDIN DELIMITER '\t' NULL 'null'`,
				query:       `COPY test ({{fields}}) FROM STDIN DELIMITER '\t' NULL 'null'`,
				fields: []string{
					"service=srv",
					"msg",
					"error=err:string",
					"timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
				},
				msg: msg1,
			},
			{
				queryTarget: `INSERT INTO my_table (service, msg, error, timestamp) VALUES(?, ?, ?, ?)`,
				query:       `INSERT INTO my_table ({{fields}}) VALUES({{values}})`,
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
				msg: msg1,
			},
			// Patter generation
			{
				queryTarget: `INSERT INTO test_target (service, msg, error, timestamp) VALUES(?, ?, ?, toTimestamp(?))`,
				query:       `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
				target:      "test_target",
				fields:      "service=srv,msg,error=err,timestamp=@toTimestamp({{timestamp:date}})",
			},
			{
				queryTarget: `INSERT INTO test_target (service, msg, error, timestamp) VALUES(?, ?, ?, toTimestamp('?'))`,
				query:       `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
				target:      "test_target",
				fields: []string{
					"service=srv",
					"msg",
					"error=err:string",
					"timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
				},
			},
			{
				queryTarget: `INSERT INTO test_target (error, msg, service, timestamp) VALUES(?, ?, ?, ?)`,
				query:       `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
				target:      "test_target",
				fields: map[string]any{
					"service":   "srv",
					"msg":       "msg",
					"error":     "err",
					"timestamp": "timestamp:date|2006-01-02",
				},
			},
			{
				queryTarget: `INSERT INTO test_target (error, msg, service, timestamp) VALUES(?, ?, ?, ?)`,
				query:       `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
				target:      "test_target",
				fields: []any{map[string]any{
					"service":   "srv",
					"msg":       "msg",
					"error":     "err",
					"timestamp": "timestamp:date|2006-01-02",
				}},
			},
			// Iterated message queries
			{
				queryTarget: `INSERT INTO test_target (service, msg, error, timestamp, ext) VALUES(?, ?, ?, toTimestamp('?'), ?)`,
				query:       `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
				target:      "test_target",
				iterateBy:   "ext",
				fields: []string{
					"service=srv",
					"msg",
					"error=err:string",
					"timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
					"ext=$iter:string",
				},
				countTarget: 2,
				msg:         msg1,
			},
			{
				queryTarget: `INSERT INTO test_target (service, msg, error, timestamp, ext) VALUES(?, ?, ?, toTimestamp('?'), ?)`,
				query:       `INSERT INTO {{target}} ({{fields}}) VALUES({{values}})`,
				target:      "test_target",
				iterateBy:   "ext",
				fields: []string{
					"service=srv",
					"msg",
					"error=err:string",
					"timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
					"ext=$iter.ext:string",
				},
				countTarget: 2,
				msg:         msg1,
			},
		}
	)
	for _, test := range tests {
		query, err := NewQuery(test.query,
			QWithTarget(test.target),
			QWithMessageTmpl(test.fields),
			QWithIterateBy(test.iterateBy))
		if !assert.NoError(t, err) {
			continue
		}
		if !assert.Equal(t, test.queryTarget, query.query) {
			continue
		}
		results := query.Extract(test.msg)
		assert.Equal(t, tMax(1, test.countTarget), len(results))

		resultPrams := query.ParamsBy(test.msg)
		assert.Equal(t, tMax(1, test.countTarget), len(resultPrams))
		resultPrams.release()
	}
}
