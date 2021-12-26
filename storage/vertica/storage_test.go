package vertica

import (
	"testing"

	sqlstorage "github.com/geniusrabbit/eventstream/storage/sql"
	"github.com/stretchr/testify/assert"
)

func TestOpenError(t *testing.T) {
	storage, err := Open(`error`, sqlstorage.WithDebug(true))
	if !assert.Error(t, err) {
		assert.NoError(t, storage.Close())
	}
}
