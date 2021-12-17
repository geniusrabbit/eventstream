# Eventstream message pipeline service

[![Build Status](https://github.com/geniusrabbit/eventstream/workflows/run%20tests/badge.svg)](https://github.com/geniusrabbit/eventstream/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/eventstream)](https://goreportcard.com/report/github.com/geniusrabbit/eventstream)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/eventstream?status.svg)](https://godoc.org/github.com/geniusrabbit/eventstream)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/eventstream/badge.svg)](https://coveralls.io/github/geniusrabbit/eventstream)

Eventstream pipeline for storing and re-sending events inside the system.

> License Apache 2.0
> Copyright 2017 GeniusRabbit Dmitry Ponomarev <demdxx@gmail.com>

```sh
go get -v -u github.com/geniusrabbit/eventstream/cmd/eventstream
```

## Run eventstream service in docker

```sh
docker run -d -it -v ./custom.config.hcl:/config.hcl \
  geniusrabbit/eventstream
```

## Config example

Supports two file formats YAML & HCL

```js
stores {
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
    // Transforms into:
    //   INSERT INTO testlog (service, msg, error, timestamp) VALUES($srv, $msg, $err, @toDateTime($timestamp))
    fields = "service=srv,msg,error=err,timestamp=@toDateTime({{timestamp:date}})"
    where  = "srv == \"main\""
  }
}
```

## TODO

- [ ] Prepare evetstream as Framework extension
- [X] Add Kafka stream writer support
- [X] Add NATS stream writer support
- [X] Add Redis stream source/storage support
- [ ] Add RabbitMQ stream writer support
- [ ] Add RabbitMQ queue source support
- [ ] Add health check API
- [ ] Add metrics support (prometheus)
- [x] Add 'where' stream condition (http://github.com/Knetic/govaluate)
- [X] Ack message only if success
- [X] Buffering all data until be stored
- [ ] ~~Fix HDFS writer~~
- [X] Add support HCL config
