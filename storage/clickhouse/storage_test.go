package clickhouse

import (
	"context"
	"testing"

	"github.com/geniusrabbit/eventstream/storage/sql"
	"github.com/stretchr/testify/assert"
)

func TestOpenError(t *testing.T) {
	storage, err := Open(context.Background(), `error`,
		WithQuery(sql.QWithTarget(`test`)))
	if !assert.Error(t, err) {
		assert.NoError(t, storage.Close())
	}
}
