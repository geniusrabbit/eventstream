 
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
    options { # Optional
      buffer = 1000
    }
  }
  hdfs_1 {
    connect = "hdfs://hdfs:8020/"
    options { # Optional
      buffer = 1000
      tmpdir = "/tmp/hdfs/"
    }
  }
}

// Source could be any supported stream service like kafka, nats, etc...
sources {
  nats_1 {
    connect = "nats://nats:4222/?topics=topic1,topic2"
    format  = "json"
  }
  kafka_1 {
    connect = "nats://nats:4222/group?topics=topic1"
  }
}

// Streams it's pipelines which have source and destination store
streams {
  // log_1 {
  //   store = "hdfs_1"
  //   source = "nats_1"
  //   target = "/path/log{{date|2006-01-02_15}}{{iterator|_%iter%}}.log.gz"
  //   fields = "service=srv:int,msg,timestamp:date|2006-01-02 15:04:05.999999-07:00"

  //   options {
  //     timefield = "timestamp"
  //     separator = ",\n"
  //     # 1Mb = 1024 * 1024 * 1
  //     maxsize   = 100
  //   }
  // }

  log_2 {
    store   = "clickhouse_1"
    source  = "nats_1"
    rawitem = <<Q
      INSERT INTO testlog (service, msg, error, timestamp)
        VALUES({{srv}}, {{msg}}, {{err}}, toTimestamp({{timestamp:date}}))
    Q
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
}
