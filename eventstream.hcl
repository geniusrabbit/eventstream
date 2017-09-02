
variable "fmt" {
  datetime = "2016-01-02 15:04:05"
}

stores {
  clickhouse {
    connect = "clickhouse://clickhouse:9000/stats"
    driver  = "clickhouse"
    buffer  = 1000
  }

  metrics {
    connect = "statsd://metrics:8125/?service=myservice"
    driver  = "metrics"
    format  = "influxdb"
  }
}

// Source could be any supported stream service like kafka, nats, etc...
sources {
  nats_actions {
    connect = "nats://nats:4222/?topics=actions"
    driver  = "nats"
    format  = "json"
  }
}

// Streams it's pipelines which have source and destination store
streams {
  imps {
    store   = "clickhouse"
    source  = "nats_actions"
    target  = "actions"
    fields  = [
      "timemark=tm:unixnano",               // DateTime
      "action=act:uint",                    // UInt32
      "service=srv:fix*16",                 // FixedString(16)
      "aucid=auc:uuid",                     // FixedString(16)  -- Internal Auction ID
      "bidid=bid",                          // String
      "impid=i",                            // String
      "rtb_source=sid:uint",                // UInt64
      "rtb_access_point=acp:uint",          // UInt64
      "project=pr:uint",                    // UInt64
      "pub_company=pcb:uint",               // UInt64
      "adv_company=acv:uint",               // UInt64
      "network=net",                        // String
      "pixel=pxl",                          // String
      "platform=pl:uint",                   // UInt64
      "domain=dm",                          // String
      "app:int",                            // UInt64
      "zone=z:int",                         // UInt64
      "campaign=cmp:int",                   // UInt64
      "url=url",                            // String
      "winurl=wurl",                        // String
      "ad=ad:uint",                         // UInt64
      "ad_w=aw:int",                        // UInt32
      "ad_h=ah:int",                        // UInt32
      "jumper:int",                         // UInt64
      "ad_type=dt:int",                     // UInt8
      "pricing_model=pm:uint",              // UInt8
      "price=bpr:int",                      // UInt64           -- Number
      "cpmbid=mbp:int",                     // UInt64           -- Number
      "second_price=sbp:int",               // UInt64           -- Number
      "revenue=rev:int",                    // UInt32           -- In percents Percent * 1_000
      "potential=pt:int",                   // UInt32           -- Percent of avaited descripancy
      "ecpm=ecpm:int",                      // UInt32           -- Number
      "udid=udi",                           // FixedString(16)
      "uuid=uui:uuid",                      // FixedString(16)
      "sessid=ses:uuid",                    // FixedString(16)
      "fingerprint=fpr:uuid",               // String
      "etag=etg",                           // String
      "carrier=car:uint",                   // UInt64
      "country=cc:fix*2",                   // FixedString(2)
      "city=ct:fix*5",                      // FixedString(5)
      "latitude=lt:float",                  // Float64
      "longitude=lg:float",                 // Float64
      "language=lng:fix*5",                 // FixedString(5)
      "ip:ip",                              // String
      "ref",                                // String
      "ua",                                 // String
      "device_type=dvt:int",                // UInt32
      "device=dv:int",                      // UInt32
      "os:uint",                            // UInt32
      "browser=br:uint",                    // UInt32
      "categories=c:[]int32",               // Array(Int32)
      "adblock=ab:uint",                    // UInt8
      "private=prv:uint",                   // UInt8
      "robot=r:int",                        // UInt8
      "backup=b:int",                       // UInt8
      "x:int",                              // Int32
      "y:int",                              // Int32
      "w:int",                              // Int32
      "h:int",                              // Int32
      "subid=sd:uint",                      // UInt32 default 0
    ]
  }

  metrics {
    store   = "metrics"
    source  = "nats_actions"
    metrics = [
      {
        name = "action.counter"
        type = "increment"
        tags {
          action  = "{{act}}"
          os      = "{{os}}"
        }
        // value = "key of message field"
      }
    ]
  }
}
