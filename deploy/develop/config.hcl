# Version of the config
version = "1"
description = "Example of the logs writing"

stores {
  clickhouse_1 {
    connect = "{{@env:CLICKHOUSE_STORE_CONNECT}}"
    buffer  = 1000
    init_query = [
      "CREATE DATABASE IF NOT EXISTS stat;",
      <<Q
        CREATE TABLE IF NOT EXISTS stat.testlog (
          timestamp        DateTime
        , datemark         Date default toDate(timestamp)
        , service          String
        , msg              String
        , error            String
        , ext              String
        , created_at       DateTime default now()
        ) Engine=Memory COMMENT 'The test table';
      Q
    ]
  }

  nats_store {
    connect = "{{@env:NATS_STORE_CONNECT}}"
  }
}

// Source could be any supported stream service like kafka, nats, etc...
sources {
  nats_1 {
    connect = "{{@env:NATS_STREAM_CONNECT}}"
    format  = "json"
  }
}

variable "service_info" {
  default = "info"
}

// Streams it's pipelines which have source and destination store
streams {
  log_2 {
    store   = "clickhouse_1"
    source  = "nats_1"

    sql_query = <<Q
      INSERT INTO stat.testlog (service, msg, error, timestamp)
        VALUES({{srv}}, {{msg}}, {{err}}, toTimestamp({{timestamp:date}}))
    Q

    where   = "service==\"info\""
  }

  log_3 {
    store  = "clickhouse_1"
    source = "nats_1"

    target = "stat.testlog"
    # Optional if fields in log and in message the same
    fields = "service=srv,msg,error=err,timestamp=@toTimestamp({{timestamp:date}})"
  }

  log_4 {
    store  = "clickhouse_1"
    source = "nats_1"

    target = "stat.testlog"
    fields = [
      "service=srv",
      "msg",
      "error=err:string",
      "timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
    ]
  }

  log_5 {
    store  = "clickhouse_1"
    source = "nats_1"

    target = "stat.testlog"
    iterate_by = "iterator"
    fields = [
      "service=srv",
      "msg",
      "error=err:string",
      "ext=$iter.iterator:string",
      "timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
    ]
    where   = <<ST
      type=="iterator"
    ST
  }

  nats_1 {
    # Write two messages into the nats_1 store
    # message1 = {"type": "message", "service": "...", message: "..."}
    # message2 = {"type": "error",   "service": "...", message: "..."}
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
