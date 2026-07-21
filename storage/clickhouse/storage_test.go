package clickhouse

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenCreatesHandle(t *testing.T) {
	// database/sql Open is lazy; DSN validity is checked on Ping/use.
	storage, err := Open(context.Background(), `clickhouse://localhost:9000/default`)
	assert.NoError(t, err)
	if assert.NotNil(t, storage) {
		assert.NoError(t, storage.Close())
	}
}
