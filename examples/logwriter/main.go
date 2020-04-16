package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/geniusrabbit/eventstream/source/ncstreams"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/clickhouse"
)

var (
	commit     string
	appVersion string
)

type config struct {
	LogLevel       string `env:"LOGGER_LEVEL"`
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

}

func commandLogwriter() {
	var conf config

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init new logger object
	logger, err := newLogger(conf.isDebug(), conf.LogLevel, zap.Fields(
		zap.String("commit", commit),
		zap.String("version", appVersion),
	))

	fatalError("init logger", err)
	logger.Info(`RUN`)

	// Open source connection
	datasource, err := ncstreams.Open(conf.SourceConnect)
	fatalError(`source connect`, err)

	// Open clickhouse storage connect
	datastorage, err := clickhouse.Open(conf.StorageConnect, storage.WithDebug(conf.isDebug()))
	fatalError(`clickhouse storage connect`, err)

	// Get new stream writer for the storage
	stream, err := datastorage.Stream(clickhouse.WithQueryByTarget(`logs.common`, &logMessage{}))
	fatalError(`get new writer stream`, err)

	// Subscribe stream writer of the clickhouse
	err = datasource.Subscribe(ctx, stream)
	fatalError(`subscribe storage target`, err)

	err = datasource.Start(ctx)
	fatalError(`start source`, err)

	<-ctx.Done()
}

func newLogger(debug bool, loglevel string, options ...zap.Option) (logger *zap.Logger, err error) {
	if debug {
		return zap.NewDevelopment(options...)
	}
	var (
		level         zapcore.Level
		loggerEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		})
	)
	if err := level.UnmarshalText([]byte(loglevel)); err != nil {
		logger.Error("parse log level error", zap.Error(err))
	}
	core := zapcore.NewCore(loggerEncoder, os.Stdout, level)
	logger = zap.New(core, options...)

	return logger, nil
}

func fatalError(block string, err error) {
	if err != nil {
		log.Fatal("[main] fatal: ", block+" ", err)
	}
}
