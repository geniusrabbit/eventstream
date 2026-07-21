package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/geniusrabbit/blaze-api/pkg/zlogger"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/pkg/pipeline"
	"github.com/geniusrabbit/eventstream/source/nats"
	"github.com/geniusrabbit/eventstream/storage/clickhouse"
	"github.com/geniusrabbit/eventstream/storage/sql"
)

var (
	commit     string
	appVersion string
)

type config struct {
	LogLevel       string
	LogAddr        string
	LogEncoder     string
	SourceConnect  string
	StorageConnect string
}

func loadConfig() config {
	return config{
		LogLevel:       envOr("LOGGER_LEVEL", "debug"),
		LogAddr:        os.Getenv("LOG_ADDR"),
		LogEncoder:     os.Getenv("LOG_ENCODER"),
		SourceConnect:  os.Getenv("SOURCE_CONNECT"),
		StorageConnect: os.Getenv("STORAGE_CONNECT"),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type logMessage struct {
	Time    time.Time `field:"time" json:"time"`
	UUID    string    `field:"uuid" type:"uuid" json:"uuid"`
	Type    string    `field:"type" json:"type"`
	Message string    `field:"message" json:"message"`
	UserID  uint64    `field:"user_id" json:"user_id"`
	Data    string    `field:"data" json:"data"`
}

func main() {
	conf := loadConfig()

	logger, err := zlogger.New("example", conf.LogEncoder, conf.LogLevel, conf.LogAddr, zap.Fields(
		zap.String("commit", commit),
		zap.String("version", appVersion),
	))
	fatalError("init logger", err)
	zap.ReplaceGlobals(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	engine, err := pipeline.New(ctx,
		pipeline.WithLogger(logger),
		pipeline.WithDebug(conf.LogLevel == "debug"),
		pipeline.SourceFunc("global", func(ctx context.Context) (eventstream.Sourcer, error) {
			return nats.Open(ctx, conf.SourceConnect)
		}),
		pipeline.StorageFunc("clickhouse", func(ctx context.Context) (eventstream.Storager, error) {
			return clickhouse.Open(ctx, conf.StorageConnect)
		}),
		pipeline.Stream("logs",
			pipeline.WithSources("global"),
			pipeline.WithStorage("clickhouse", clickhouse.WithQuery(
				sql.QWithTarget(`logs.common`),
				sql.QWithMessageTmplOf[logMessage](),
			)),
		),
	)
	fatalError("pipeline", err)

	logger.Info("RUN")
	fatalError("run", engine.Run(ctx))
}

func fatalError(block string, err error) {
	if err != nil {
		log.Fatal("[main] fatal: ", block+" ", err)
	}
}
