FROM golang:latest

RUN mkdir -p /project

ENV GOPATH=/project/ \
    GOBIN=/project/bin \
    PATH="$PATH:$GOBIN" \
    GO111MODULE=on

WORKDIR /project/eventstream

CMD ["go", "run", "-tags", "all", "cmd/eventstream/main.go", \
    "--config=config.example.hcl", "--profiler=:6060", "--debug"]
