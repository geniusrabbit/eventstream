//
// @project geniusrabbit::eventstream 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/kshvakov/clickhouse"
	_ "github.com/lib/pq"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/context"
	"github.com/geniusrabbit/eventstream/source"
	_ "github.com/geniusrabbit/eventstream/source/kafka"
	_ "github.com/geniusrabbit/eventstream/source/nats"
	"github.com/geniusrabbit/eventstream/storage"
	_ "github.com/geniusrabbit/eventstream/storage/clickhouse"
	_ "github.com/geniusrabbit/eventstream/storage/metrics"
	_ "github.com/geniusrabbit/eventstream/storage/vertica"
)

var (
	flagConfigFile = flag.String("config", "config.hcl", "Configuration file path")
	flagDebug      = flag.Bool("debug", false, "is debug mode on")
)

func init() {
	// Parse flags
	flag.Parse()

	// Load config
	fatalError("config.load", context.Config.Load(*flagConfigFile))

	// Validate config
	fatalError("config.validate", context.Config.Validate())

	if *flagDebug {
		context.Config.Debug = *flagDebug
		fmt.Println("Config:", context.Config.String())
	}

	// Register stores connections
	for name, conf := range context.Config.Stores {
		storageConf := &storage.Config{Debug: context.Config.Debug}
		fatalError("storage config decode <"+name+">", conf.Decode(storageConf))
		fatalError("register store <"+name+">", storage.Register(name, storageConf))
	}

	// Register sources subscribers
	for name, conf := range context.Config.Sources {
		sourceConf := &source.Config{Debug: context.Config.Debug}
		fatalError("source config decode <"+name+">", conf.Decode(sourceConf))
		fatalError("register source <"+name+">", source.Register(name, sourceConf))
	}
}

func main() {
	fmt.Println("> Run eventstream service")

	var err error

	// Register streams
	for name, strmConf := range context.Config.Streams {
		var (
			baseConf = &storage.StreamConfig{Name: name, Debug: context.Config.Debug}
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
		if err = source.Subscribe(baseConf.Source, strm); nil != err {
			fatalError(fmt.Sprintf("[stream] "+name+" subscribe <%s>", baseConf.Source), err)
			break
		}

		log.Printf("[stream] %s run stream listener on <%s>", name, baseConf.Source)
		go func(name string) { fatalError("[strean] "+name+" run", strm.Run()) }(name)
	} // end for

	// Run source listener's
	source.Listen()
	close()
}

func newStream(conf *storage.StreamConfig) (eventstream.Streamer, error) {
	store := storage.Storage(conf.Store)
	if store != nil {
		return store.Stream(conf)
	}
	return nil, fmt.Errorf("[stream] %s undefined storage [%s]", conf.Name, conf.Store)
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func close() {
	source.Close()
	storage.Close()
}

func fatalError(block string, err error) {
	if err != nil {
		close()
		log.Fatal("[main] fatal: ", block+" ", err)
	}
}
