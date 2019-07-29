module github.com/geniusrabbit/eventstream

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/demdxx/gocast v0.0.0-20160708134729-106586117e3c
	github.com/geniusrabbit/notificationcenter v0.0.0-20181001112336-107e5ac01489
	github.com/hashicorp/hcl v1.0.0
	github.com/kshvakov/clickhouse v1.3.6
	github.com/lib/pq v1.0.0
	github.com/myesui/uuid v1.0.0
	github.com/nats-io/go-nats v1.7.2 // indirect
	github.com/nats-io/nats v1.7.2
	github.com/nats-io/nats.go v1.8.1
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	gopkg.in/alexcesaro/statsd.v2 v2.0.0
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/geniusrabbit/notificationcenter => ../notificationcenter
