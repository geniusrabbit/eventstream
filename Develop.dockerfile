FROM golang:1.8.5

RUN mkdir -p /project

ENV GOPATH=/project/ \
    GOBIN=/project/bin \
    PATH="$PATH:$GOBIN"

WORKDIR /project/src/github.com/geniusrabbit/eventstream
ENTRYPOINT PATH="$PATH:$GOBIN" && bash
