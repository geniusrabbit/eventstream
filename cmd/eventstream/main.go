// Command eventstream runs the file-driven pipeline service
// (server / validate) with tag-selected source and storage drivers.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/demdxx/goconfig"
	"github.com/geniusrabbit/blaze-api/pkg/appcmd"
	"github.com/geniusrabbit/blaze-api/pkg/context/ctxlogger"
	"github.com/geniusrabbit/blaze-api/pkg/zlogger"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream/cmd/eventstream/appcontext"
	"github.com/geniusrabbit/eventstream/cmd/eventstream/commands"

	_ "github.com/geniusrabbit/eventstream/cmd/eventstream/drivers"
)

var (
	buildCommit  string
	buildVersion string
	buildDate    string
)

const bannerString = `
███████ ██    ██ ███████ ███    ██ ████████ ███████ ████████ ██████  ███████  █████  ███    ███
██      ██    ██ ██      ████   ██    ██    ██         ██    ██   ██ ██      ██   ██ ████  ████
█████   ██    ██ █████   ██ ██  ██    ██    ███████    ██    ██████  █████   ███████ ██ ████ ██
██       ██  ██  ██      ██  ██ ██    ██         ██    ██    ██   ██ ██      ██   ██ ██  ██  ██
███████   ████   ███████ ██   ████    ██    ███████    ██    ██   ██ ███████ ██   ██ ██      ██

`

func main() {
	fmt.Println(bannerString + buildVersion)

	var boot appcontext.ConfigType
	if err := goconfig.Load(&boot); err != nil {
		fatalError(err, "config.load")
	}

	if boot.IsDebug() {
		fmt.Println("Config:", boot.String())
	}

	loggerObj, err := zlogger.New(boot.ServiceName, boot.LogEncoder,
		boot.LogLevel, boot.LogAddr, zap.Fields(
			zap.String("commit", buildCommit),
			zap.String("build_version", buildVersion),
			zap.String("build_date", buildDate),
		))
	fatalError(err, "init logger")
	zap.ReplaceGlobals(loggerObj)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := &appcmd.App{
		Name:        "eventstream",
		Description: "Eventstream pipeline for storing and re-sending events",
		Version:     buildVersion,
		BuildCommit: buildCommit,
		BuildDate:   buildDate,
		CmdList: appcmd.ICommands{
			commands.ServerCommand,
			commands.ValidateCommand,
		},
		BeforeCommandRun: func(ctx context.Context, cmd appcmd.ICommand) (context.Context, error) {
			fmt.Println()
			fmt.Println("░█ Log Level:\x1b[32m", boot.LogLevel, "\x1b[0m")
			fmt.Println("░█ Run command:\x1b[31m", cmd.Cmd(), "\x1b[0m")
			return ctxlogger.WithLogger(ctx, zap.L()), nil
		},
	}

	fatalError(app.Run(ctx, os.Args), "application run")
}

func fatalError(err error, msgs ...any) {
	if err != nil {
		zap.L().Fatal(fmt.Sprint(msgs...), zap.Error(err))
	}
}
