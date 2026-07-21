package scenario_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geniusrabbit/eventstream/cmd/eventstream/scenario"
)

const sampleYAML = `
version: "1"
description: test
stores:
  ping:
    driver: ping
    connect: "http://localhost:1"
sources:
  nats_1:
    driver: nats
    connect: "nats://localhost:4222/group?topics=t"
    format: json
streams:
  log_1:
    store: ping
    source: nats_1
    target: unused
`

const sampleTOML = `
version = "1"
description = "test"

[stores.ping]
driver = "ping"
connect = "http://localhost:1"

[sources.nats_1]
driver = "nats"
connect = "nats://localhost:4222/group?topics=t"
format = "json"

[streams.log_1]
store = "ping"
source = "nats_1"
target = "unused"
`

const sampleHCL = `
version = "1"
description = "test"

stores {
  ping {
    driver  = "ping"
    connect = "http://localhost:1"
  }
}

sources {
  nats_1 {
    driver  = "nats"
    connect = "nats://localhost:4222/group?topics=t"
    format  = "json"
  }
}

streams {
  log_1 {
    store  = "ping"
    source = "nats_1"
    target = "unused"
  }
}
`

func TestLoadFormats(t *testing.T) {
	cases := []struct {
		name string
		ext  string
		body string
	}{
		{"yaml", ".yml", sampleYAML},
		{"toml", ".toml", sampleTOML},
		{"hcl", ".hcl", sampleHCL},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sc, err := scenario.Load([]byte(tc.body), tc.ext)
			require.NoError(t, err)
			require.NoError(t, sc.Validate())
			assert.Equal(t, "1", sc.Version)
			assert.Contains(t, sc.Stores, "ping")
			assert.Contains(t, sc.Sources, "nats_1")
			assert.Contains(t, sc.Streams, "log_1")
			assert.Equal(t, "ping", sc.Streams["log_1"]["store"])
		})
	}
}

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scene.yml")
	require.NoError(t, os.WriteFile(path, []byte(sampleYAML), 0o600))
	sc, err := scenario.LoadFile(path)
	require.NoError(t, err)
	assert.NoError(t, sc.Validate())
}

func TestValidateEmpty(t *testing.T) {
	assert.ErrorIs(t, (&scenario.Scenario{}).Validate(), scenario.ErrEmptyStores)
}

func TestLoadDevelopConfigHCL(t *testing.T) {
	sc, err := scenario.LoadFile("../../../deploy/develop/config.hcl")
	require.NoError(t, err)
	require.NoError(t, sc.Validate())
	assert.Contains(t, sc.Stores, "clickhouse_1")
	assert.Contains(t, sc.Sources, "nats_1")
	assert.GreaterOrEqual(t, len(sc.Streams), 3)
}
