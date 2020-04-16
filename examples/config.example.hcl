 
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
    connect = "clickhouse://clickhouse:9000/stat"
    driver  = "clickhouse"
    buffer  = 1000
  }

  // hdfs_1 {
  //   connect = "hdfs://hdfs:8020/"
  //   driver  = "hdfs"
  //   buffer  = 1000
  //   tmpdir  = "/tmp/hdfs/"
  // }

  nats_1 {
    connect = "nats://nats:4222/?topics=metrics"
    driver  = "nats"
  }
}

// Source could be any supported stream service like kafka, nats, etc...
sources {
  nats_1 {
    connect = "nats://nats:4222/?topics=topic1,topic2"
    format  = "json"
    driver  = "nats"
  }
  kafka_1 {
    connect = "nats://nats:4222/group?topics=topic1"
    driver  = "nats"
  }
}

// Streams it's pipelines which have source and destination store
streams {
  log_2 {
    store   = "clickhouse_1"
    source  = "nats_1"

    rawitem = <<Q
      INSERT INTO testlog (service, msg, error, timestamp)
        VALUES({{srv}}, {{msg}}, {{err}}, toTimestamp({{timestamp:date}}))
    Q

    where   = "service = ${"\"info\""}"
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
    store   = "nats_1"
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
