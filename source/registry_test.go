package source

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/mocks"
)

func testConnector(ctrl *gomock.Controller) func(config *Config) (eventstream.Sourcer, error) {
	return func(config *Config) (eventstream.Sourcer, error) {
		return mocks.NewMockSourcer(ctrl), nil
	}
}

func TestRegistryDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert.NotPanics(t, func() { RegisterConnector(testConnector(ctrl), `global-test`) })
	assert.NoError(t, Register(`global-test`, WithConnectURL(`global-test://host`)))
	assert.Error(t, Register(`global-test-error`, WithConnectURL(`global-test-error://host`)))
	source := Source(`global-test`)
	if assert.NotNil(t, source) {
		source.(*mocks.MockSourcer).EXPECT().Close().Return(nil)
		assert.NoError(t, Close())
	}
}
