//
// @project geniusrabbit::eventstream 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
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
	_ "github.com/geniusrabbit/eventstream/storage/hdfs"
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
	fatalError(context.Config.Load(*flagConfigFile))

	// Validate config
	fatalError(context.Config.Validate())

	if *flagDebug {
		context.Config.Debug = *flagDebug
		fmt.Println("Config:", context.Config.String())
	}

	// Register stores connections
	for name, conf := range context.Config.Stores {
		fatalError(storage.Register(name, conf, *flagDebug))
	}

	// Register sources subscribers
	for name, conf := range context.Config.Sources {
		fatalError(source.Register(name, conf, *flagDebug))
	}
}

func main() {
	fmt.Println("> RUN APP")

	// Register streams
	for _, strmConf := range context.Config.Streams {
		if strm, err := newStream(strmConf); err == nil {
			sourceName := strmConf.String("source", "")
			if err = source.Subscribe(sourceName, strm); nil != err {
				fatalError(err)
				break
			}

			go func() {
				if err = strm.Run(); err != nil {
					fatalError(err)
					return
				}
			}()
		} else {
			fatalError(err)
			return
		}
	} // end for

	// Run source listener's
	source.Listen()
	close()
}

func newStream(conf eventstream.ConfigItem) (eventstream.Streamer, error) {
	store := storage.Storage(conf.String("store", ""))
	if store != nil {
		return store.Stream(conf)
	}
	return nil, fmt.Errorf("Undefined storage [%s]", conf.String("store", ""))
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func close() {
	source.Close()
	storage.Close()
}

func fatalError(err error) {
	if err != nil {
		defer log.Fatal("[main] fatal:", err)
		close()
	}
}
