# Example eventstream scenario: ClickHouse, NATS, Kafka, Redis (KV), Redis Streams, metrics.
#
# Validate without starting the pipeline:
#   eventstream validate --config=examples/config.example.hcl
#
# Run (requires binary built with all drivers, e.g. -tags all):
#   eventstream server --config=examples/config.example.hcl
#
# Prometheus metrics: enable SERVER_PROFILE_MODE=net and open /metrics
#
# Note: "redis" as a source is an alias for Redis Streams (redistream).
#       "redis" as a store is key/value SET or INCR — a different driver.

version     = "1"
description = "Example data-flow scenarios: ClickHouse, NATS, Kafka, Redis, Redis Streams, metrics"

stores {
  clickhouse_analytics {
    driver  = "clickhouse"
    connect = "{{@env:CLICKHOUSE_CONNECT}}"
    init_query = [
      "CREATE DATABASE IF NOT EXISTS events;",
      <<Q
        CREATE TABLE IF NOT EXISTS events.raw_events (
          timestamp  DateTime
        , service    String
        , level      String
        , message    String
        , error      String
        , created_at DateTime default now()
        ) Engine=Memory COMMENT 'Example analytics table';
      Q
    ]
  }

  nats_out {
    driver  = "nats"
    connect = "{{@env:NATS_STORE_CONNECT}}"
  }

  kafka_out {
    driver  = "kafka"
    connect = "{{@env:KAFKA_STORE_CONNECT}}"
  }

  redistream_out {
    driver  = "redistream"
    connect = "{{@env:REDIS_STREAM_STORE_CONNECT}}"
  }

  # Redis key/value store (SET or INCR) — not the same as redistream above.
  redis_cache {
    driver  = "redis"
    connect = "{{@env:REDIS_KV_CONNECT}}"
  }
}

sources {
  nats_in {
    driver  = "nats"
    connect = "{{@env:NATS_SOURCE_CONNECT}}"
    format  = "json"
  }

  kafka_in {
    driver  = "kafka"
    connect = "{{@env:KAFKA_SOURCE_CONNECT}}"
    format  = "json"
  }

  # Redis Streams consumer (driver "redis" is also accepted as an alias).
  redistream_in {
    driver  = "redistream"
    connect = "{{@env:REDIS_STREAM_SOURCE_CONNECT}}"
    format  = "json"
  }
}

streams {
  # NATS → ClickHouse: target + fields mapping, where filter, Prometheus metrics.
  events_to_clickhouse {
    store  = "clickhouse_analytics"
    source = "nats_in"

    target = "events.raw_events"
    fields = "service,level,message=msg,error=err,timestamp=@toTimestamp({{timestamp:date}})"
    where  = "level!=\"debug\""

    metrics = [
      {
        name      = "events_total"
        type      = "counter"
        namespace = "eventstream"
        subsystem = "example"
        tags = [
          { service = "{{service}}", level = "{{level}}" }
        ]
      }
    ]
  }

  # Kafka → ClickHouse: raw INSERT SQL with message field placeholders.
  kafka_to_clickhouse {
    store  = "clickhouse_analytics"
    source = "kafka_in"

    sql_query = <<Q
      INSERT INTO events.raw_events (service, level, message, error, timestamp)
        VALUES({{service}}, {{level}}, {{msg}}, {{err}}, toTimestamp({{timestamp:date}}))
    Q

    where = "service!=\"\""
  }

  # Redis Streams → ClickHouse: target + fields (array syntax).
  redisstream_to_clickhouse {
    store  = "clickhouse_analytics"
    source = "redistream_in"

    target = "events.raw_events"
    fields = [
      "service",
      "level",
      "message=msg",
      "error=err:string",
      "timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
    ]
  }

  # NATS → NATS: fan-out with per-target fields and where conditions.
  nats_fanout {
    store  = "nats_out"
    source = "nats_in"

    targets = [
      {
        fields {
          type    = "message"
          service = "{{service}}"
          level   = "{{level}}"
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

  # Kafka → Kafka: simple republish of the incoming message.
  kafka_republish {
    store  = "kafka_out"
    source = "kafka_in"

    targets = [
      {
        fields {
          service = "{{service}}"
          level   = "{{level}}"
          message = "{{msg}}"
          error   = "{{err}}"
        }
      }
    ]
  }

  # NATS → Redis Streams: publish transformed events to an output stream.
  nats_to_redisstream {
    store  = "redistream_out"
    source = "nats_in"

    targets = [
      {
        fields {
          event   = "forwarded"
          service = "{{service}}"
          level   = "{{level}}"
          payload = "{{msg}}"
        }
        where = "level==\"info\" || level==\"warn\""
      }
    ]
  }

  # NATS → Redis KV: cache message payload by id with TTL.
  # For counters use a separate stream with incrementor = true instead of expiration.
  cache_and_count {
    store      = "redis_cache"
    source     = "nats_in"
    key        = "events:{{service}}:{{id}}"
    expiration = "24h"
    where      = "id!=\"\""
  }

  event_counter {
    store       = "redis_cache"
    source      = "nats_in"
    key         = "events:{{service}}:count"
    incrementor = true
  }
}
