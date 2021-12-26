 
stores {
  # CREATE TABLE stat.testlog (
  #    timestamp        DateTime
  #  , datemark         Date default toDate(timestamp)
  #  , service          String
  #  , msg              String
  #  , error            String
  #  , created_at       DateTime default now()
  # ) Engine=MergeTree(datemark, (service), 8192);

  clickhouse_1 {
    connect = "@env:CLICKHOUSE_STORE_CONNECT"
    buffer  = 1000
  }

  nats_store {
    connect = "@env:NATS_STORE_CONNECT"
  }
}

// Source could be any supported stream service like kafka, nats, etc...
sources {
  nats_1 {
    connect = "@env:NATS_STREAM_CONNECT"
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
      INSERT INTO testlog (service, msg, error, timestamp)
        VALUES({{srv}}, {{msg}}, {{err}}, toTimestamp({{timestamp:date}}))
    Q

    where   = "service==\"info\""
  }

  log_3 {
    store  = "clickhouse_1"
    source = "nats_1"

    target = "testlog"
    # Optional if fields in log and in message the same
    fields = "service=srv,msg,error=err,timestamp=@toTimestamp({{timestamp:date}})"
  }

  log_4 {
    store  = "clickhouse_1"
    source = "nats_1"

    target = "testlog"
    fields = [
      "service=srv",
      "msg",
      "error=err:string",
      "timestamp=@toTimestamp('{{timestamp:date|2006-01-02 15:04:05}}')",
    ]
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
