package clickhouse

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenCreatesHandle(t *testing.T) {
	// database/sql Open is lazy; DSN validity is checked on Ping/use.
	storage, err := Open(context.Background(), `clickhouse://localhost:9000/default`)
	assert.NoError(t, err)
	if assert.NotNil(t, storage) {
		assert.NoError(t, storage.Close())
	}
}

func TestCloseNilConn(t *testing.T) {
	assert.NoError(t, (&Clickhouse{}).Close())
}

func TestConnectionConcurrent(t *testing.T) {
	t.Parallel()

	storage := &Clickhouse{connect: `clickhouse://127.0.0.1:1/default`}
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := storage.Connection()
			if err == nil {
				assert.NotNil(t, conn)
			}
		}()
	}
	wg.Wait()
}

func TestConnectionStaleReconnect(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("clickhouse", `clickhouse://127.0.0.1:1/default`)
	require.NoError(t, err)
	storage := &Clickhouse{connect: `clickhouse://127.0.0.1:1/default`, conn: db}

	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = storage.Connection()
		}()
	}
	wg.Wait()
}
