package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenError(t *testing.T) {
	storage, err := Open(`error`, WithQueryByTarget(`test`, nil))
	if !assert.Error(t, err) {
		assert.NoError(t, storage.Close())
	}
}
