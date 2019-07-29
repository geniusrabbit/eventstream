
PROJDIR ?= $(CURDIR)/../

buildapp:
	docker run -it --rm --env CGO_ENABLED=0 --env GO111MODULE=on \
    -v="$(CURDIR):/eventstream" -w="/eventstream" golang:latest \
    go build -tags all -a -installsuffix cgo -gcflags '-B' \
			-ldflags '-s -w' -o ".build/eventstream" "cmd/eventstream/main.go"

builddocker:
	docker build -t geniusrabbit/eventstream -f deploy/docker/Dockerfile .

build: buildapp builddocker

destroy:
	-docker rmi -f geniusrabbit/eventstream

run:
	go run -tags all cmd/eventstream/main.go --config=config.example.hcl --profiler=:6060 --debug

dcbuild:
	docker build -t eventstream -f Develop.dockerfile .

dcrun: dcbuild
	docker run --rm -it -e DEBUG=true -p 6060:6060 --name eventstream \
		--link nats:nats-streaming --link clickhouse \
		-v $(PROJDIR)/:/project eventstream
