# Eventstream

[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/eventstream)](https://goreportcard.com/report/github.com/geniusrabbit/eventstream)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/eventstream?status.svg)](https://godoc.org/github.com/geniusrabbit/eventstream)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/eventstream/badge.svg)](https://coveralls.io/github/geniusrabbit/eventstream)

Eventstream pipeline for storing and re-sending events inside the system.

> License Apache 2.0 
> Copyright 2017 GeniusRabbit Dmitry Ponomarev <demdxx@gmail.com>

```sh
go get -v -u -t github.com/geniusrabbit/eventstream/cmd/eventstream
```

## Config

Supports two file formats YAML & HCL

```js
stores {
  // CREATE TABLE stat.testlog (
  //    timestamp        DateTime
  //  , datemark         Date default toDate(timestamp)
  //  , service          String
  //  , msg              String
  //  , error            String
  //  , created_at       DateTime default now()
  // ) Engine=MergeTree(datemark, (service), 8192);

  clickhouse_1 {
    connect = "clickhouse://clickhouse:9000/stat"
    options { # Optional
      buffer = 1000
    }
  }
}

// Source could be any supported stream service like kafka, nats, etc...
sources {
  nats_1 {
    connect = "nats://nats:4222/?topics=topic1,topic2"
    format  = "json"
  }
}

// Streams it's pipelines which have source and destination store
streams {
  log_1 {
    store  = "clickhouse_1"
    source = "nats_1"
    target = "testlog"
    // Optional if fields in log and in message the same
    fields = "service=srv,msg,error=err,timestamp=@toTimestamp({{timestamp:date}})"
  }
}
```

## TODO

 - [ ] Ack message only if success
 - [ ] Buffering all data until be stored
 - [ ] Fix HDFS writer
 - [x] Add support HCL config
