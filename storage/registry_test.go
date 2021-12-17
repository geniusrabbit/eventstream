package storage

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/mocks"
)

func testConnector(ctrl *gomock.Controller) func(ctx context.Context, config *Config) (eventstream.Storager, error) {
	return func(ctx context.Context, config *Config) (eventstream.Storager, error) {
		return mocks.NewMockStorager(ctrl), nil
	}
}

func TestRegistryDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	assert.NotPanics(t, func() { RegisterConnector(`global-test`, testConnector(ctrl)) })
	assert.NoError(t, Register(ctx, `global-test`, WithConnectURL(`global-test://host`)))
	assert.Error(t, Register(ctx, `global-test-error`, WithConnectURL(`global-test-error://host`)))
	storage := Storage(`global-test`)
	if assert.NotNil(t, storage) {
		storage.(*mocks.MockStorager).EXPECT().Close().Return(nil)
		assert.NoError(t, Close())
	}
}
