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

## Source list

- **kafka**
- **NATS** & **NATS stream**
- **Redis** stream

## Storage list

- **Clickhouse**
- **Vertica**
- **kafka**
- **NATS**
- **Redis** stream

## Config example

Supports two file formats YAML & HCL

```js
stores {
  clickhouse_1 {
    connect = "@env:CLICKHOUSE_STORE_CONNECT"
    options { # Optional
      buffer = 1000
    }
  }
  kafka_1 {
    connect = "@env:KAFKA_EVENTS_CONNECT"
  }
}

// Source could be any supported stream service like kafka, nats, etc...
sources {
  nats_1 {
    connect = "@env:NATS_SOURCE_CONNECT"
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
    metrics = [
      {
        name = "log.counter"
        type = "counter"
        tags {
          server  = "{{srv}}"
        }
      }
    ]
  }
  kafka_retranslate {
    store  = "kafka_1"
    source = "nats_1"
    targets = [
      {
        fields = {
          server = "{{srv}}"
          timestamp = "{{timestamp}}"
        }
        where = "type = \"statistic\""
      }
    ]
    where = "srv = \"events\""
  }
}
```

## Metrics

Metrics helps analyze some events during processing and monitor streams state.
Every stream can process metrics with the keyword `metrics`.

Example:
```js
metrics = [
  {
    name = "log.counter"
    type = "counter"
    tags { server = "{{srv}}" }
  },
  {
    name = "actions.counter"
    type = "counter"
    tags { action = "{{action}}" }
  },
  {...}
]
```

All metrics available by URL `/metrics` with prometheus protocol.
To activate metrics need to define profile connection port.

```env
SERVER_PROFILE_MODE=net
SERVER_PROFILE_LISTEN=:6060
```

## TODO

- [ ] Add MySQL database storage
- [ ] Add PostgreSQL database storage
- [ ] Add MongoDB database storage
- [ ] Add Redis database storage
- [X] Prepare evetstream as Framework extension
- [X] Add Kafka stream writer support
- [X] Add NATS stream writer support
- [X] Add Redis stream source/storage support
- [ ] Add RabbitMQ stream source/storage support
- [ ] Add health check API
- [X] Add customizable prometheus metrics
- [x] Add 'where' stream condition (http://github.com/Knetic/govaluate)
- [X] Ack message only if success
- [X] Buffering all data until be stored
- [X] Add support HCL config
