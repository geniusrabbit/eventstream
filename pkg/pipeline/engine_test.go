package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geniusrabbit/eventstream/internal/mocks"
	"github.com/geniusrabbit/eventstream/pkg/message"
	"github.com/geniusrabbit/eventstream/pkg/pipeline"
	"github.com/geniusrabbit/eventstream/stream"
)

type stubStreamer struct {
	id   string
	puts []message.Message
}

func (s *stubStreamer) ID() string { return s.id }
func (s *stubStreamer) Put(_ context.Context, msg message.Message) error {
	s.puts = append(s.puts, msg)
	return nil
}
func (s *stubStreamer) Check(context.Context, message.Message) bool { return true }
func (s *stubStreamer) Run(context.Context) error                   { return nil }
func (s *stubStreamer) Close() error                                { return nil }

func TestNewWithOptionsAndRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	src := mocks.NewMockSourcer(ctrl)
	store := mocks.NewMockStorager(ctrl)
	sink := &stubStreamer{id: "sink-1"}

	src.EXPECT().Subscribe(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, streams ...stream.Streamer) error {
			require.Len(t, streams, 1)
			return nil
		},
	)
	src.EXPECT().Start(gomock.Any()).Return(nil)
	src.EXPECT().Close().Return(nil)
	store.EXPECT().Close().Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine, err := pipeline.New(ctx,
		pipeline.Source("s1", src),
		pipeline.Storage("db", store),
		pipeline.Stream("st1",
			pipeline.WithSources("s1"),
			pipeline.WithStreamer(sink),
		),
	)
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() { done <- engine.Run(ctx) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	require.NoError(t, <-done)
}

func TestRegisterStreamRequiresSource(t *testing.T) {
	ctx := context.Background()
	engine, err := pipeline.New(ctx)
	require.NoError(t, err)
	err = engine.RegisterStream(ctx, "x", pipeline.WithStreamer(&stubStreamer{id: "s"}))
	assert.Error(t, err)
}
