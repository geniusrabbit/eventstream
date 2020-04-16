//
// @project geniusrabbit::eventstream 2017 - 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2020
//

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/lib/pq"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/cmd/eventstream/appcontext"
	"github.com/geniusrabbit/eventstream/source"
	_ "github.com/geniusrabbit/eventstream/source/ncstreams"
	"github.com/geniusrabbit/eventstream/storage"
	_ "github.com/geniusrabbit/eventstream/storage/clickhouse"
	_ "github.com/geniusrabbit/eventstream/storage/ncstreams"
	_ "github.com/geniusrabbit/eventstream/storage/vertica"
	"github.com/geniusrabbit/eventstream/stream"
)

var (
	flagConfigFile = flag.String("config", "config.hcl", "Configuration file path")
	flagProfiler   = flag.String("profiler", "", "The hostname and port of golang profiler, for example: :6060")
)

func init() {
	// Parse flags
	flag.Parse()

	// Load config
	config := &appcontext.Config
	err := config.Load(*flagConfigFile)
	fatalError("config.load", err)

	// Validate config
	err = config.Validate()
	fatalError("config.validate", err)

	if config.IsDebug() {
		fmt.Println("Config:", config.String())
	}

	// Register stores connections
	for name, conf := range config.Stores {
		log.Printf("[storage] %s register", name)
		storageConf := &storage.Config{Debug: config.IsDebug()}
		err = conf.Decode(storageConf)
		fatalError("storage config decode <"+name+">", err)
		err = storage.Register(name, storage.WithConfig(storageConf))
		fatalError("register store <"+name+">", err)
	}

	// Register sources subscribers
	for name, conf := range config.Sources {
		log.Printf("[source] %s register", name)
		sourceConf := &source.Config{Debug: config.IsDebug()}
		err = conf.Decode(sourceConf)
		fatalError("source config decode <"+name+">", err)
		err = source.Register(name, source.WithConfig(sourceConf))
		fatalError("register source <"+name+">", err)
	}
}

func main() {
	var (
		err         error
		config      = &appcontext.Config
		ctx, cancel = context.WithCancel(context.Background())
	)

	defer cancel()

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

	// Run profiler
	if *flagProfiler != "" {
		go func() {
			fmt.Println("Run profile: " + *flagProfiler)
			fatalError("profiler", http.ListenAndServe(*flagProfiler, nil))
		}()
	}

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

func fatalError(block string, err error) {
	if err != nil {
		log.Fatal("[main] fatal: ", block+" ", err)
	}
}
