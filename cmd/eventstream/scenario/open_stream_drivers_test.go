//go:build all || allsource || allstorage || clickhouse || vertica || nats || kafka || natstream || redistream
// +build all allsource allstorage clickhouse vertica nats kafka natstream redistream

package scenario

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/geniusrabbit/eventstream/cmd/eventstream/drivers"
)

func TestSQLStreamOptsFromTargetFields(t *testing.T) {
	item := Item{
		"store":  "clickhouse_1",
		"source": "nats_1",
		"target": "stat.testlog",
		"fields": []any{"service=srv", "msg"},
	}
	opts, err := streamOptsForStore("log_3", "clickhouse", item, true)
	require.NoError(t, err)
	assert.NotEmpty(t, opts)
}

func TestNCStreamOptsFromTargets(t *testing.T) {
	item := Item{
		"store":  "nats_store",
		"source": "nats_1",
		"targets": []any{
			map[string]any{
				"fields": map[string]any{"type": "message"},
				"where":  "!error",
			},
		},
	}
	opts, err := streamOptsForStore("nats_1", "nats", item, false)
	require.NoError(t, err)
	assert.Len(t, opts, 2)
}
