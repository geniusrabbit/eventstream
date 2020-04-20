package vertica

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenError(t *testing.T) {
	storage, err := Open(`error`, WithDebug(true))
	if !assert.Error(t, err) {
		assert.NoError(t, storage.Close())
	}
}
