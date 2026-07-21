// Package commands implements CLI commands for the eventstream service binary.
package commands

import (
	"context"
	"fmt"

	"github.com/geniusrabbit/blaze-api/pkg/appcmd"
	"github.com/geniusrabbit/blaze-api/pkg/context/ctxlogger"
	"github.com/geniusrabbit/blaze-api/pkg/profiler"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream/cmd/eventstream/appcontext"
	"github.com/geniusrabbit/eventstream/cmd/eventstream/scenario"
	"github.com/geniusrabbit/eventstream/pkg/pipeline"
)

// ServerCommand runs the eventstream pipeline service.
var ServerCommand = &appcmd.Command[appcontext.ConfigType]{
	Name:     "server",
	HelpDesc: "Run eventstream pipeline server",
	Exec:     serverCommand,
}

func serverCommand(ctx context.Context, _ []string, conf *appcontext.ConfigType) error {
	if conf.ConfigFile == "" {
		return fmt.Errorf("scenario config file is required (use --config)")
	}

	logger := ctxlogger.Get(ctx)
	profiler.Run(conf.Profile.Mode, conf.Profile.Listen, logger, true)

	sc, err := scenario.LoadFile(conf.ConfigFile)
	if err != nil {
		return err
	}

	engine, err := pipeline.New(ctx,
		pipeline.WithDebug(conf.IsDebug()),
		pipeline.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	if err := scenario.ApplyToEngine(ctx, engine, sc); err != nil {
		return err
	}

	logger.Info("starting eventstream server",
		zap.String("config", conf.ConfigFile),
		zap.String("version", sc.Version))

	fmt.Println("> Run eventstream service")
	return engine.Run(ctx)
}
