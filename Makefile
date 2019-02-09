
PROJDIR ?= $(CURDIR)/../../../../

buildapp:
	docker run -it --rm --env CGO_ENABLED=0 --env GOPATH="/project" \
    -v="`pwd`/../../../..:/project" -w="/project/src/github.com/geniusrabbit/eventstream" golang:latest \
    go build -a -installsuffix cgo -gcflags '-B' -ldflags '-s -w' -o ".build/eventstream" "cmd/eventstream/main.go"

builddocker:
	docker build -t geniusrabbit/eventstream -f deploy/Dockerfile .

build: buildapp builddocker

destroy:
	-docker rmi -f geniusrabbit/eventstream

run:
	go run -tags all cmd/eventstream/main.go --config=config.example.hcl --debug

dcbuild:
	docker build -t eventstream -f Develop.dockerfile .

dcrun: dcbuild
	docker run --rm -it -e DEBUG=true --name eventstream \
		--link nats:nats-streaming --link clickhouse \
		-v $(PROJDIR)/:/project eventstream
