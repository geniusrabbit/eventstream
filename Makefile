
PROJDIR ?= $(CURDIR)/../../../../../www

buildapp:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o .build/eventstream cmd/eventstream/main.go

builddocker:
	docker build -t geniusrabbit/eventstream .

build: buildapp builddocker

run:
	docker run --rm -it -e DEBUG=true \
		-v .build/config.yml:/config.yml \
		geniusrabbit/eventstream

destroy:
	-docker rmi -f geniusrabbit/eventstream

drun:
	# go run cmd/eventstream/main.go --config=config.example.hcl --debug
	go run cmd/eventstream/main.go --config=eventstream.hcl --debug

dcbuild:
	docker build -t eventstream -f Develop.docker .

dcrun:
	docker run --rm -it -e DEBUG=true --name eventstream \
		--link nats:nats --link grclickhouse:clickhouse \
		-v $(PROJDIR)/:/project eventstream
	docker network connect --link telegraf:metrics influxdb eventstream
