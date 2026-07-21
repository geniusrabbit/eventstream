package scenario

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLStreamOptsRequiresQueryOrTarget(t *testing.T) {
	_, err := streamOptsForStore("x", "clickhouse", Item{"store": "ch"}, false)
	assert.Error(t, err)
}

func TestParseTargetFieldsMap(t *testing.T) {
	fields, err := parseTargetFields(map[string]any{
		"type":    "message",
		"service": "{{service}}",
	})
	require.NoError(t, err)
	assert.Equal(t, "message", fields["type"])
	assert.Equal(t, "{{service}}", fields["service"])
}

func TestParseTargetFieldsHCLBlock(t *testing.T) {
	// hashicorp/hcl v1 encodes nested blocks as arrays of objects.
	fields, err := parseTargetFields([]any{
		map[string]any{
			"type":    "message",
			"service": "{{service}}",
			"message": "{{msg}}",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"type":    "message",
		"service": "{{service}}",
		"message": "{{msg}}",
	}, fields)
}

func TestParseStreamTargetsHCLShape(t *testing.T) {
	targets, err := parseStreamTargets([]any{
		map[string]any{
			"fields": []any{
				map[string]any{"type": "message", "service": "{{service}}"},
			},
			"where": "!error",
		},
		map[string]any{
			"fields": []any{
				map[string]any{"type": "error", "error": "{{msg}}"},
			},
			"where": "error",
		},
	})
	require.NoError(t, err)
	require.Len(t, targets, 2)
	assert.Equal(t, "message", targets[0].Fields["type"])
	assert.Equal(t, "!error", targets[0].Where)
	assert.Equal(t, "error", targets[1].Fields["type"])
}
