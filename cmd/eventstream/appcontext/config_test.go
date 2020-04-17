package appcontext

import (
	"os"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func init() {
	os.Args = os.Args[:1]
}

func TestLoad(t *testing.T) {
	assert.NoError(t, Config.Load())
	assert.True(t, len(Config.String()) > 2)
	assert.Error(t, Config.Validate())
	assert.True(t, Config.IsDebug())
}
