//
// @project geniusrabbit::eventstream 2017 - 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2021
//

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/lib/pq"
	"github.com/pkg/profile"
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/cmd/eventstream/appcontext"
	"github.com/geniusrabbit/eventstream/internal/zlogger"
	"github.com/geniusrabbit/eventstream/source"
	_ "github.com/geniusrabbit/eventstream/source/ncstreams"
	"github.com/geniusrabbit/eventstream/storage"
	_ "github.com/geniusrabbit/eventstream/storage/clickhouse"
	_ "github.com/geniusrabbit/eventstream/storage/ncstreams"
	_ "github.com/geniusrabbit/eventstream/storage/vertica"
	"github.com/geniusrabbit/eventstream/stream"
)

var (
	appVersion   string
	buildCommit  string
	buildVersion string
	buildDate    string
)

func init() {
	// Load config
	config := &appcontext.Config
	err := config.Load()
	fatalError("config.load", err)

	// Validate config
	err = config.Validate()
	fatalError("config.validate", err)

	if config.IsDebug() {
		fmt.Println("Config:", config.String())
	}

	// Init new logger object
	loggerObj, err := zlogger.New(config.ServiceName, config.LogEncoder,
		config.LogLevel, config.LogAddr, zap.Fields(
			zap.String("commit", buildCommit),
			zap.String("version", appVersion),
			zap.String("build_version", buildVersion),
			zap.String("build_date", buildDate),
		))
	fatalError("init logger", err)

	// Register global logger
	zap.ReplaceGlobals(loggerObj)
}

func main() {
	var (
		err         error
		config      = &appcontext.Config
		logger      = zap.L()
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()

	// Register stores connections
	for name, conf := range config.Stores {
		log.Printf("[storage] %s register", name)
		storageConf := &storage.Config{Debug: config.IsDebug()}
		err = conf.Decode(storageConf)
		fatalError("storage config decode <"+name+">", err)
		err = storage.Register(ctx, name, storage.WithConfig(storageConf))
		fatalError("register store <"+name+">", err)
	}

	// Register sources subscribers
	for name, conf := range config.Sources {
		log.Printf("[source] %s register", name)
		sourceConf := &source.Config{Debug: config.IsDebug()}
		err = conf.Decode(sourceConf)
		fatalError("source config decode <"+name+">", err)
		err = source.Register(ctx, name, source.WithConfig(sourceConf))
		fatalError("register source <"+name+">", err)
	}

	// Register streams
	for name, strmConf := range config.Streams {
		var (
			baseConf = &stream.Config{Name: name, Debug: config.IsDebug()}
			strm     eventstream.Streamer
		)
		if err = strmConf.Decode(baseConf); err != nil {
			fatalError("[stream] "+name+" decode error", err)
			break
		}

		if err = baseConf.Validate(); err != nil {
			fatalError("[stream] "+name+" decode error", err)
			break
		}

		if strm, err = newStream(baseConf); err != nil {
			fatalError(fmt.Sprintf("[stream] %s new init", name), err)
			return
		}

		log.Printf("[stream] %s subscribe on <%s>", name, baseConf.Source)
		if err = source.Subscribe(ctx, baseConf.Source, strm); err != nil {
			fatalError(fmt.Sprintf("[stream] "+name+" subscribe <%s>", baseConf.Source), err)
			break
		}

		log.Printf("[stream] %s run stream listener on <%s>", name, baseConf.Source)
		go func(name string) { fatalError("[stream] "+name+" run", strm.Run(ctx)) }(name)
	} // end for

	// Profiling server of collector
	runProfile(config, logger)

	// Run source listener's
	fmt.Println("> Run eventstream service")
	fatalError("profiler", source.Listen(ctx))
}

func newStream(conf *stream.Config) (eventstream.Streamer, error) {
	store := storage.Storage(conf.Store)
	if store != nil {
		return store.Stream(conf)
	}
	return nil, fmt.Errorf("[stream] %s undefined storage [%s]", conf.Name, conf.Store)
}

func runProfile(conf *appcontext.ConfigType, logger *zap.Logger) {
	switch conf.Profile.Mode {
	case "cpu":
		defer profile.Start(profile.CPUProfile).Stop()
	case "mem", "memory":
		defer profile.Start(profile.MemProfile).Stop()
	case "mutex":
		defer profile.Start(profile.MutexProfile).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile).Stop()
	case "net":
		go func() {
			fmt.Printf("Run profile (port %s)\n", conf.Profile.Listen)
			if err := http.ListenAndServe(conf.Profile.Listen, nil); err != nil {
				logger.Error("profile server error", zap.Error(err))
			}
		}()
	}
}

func fatalError(block string, err error) {
	if err != nil {
		log.Fatal("[main] fatal: ", block+" ", err)
	}
}
