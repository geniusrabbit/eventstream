package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/geniusrabbit/notificationcenter/v2/nats"
)

var (
	flagStreamConnect = flag.String(`stream`, ``, `scheme://host:port/group?topics=name`)
	flagMessage       = flag.String(`message`, ``, `{"type":"event","message":"OK"}`)
	flagRepeat        = flag.Int(`repeat`, 1, `"--repeat=10" times`)
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pub, err := nats.NewPublisher(nats.WithNatsURL(*flagStreamConnect))
	fatalError(`publisher`, err)

	for i := 0; i < *flagRepeat; i++ {
		fmt.Println("> SEND", *flagMessage)
		fatalError(`publish`, pub.Publish(ctx, json.RawMessage(*flagMessage)))
		time.Sleep(time.Second)
	}
}

func fatalError(block string, err error) {
	if err != nil {
		log.Fatal("[main] fatal: ", block+" ", err)
	}
}
