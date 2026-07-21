package vertica

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenError(t *testing.T) {
	storage, err := Open(context.Background(), `:::bad`, WithDebug(true))
	if !assert.Error(t, err) {
		assert.NoError(t, storage.Close())
	}
}
