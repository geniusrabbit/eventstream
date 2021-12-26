package main

import (
	"context"
	"log"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream/internal/zlogger"
	"github.com/geniusrabbit/eventstream/source"
	_ "github.com/geniusrabbit/eventstream/source/ncstreams"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/clickhouse"
)

var (
	commit     string
	appVersion string
)

type config struct {
	LogLevel       string `env:"LOGGER_LEVEL"`
	LogAddr        string `env:"LOG_ADDR"`
	LogEncoder     string `env:"LOG_ENCODER"`
	SourceConnect  string `env:"SOURCE_CONNECT"`
	StorageConnect string `env:"STORAGE_CONNECT"`
}

func (cnf *config) isDebug() bool {
	return cnf.LogLevel == `debug`
}

type logMessage struct {
	Time    time.Time `field:"time"`
	UUID    string    `field:"uuid" type:"uuid"`
	Type    string    `field:"type"`
	Message string    `field:"message"`
	UserID  uint64    `field:"user_id"`
	Data    string    `field:"data"`
}

func main() {
	var conf config

	// Init new logger object
	logger, err := zlogger.New("example", conf.LogEncoder, conf.LogLevel, conf.LogAddr, zap.Fields(
		zap.String("commit", commit),
		zap.String("version", appVersion),
	))
	fatalError("init logger", err)

	commandLogwriter(&conf, logger)
}

func commandLogwriter(conf *config, logger *zap.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info(`RUN`)

	// Open source connection
	err := source.Register(ctx, "global", source.WithConnectURL(conf.SourceConnect))
	fatalError(`source connect`, err)

	// Open clickhouse storage connect
	datastorage, err := clickhouse.Open(conf.StorageConnect, storage.WithDebug(conf.isDebug()))
	fatalError(`clickhouse storage connect`, err)

	// Get new stream writer for the storage
	stream, err := datastorage.Stream(clickhouse.WithQueryByTarget(`logs.common`, &logMessage{}))
	fatalError(`get new writer stream`, err)

	// Subscribe stream writer of the clickhouse
	err = source.Subscribe(ctx, "global", stream)
	fatalError(`subscribe storage target`, err)

	err = source.Listen(ctx)
	fatalError(`start source`, err)

	<-ctx.Done()
}

func fatalError(block string, err error) {
	if err != nil {
		log.Fatal("[main] fatal: ", block+" ", err)
	}
}
