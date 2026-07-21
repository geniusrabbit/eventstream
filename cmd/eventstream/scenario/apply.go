package scenario

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream/pkg/pipeline"
)

// ApplyToEngine registers stores, sources and streams from a file scenario
// into a pipeline.Engine. This helper is cmd/service-only.
func ApplyToEngine(ctx context.Context, engine *pipeline.Engine, sc *Scenario) error {
	if err := sc.Validate(); err != nil {
		return err
	}
	logger := engine.Logger
	if logger == nil {
		logger = zap.L()
	}

	storeDrivers := make(map[string]string, len(sc.Stores))

	for name, item := range sc.Stores {
		logger.Info("[storage] register", zap.String("store", name))
		store, err := openStore(ctx, name, item, engine.Debug)
		if err != nil {
			return fmt.Errorf("register store <%s>: %w", name, err)
		}
		if err := engine.RegisterStorage(ctx, name, store); err != nil {
			return fmt.Errorf("register store <%s>: %w", name, err)
		}
		storeDrivers[name] = resolveDriver(item)
	}

	for name, item := range sc.Sources {
		logger.Info("[source] register", zap.String("source", name))
		src, err := openSource(ctx, name, item)
		if err != nil {
			return fmt.Errorf("register source <%s>: %w", name, err)
		}
		if err := engine.RegisterSource(ctx, name, src); err != nil {
			return fmt.Errorf("register source <%s>: %w", name, err)
		}
	}

	for name, item := range sc.Streams {
		storeName := itemString(item, "store")
		sourceName := itemString(item, "source")
		if storeName == "" || sourceName == "" {
			return fmt.Errorf("[stream] %s: store and source are required", name)
		}
		storeDriver, ok := storeDrivers[storeName]
		if !ok {
			return fmt.Errorf("[stream] %s: undefined store %q", name, storeName)
		}
		streamOpts, err := streamOptsForStore(name, storeDriver, item, engine.Debug)
		if err != nil {
			return fmt.Errorf("[stream] %s: %w", name, err)
		}
		opts := []pipeline.StreamOption{
			pipeline.WithSources(sourceName),
			pipeline.WithStorage(storeName, streamOpts...),
		}
		if where := itemString(item, "where"); where != "" {
			opts = append(opts, pipeline.WithWhere(where))
		}
		metricList, err := itemMetrics(item)
		if err != nil {
			return fmt.Errorf("[stream] %s: %w", name, err)
		}
		if len(metricList) > 0 {
			opts = append(opts, pipeline.WithMetrics(metricList))
		}
		if err := engine.RegisterStream(ctx, name, opts...); err != nil {
			return fmt.Errorf("[stream] %s: %w", name, err)
		}
	}
	return nil
}
