// Package pipeline provides Engine that wires sources, storages and streams
// into a runnable event processing graph.
package pipeline

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/geniusrabbit/eventstream/pkg/condition"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/geniusrabbit/eventstream/stream/wrapper"
)

// Engine owns private source/storage registries and stream bindings.
// Clients register via RegisterSource / RegisterStorage / RegisterStream
// or New(...Option); registries are not part of the public API.
type Engine struct {
	Debug  bool
	Logger *zap.Logger

	sources  *sourceRegistry
	storages *storageRegistry

	mu      sync.Mutex
	streams []stream.Streamer
}

// New creates an engine with private registries and applies options in order:
// config → sources → storages → streams. Options receive the provided ctx.
func New(ctx context.Context, opts ...Option) (*Engine, error) {
	e := &Engine{
		sources:  newSourceRegistry(),
		storages: newStorageRegistry(),
		Logger:   zap.L(),
	}

	ranked := append([]Option(nil), opts...)
	sort.SliceStable(ranked, func(i, j int) bool {
		return ranked[i].phase() < ranked[j].phase()
	})

	for _, opt := range ranked {
		if err := opt.apply(ctx, e); err != nil {
			_ = e.Close()
			return nil, err
		}
	}
	if e.Logger == nil {
		e.Logger = zap.L()
	}
	return e, nil
}

// RegisterSource adds a named source to this engine.
func (e *Engine) RegisterSource(_ context.Context, name string, src eventstream.Sourcer) error {
	if name == "" {
		return fmt.Errorf("pipeline: empty source name")
	}
	if src == nil {
		return fmt.Errorf("pipeline: nil source %q", name)
	}
	e.sources.set(name, src)
	e.Logger.Info("[source] register", zap.String("source", name))
	return nil
}

// RegisterStorage adds a named storage to this engine.
func (e *Engine) RegisterStorage(_ context.Context, name string, store eventstream.Storager) error {
	if name == "" {
		return fmt.Errorf("pipeline: empty storage name")
	}
	if store == nil {
		return fmt.Errorf("pipeline: nil storage %q", name)
	}
	e.storages.set(name, store)
	e.Logger.Info("[storage] register", zap.String("store", name))
	return nil
}

// RegisterStream binds sources to a sink and starts the stream runner.
func (e *Engine) RegisterStream(ctx context.Context, name string, opts ...StreamOption) error {
	var conf streamConfig
	for _, opt := range opts {
		opt(&conf)
	}
	if name == "" {
		return fmt.Errorf("pipeline: empty stream name")
	}
	if len(conf.sources) == 0 {
		return fmt.Errorf("pipeline: stream %q has no sources", name)
	}
	if !conf.haveStorage && !conf.haveStreamer {
		return fmt.Errorf("pipeline: stream %q has no sink (storage or streamer)", name)
	}
	if conf.haveStorage && conf.haveStreamer {
		return fmt.Errorf("pipeline: stream %q: use either WithStorage or WithStreamer", name)
	}

	sink, err := e.buildSink(ctx, name, &conf)
	if err != nil {
		return err
	}

	cond := conf.whereCond
	if cond == nil && conf.where != "" {
		cond, err = condition.NewExpression(conf.where)
		if err != nil {
			return fmt.Errorf("pipeline: stream %q where: %w", name, err)
		}
	}
	var metricExec metrics.Metricer
	if len(conf.metrics) > 0 {
		metricExec, err = conf.metrics.Metric()
		if err != nil {
			return fmt.Errorf("pipeline: stream %q metrics: %w", name, err)
		}
	}
	if cond != nil || metricExec != nil {
		sink = wrapper.NewStreamWrapper(sink, cond, metricExec)
	}
	sink = wrapMapper(sink, conf.mapper)

	for _, srcName := range conf.sources {
		e.Logger.Info("[stream] subscribe",
			zap.String("stream", name),
			zap.String("source", srcName))
		if err := e.sources.subscribe(ctx, srcName, sink); err != nil {
			return fmt.Errorf("pipeline: stream %q subscribe %q: %w", name, srcName, err)
		}
	}

	e.mu.Lock()
	e.streams = append(e.streams, sink)
	e.mu.Unlock()

	e.Logger.Info("[stream] run", zap.String("stream", name))
	go func(name string, strm stream.Streamer) {
		if err := strm.Run(ctx); err != nil && ctx.Err() == nil {
			e.Logger.Error("[stream] run failed", zap.String("stream", name), zap.Error(err))
		}
	}(name, sink)
	return nil
}

func (e *Engine) buildSink(_ context.Context, name string, conf *streamConfig) (stream.Streamer, error) {
	if conf.haveStreamer {
		return conf.streamer, nil
	}
	store := e.storages.get(conf.storageName)
	if store == nil {
		return nil, fmt.Errorf("pipeline: stream %q undefined storage %q", name, conf.storageName)
	}
	return store.Stream(conf.storageOpts...)
}

// Run starts all sources and blocks until ctx is canceled.
func (e *Engine) Run(ctx context.Context) error {
	defer func() { _ = e.Close() }()
	return e.sources.listen(ctx)
}

// Close shuts down streams, sources and storages.
func (e *Engine) Close() error {
	e.mu.Lock()
	streams := e.streams
	e.streams = nil
	e.mu.Unlock()

	var firstErr error
	for _, strm := range streams {
		if err := strm.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if err := e.sources.closeAll(); err != nil && firstErr == nil {
		firstErr = err
	}
	if err := e.storages.closeAll(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}
