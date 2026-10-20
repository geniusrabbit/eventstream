# Eventstream message pipeline

![License](https://img.shields.io/github/license/geniusrabbit/eventstream)
[![Docker Pulls](https://img.shields.io/docker/pulls/geniusrabbit/eventstream.svg?maxAge=604800)](https://hub.docker.com/r/geniusrabbit/eventstream)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/eventstream)](https://goreportcard.com/report/github.com/geniusrabbit/eventstream)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/eventstream/badge.svg?branch=master)](https://coveralls.io/github/geniusrabbit/eventstream?branch=master)
[![Testing Status](https://github.com/geniusrabbit/eventstream/workflows/Tests/badge.svg)](https://github.com/geniusrabbit/eventstream/actions?workflow=Tests)
[![Publish Docker Status](https://github.com/geniusrabbit/eventstream/workflows/Publish/badge.svg)](https://github.com/geniusrabbit/eventstream/actions?workflow=Publish)

Eventstream is a Go **framework** for building event pipelines (source → map → storage),
plus an optional **service binary** that loads file-based scenarios.

```sh
go get github.com/geniusrabbit/eventstream
```

## Use as a library

The library API is code-first: `pipeline.New(ctx, ...Option)` with private registries
(no global connection registry required).

```go
engine, err := pipeline.New(ctx,
    pipeline.WithLogger(log),
    pipeline.SourceFunc("nats", func(ctx context.Context) (eventstream.Sourcer, error) {
        return nats.Open(ctx, natsURL)
    }),
    pipeline.StorageFunc("ch", func(ctx context.Context) (eventstream.Storager, error) {
        return clickhouse.Open(ctx, dsn)
    }),
    pipeline.Stream("logs",
        pipeline.WithSources("nats"),
        pipeline.WithStorage("ch", clickhouse.WithQuery(
            sql.QWithTarget("logs.common"),
            sql.QWithMessageTmplOf[LogRow](),
        )),
        pipeline.WithMapper(pipeline.MapperFunc(func(ctx context.Context, in *In) (*Out, error) {
            return &Out{ID: in.ID}, nil
        })),
    ),
)
if err != nil {
    return err
}
return engine.Run(ctx)
```

Imperative style is equivalent:

```go
engine, _ := pipeline.New(ctx, pipeline.WithLogger(log))
_ = engine.RegisterSource(ctx, "nats", src)
_ = engine.RegisterStorage(ctx, "ch", store)
_ = engine.RegisterStream(ctx, "logs",
    pipeline.WithSources("nats"),
    pipeline.WithStorage("ch", clickhouse.WithQuery(...)),
)
return engine.Run(ctx)
```

Open drivers directly with typed options (no Config bags / global connector registry):

```go
src, _ := nats.Open(ctx, natsURL, source.WithFormat("json"))       // source/nats
store, _ := clickhouse.Open(ctx, dsn, clickhouse.WithInitQuery(initSQL))
pub, _ := pubnats.Open(ctx, natsURL)                               // storage/nats
```

File scenarios map YAML/HCL fields onto these same typed APIs inside `cmd/eventstream`
via tag-selected driver catalogs in `cmd/eventstream/{source,storage}`
(`-tags all` / `nats` / `clickhouse` / …). Library drivers under `source/*` and
`storage/*` have no build tags — include them by importing the package.

## Run the service binary

File scenarios (YAML / TOML / HCL) are a **service-only** feature under `cmd/eventstream`.
They are not part of the importable framework API.

```sh
eventstream server --config=./config.hcl
eventstream validate --config=./config.yml
```

Docker:

```sh
docker run -d -it --rm -v ./custom.config.hcl:/config.hcl \
  geniusrabbit/eventstream
```

See [`deploy/develop/config.hcl`](deploy/develop/config.hcl) (also `.yml` / `.toml`).

## Scenario config reference

### File structure

The format is chosen by file extension: `.hcl`, `.yml` / `.yaml`, `.toml`, `.json`.

| Key | Required | Purpose |
|-----|----------|---------|
| `version` | no | Config version label |
| `description` | no | Free-form description |
| `stores` | yes | Named destinations (at least one) |
| `sources` | yes | Named inputs (at least one) |
| `streams` | yes | Named pipelines `source → store` (at least one) |

Common rules:

- `driver` may be omitted. It is then taken from the `connect` scheme
  (`nats://host:4222` → `nats`).
- String values support environment substitution: `{{@env:VAR_NAME}}`.
  Other `{{...}}` placeholders are left as is and are message templates.
- The binary contains only the drivers selected by build tags
  (`-tags all`, `nats`, `clickhouse`, …). Using a driver that was not built in
  fails with `undefined driver "<name>" (not built into this binary)`.

```hcl
version     = "1"
description = "Logs pipeline"

stores  { ... }
sources { ... }
streams { ... }
```

### `stores.<name>`

| Param | Drivers | Purpose |
|-------|---------|---------|
| `driver` | all | `clickhouse`, `vertica`, `nats`, `natstream`, `kafka`, `redistream`, `redis`, `ping` / `http` |
| `connect` | all | Connection DSN. Required for every driver except `ping` / `http` |
| `init_query` | `clickhouse`, `vertica` | String or list of SQL statements run once on start (create database / tables) |
| `method` | `ping`, `http` | Default HTTP method for the store (`GET`, `POST`, …) |
| `url` | `ping`, `http` | Default URL; used instead of `connect` when `connect` is empty |

The block size of SQL stores is set per stream (`buffer_size`), not on the store.

```hcl
stores {
  clickhouse_1 {
    driver     = "clickhouse"
    connect    = "{{@env:CLICKHOUSE_STORE_CONNECT}}"
    init_query = [
      "CREATE DATABASE IF NOT EXISTS stat;",
      "CREATE TABLE IF NOT EXISTS stat.testlog (timestamp DateTime, service String, msg String) Engine=Memory;",
    ]
  }

  nats_store {
    connect = "nats://localhost:4222"   # driver inferred from the scheme
  }

  redis_store {
    driver  = "redis"
    connect = "redis://localhost:6379/0"
  }

  webhook {
    driver = "ping"
    url    = "http://localhost:8080/hook"
    method = "POST"
  }
}
```

### `sources.<name>`

| Param | Required | Purpose |
|-------|----------|---------|
| `driver` | no | `nats`, `natstream`, `kafka`, `redistream` (`redis` is an alias) |
| `connect` | yes | Connection DSN (topic / subject / stream is part of the DSN) |
| `format` | no | Payload decoder: `json`, `bson`, `raw` (default `raw`) |

```hcl
sources {
  nats_1 {
    driver  = "nats"
    connect = "{{@env:NATS_STREAM_CONNECT}}"
    format  = "json"
  }
}
```

### `streams.<name>` — common params

| Param | Required | Purpose |
|-------|----------|---------|
| `source` | yes | Name of an item in `sources` |
| `store` | yes | Name of an item in `stores`; its driver decides which params below apply |
| `where` | no | Condition on the message; only matching messages are written (`service=="info"`, `!error`) |
| `metrics` | no | Prometheus metrics updated per processed message |

Each `metrics` entry:

| Field | Purpose |
|-------|---------|
| `name` | Metric name (`.` is replaced with `_`) |
| `type` | Metric type; only `counter` is supported |
| `namespace`, `subsystem` | Optional Prometheus name prefixes |
| `tags` | Labels with `{{field}}` value templates. A **list with one map**; only the first map is used |

```yaml
streams:
  info_logs:
    source: nats_1
    store: clickhouse_1
    where: 'service=="info"'
    metrics:
      - name: logs.received
        type: counter
        tags:
          - service: "{{service}}"
            type: "{{type}}"
    # ...driver-specific params
```

In HCL, `tags = { service = "{{service}}", type = "{{type}}" }` works with two or
more labels. A single label currently fails to parse (the HCL loader turns it into
a plain map); use YAML / TOML / JSON for single-label metrics.

The stream name is also its ID inside the store.

### `streams.<name>` — SQL stores (`clickhouse`, `vertica`)

| Param | Purpose |
|-------|---------|
| `sql_query` | Full `INSERT` statement with message placeholders `{{field}}`, `{{field:type}}`, `{{field:date\|layout}}`. When set, `target` is ignored |
| `target` | Destination table, used with `fields` when `sql_query` is empty |
| `fields` | Column ↔ message mapping (see below). Optional with `sql_query` (used as the value template) |
| `iterate_by` | Name of an array field; one row is written per element, available as `$iter.<name>` |
| `buffer_size` | Rows per write block (default `1000`) |
| `write_timeout` | Flush interval in milliseconds (default `1000`) |

Either `sql_query` or `target` / `fields` is required.

`fields` accepts a comma-separated string, a list, or a map (`column = "value"`).
Each item has one of these forms:

| Form | Meaning |
|------|---------|
| `msg` | Column `msg` ← message field `msg` |
| `service=srv` | Column `service` ← message field `srv` |
| `error=err:string` | Same, with a cast type |
| `ts=@toTimestamp({{timestamp:date}})` | Column ← raw SQL expression with message placeholders |
| `ext=$iter.iterator:string` | Column ← current element of `iterate_by` |

Types: `string`, `fix`, `uuid`, `int`, `int8`, `int32`, `int64`, `uint`, `uint8`,
`uint32`, `uint64`, `float`, `bool`, `ip`, `date`, `unixnano`, `[]int32`, `[]int64`, `[]string`.
`date` takes an optional Go layout after `|`: `{{timestamp:date|2006-01-02 15:04:05}}`.

```hcl
streams {
  # Full query
  log_query {
    store  = "clickhouse_1"
    source = "nats_1"
    sql_query = <<Q
      INSERT INTO stat.testlog (service, msg, error, timestamp)
        VALUES({{srv}}, {{msg}}, {{err}}, toTimestamp({{timestamp:date}}))
    Q
  }

  # Table + fields, with batching
  log_fields {
    store         = "clickhouse_1"
    source        = "nats_1"
    target        = "stat.testlog"
    buffer_size   = 5000
    write_timeout = 2000
    fields = [
      "service=srv",
      "msg",
      "error=err:string",
      "timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
    ]
  }

  # One row per element of the "iterator" array
  log_iter {
    store      = "clickhouse_1"
    source     = "nats_1"
    target     = "stat.testlog"
    iterate_by = "iterator"
    fields     = "service=srv,msg,ext=$iter.iterator:string"
  }
}
```

The same in YAML:

```yaml
streams:
  log_fields:
    store: clickhouse_1
    source: nats_1
    target: stat.testlog
    buffer_size: 5000
    write_timeout: 2000
    fields:
      service: srv
      msg: msg
      error: "err:string"
```

### `streams.<name>` — publish stores (`nats`, `natstream`, `kafka`, `redistream`)

| Param | Purpose |
|-------|---------|
| `targets` | Required. List of output messages. Each target has `fields` (message shape with `{{field}}` templates) and an optional `where`; a target is published only when its `where` matches |

```hcl
streams {
  split_errors {
    store   = "nats_store"
    source  = "nats_1"
    targets = [
      {
        fields {
          type    = "message"
          service = "{{service}}"
          message = "{{msg}}"
        }
        where = "!error"
      },
      {
        fields {
          type    = "error"
          service = "{{service}}"
          error   = "{{msg}}"
        }
        where = "error"
      }
    ]
  }
}
```

```yaml
streams:
  split_errors:
    store: nats_store
    source: nats_1
    targets:
      - fields: { type: message, service: "{{service}}", message: "{{msg}}" }
        where: "!error"
      - fields: { type: error, service: "{{service}}", error: "{{msg}}" }
        where: "error"
```

### `streams.<name>` — key/value store (`redis`)

| Param | Purpose |
|-------|---------|
| `key` | Required. Key template, e.g. `stat:{{service}}:{{type}}` |
| `incrementor` | `true` → `INCR` the key per message; `false` (default) → `SET` the message JSON |
| `expiration` | Key TTL as a Go duration (`30m`, `24h`); empty means no TTL |

```hcl
streams {
  counters {
    store       = "redis_store"
    source      = "nats_1"
    key         = "stat:{{service}}:{{type}}"
    incrementor = true
    expiration  = "24h"
  }
}
```

### `streams.<name>` — HTTP store (`ping`, `http`)

| Param | Purpose |
|-------|---------|
| `url` | Overrides the store URL for this stream (`{{field}}` templates allowed) |
| `method` | Overrides the store HTTP method for this stream. `GET` sends no body; any other value sends a `POST` with the message JSON |
| `content_type` | Content type of the `POST` body (default `application/json`) |

```hcl
streams {
  notify {
    store        = "webhook"
    source       = "nats_1"
    url          = "http://localhost:8080/hook?service={{service}}"
    method       = "POST"
    content_type = "application/json"
    where        = "error"
  }
}
```

Full working examples: [`deploy/develop/config.hcl`](deploy/develop/config.hcl),
[`config.yml`](deploy/develop/config.yml), [`config.toml`](deploy/develop/config.toml).

## Sources

- **kafka**
- **NATS** & **NATS stream**
- **Redis** stream

## Storages

- **Clickhouse**
- **Vertica**
- **kafka** / **NATS** / **Redis** stream
- **ping** (HTTP)

## Metrics & health (service)

```env
SERVER_PROFILE_MODE=net
SERVER_PROFILE_LISTEN=:6060
```

- Prometheus: `/metrics`
- Health: `/healthcheck`

## Package layout

| Package | Role |
|---------|------|
| `eventstream` | Core interfaces (`Sourcer`, `Storager`; `Streamer` alias) |
| `eventstream/pkg/message` | Message types |
| `eventstream/source` | Shared source subscriber + options |
| `eventstream/source/{nats,kafka,…}` | Source drivers (`source.go`) |
| `eventstream/storage` | Shared publish storage + stream options |
| `eventstream/storage/{nats,kafka,…}` | Publish storage drivers (`storage.go`) |
| `eventstream/storage/{clickhouse,sql,…}` | DB / specialty storage drivers |
| `eventstream/stream` | `Streamer` interface + wrapper |
| `eventstream/pkg/pipeline` | `Engine` (private registries), Options, Mapper |
| `eventstream/pkg/condition` | Stream `where` conditions |
| `eventstream/pkg/converter` | Message format converters |
| `cmd/eventstream` | Service binary + tag-selected driver catalogs + scenario wiring |

## TODO

- [ ] Add processing custom error metrics
- [ ] Add MySQL / PostgreSQL / MongoDB storage
- [X] HTTP/Ping, Redis, Kafka, NATS source/storage
- [X] Framework API (`pipeline.Engine`)
- [X] Health check + Prometheus metrics
- [X] Stream `where` conditions
- [X] Service configs: HCL / YAML / TOML
