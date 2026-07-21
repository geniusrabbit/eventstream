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
