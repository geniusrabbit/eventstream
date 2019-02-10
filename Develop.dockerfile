FROM golang:latest

RUN mkdir -p /project

ENV GOPATH=/project/ \
    GOBIN=/project/bin \
    PATH="$PATH:$GOBIN"

WORKDIR /project/src/github.com/geniusrabbit/eventstream

CMD make run
