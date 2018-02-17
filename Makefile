
PROJDIR ?= $(CURDIR)/../../../../

buildapp:
	docker run -it --rm --env CGO_ENABLED=0 --env GOPATH="/project" \
    -v="`pwd`/../../../..:/project" -w="/project/src/github.com/geniusrabbit/eventstream" golang:1.9.4 \
    go build -a -installsuffix cgo -gcflags '-B' -ldflags '-s -w' -o ".build/eventstream" "cmd/eventstream/main.go"

builddocker:
	docker build -t geniusrabbit/eventstream -f deploy/Dockerfile .

build: buildapp builddocker

run:
	docker run --rm -it -e DEBUG=true \
		-v .build/config.yml:/config.yml \
		geniusrabbit/eventstream

destroy:
	-docker rmi -f geniusrabbit/eventstream

drun:
	go run cmd/eventstream/main.go --config=config.example.hcl --debug

dcbuild:
	docker build -t eventstream -f Develop.dockerfile .

dcrun: dcbuild
	docker run --rm -it -e DEBUG=true --name eventstream \
		--link nats:nats-streaming --link clickhouse \
		-v $(PROJDIR)/:/project eventstream
