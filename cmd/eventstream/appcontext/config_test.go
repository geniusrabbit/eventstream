package appcontext

import (
	"os"
	"testing"

	"github.com/demdxx/goconfig"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Args = os.Args[:1]
}

func TestLoad(t *testing.T) {
	var conf ConfigType
	assert.NoError(t, goconfig.Load(&conf))
	assert.True(t, len(conf.String()) > 2)
	assert.Equal(t, "eventstream", conf.ServiceName)
	assert.True(t, conf.IsDebug())
}
