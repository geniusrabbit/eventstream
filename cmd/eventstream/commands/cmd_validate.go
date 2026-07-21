package commands

import (
	"context"
	"fmt"

	"github.com/geniusrabbit/blaze-api/pkg/appcmd"
	"github.com/geniusrabbit/blaze-api/pkg/context/ctxlogger"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream/cmd/eventstream/appcontext"
	"github.com/geniusrabbit/eventstream/cmd/eventstream/scenario"
)

// ValidateCommand loads and validates a scenario config without starting the pipeline.
var ValidateCommand = &appcmd.Command[appcontext.ConfigType]{
	Name:     "validate",
	HelpDesc: "Validate eventstream scenario config file",
	Exec:     validateCommand,
}

func validateCommand(ctx context.Context, _ []string, conf *appcontext.ConfigType) error {
	if conf.ConfigFile == "" {
		return fmt.Errorf("scenario config file is required (use --config)")
	}

	logger := ctxlogger.Get(ctx)
	sc, err := scenario.LoadFile(conf.ConfigFile)
	if err != nil {
		return err
	}
	if err := sc.Validate(); err != nil {
		return err
	}

	logger.Info("scenario is valid",
		zap.String("config", conf.ConfigFile),
		zap.String("version", sc.Version),
		zap.Int("stores", len(sc.Stores)),
		zap.Int("sources", len(sc.Sources)),
		zap.Int("streams", len(sc.Streams)))
	fmt.Println("OK:", conf.ConfigFile)
	return nil
}
