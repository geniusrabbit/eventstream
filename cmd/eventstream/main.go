//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/kshvakov/clickhouse"
	_ "github.com/lib/pq"

	"github.com/geniusrabbit/eventstream/context"
	"github.com/geniusrabbit/eventstream/converter"
	"github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/eventstream/storage"
	_ "github.com/geniusrabbit/eventstream/storage/clickhouse"
	_ "github.com/geniusrabbit/eventstream/storage/hdfs"
	_ "github.com/geniusrabbit/eventstream/storage/vertica"
	"github.com/geniusrabbit/eventstream/stream"
	"github.com/geniusrabbit/eventstream/stream/clickhouse"
	"github.com/geniusrabbit/eventstream/stream/hdfs"
	"github.com/geniusrabbit/eventstream/stream/vertica"
	"github.com/geniusrabbit/notificationcenter"
)

var fabric = map[string]stream.NewConstructor{
	"clickhouse": clickhouse.New,
	"ch":         clickhouse.New,
	"vertica":    vertica.New,
	"vt":         vertica.New,
	"hdfs":       hdfs.New,
}

var (
	flagConfigFile = flag.String("config", "config.yml", "Configuration file path")
	flagDebug      = flag.Bool("debug", false, "is debug mode on")
)

func init() {
	// Parse flags
	flag.Parse()

	// Load config
	fatalError(context.Config.Load(*flagConfigFile))

	// Validate config
	fatalError(context.Config.Validate())

	// Register stures connections
	for name, st := range context.Config.Stores {
		fatalError(storage.Register(name, st.Connect, *flagDebug))
	}

	// Register sources subscribers
	for name, sr := range context.Config.Sources {
		fatalError(source.Register(name, sr.Connect))
	}
}

func main() {
	fmt.Println("> RUN APP")

	// Register streams
	for _, st := range context.Config.Streams {
		if s, err := newStream(st); nil == err {
			if err = source.Subscribe(st.Source, s); nil != err {
				notificationcenter.Close()
				fatalError(err)
				break
			}

			go test(s)
			go s.Process()
		} else {
			fatalError(err)
		}
	} // end for

	// Run notification listener
	notificationcenter.Listen()
}

func newStream(st context.StreamConfig) (s stream.ExtStreamer, err error) {
	var baseStream stream.Streamer
	if baseStream, err = newStreamBase(st); nil != err {
		return
	}

	source := context.Config.Sources[st.Source]
	s = stream.NewWrapper(baseStream, converter.ByName(source.Format))
	return
}

func newStreamBase(st context.StreamConfig) (stream.Streamer, error) {
	var (
		opt   = map[string]interface{}{}
		store = context.Config.Stores[st.Store]
	)

	if nil != st.Options {
		for k, v := range st.Options {
			opt[k] = v
		}
	}

	if nil != store.Options {
		for k, v := range store.Options {
			opt[k] = v
		}
	}

	if fb, _ := fabric[store.ConnectScheme()]; nil != fb {
		return fb(stream.Options{
			Connection: st.Store,
			RawItem:    st.RawItem,
			Target:     st.Target,
			Fields:     st.Fields,
			Options:    opt,
		})
	}

	return nil, fmt.Errorf("Undefined stream scheme: %s", store.ConnectScheme())
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func fatalError(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

func test(l stream.ExtStreamer) {
	for _ = range time.Tick(time.Second) {
		l.Handle(`{
			"srv": "service",
			"msg": "Msg",
			"err": "Error",
			"timestamp": "2017-05-05"
		}`)
	}
}
