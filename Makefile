SHELL := /bin/bash -o pipefail
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

BUILD_GOOS ?= linux
BUILD_GOARCH ?= amd64
BUILD_CGO_ENABLED ?= 0

COMMIT_NUMBER ?= staging # $(shell git log -1 --pretty=format:%h)

PROJDIR ?= $(CURDIR)/../
MAIN ?= eventstream

TMP_BASE := .tmp
TMP := $(TMP_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
TMP_BIN = $(TMP)/bin
TMP_ETC := $(TMP)/etc
TMP_LIB := $(TMP)/lib
TMP_VERSIONS := $(TMP)/versions

unexport GOPATH
export GOPATH=$(abspath $(TMP))
export GO111MODULE := on
export GOBIN := $(abspath $(TMP_BIN))
export PATH := $(GOBIN):$(PATH)
# Go 1.13 defaults to TLS 1.3 and requires an opt-out.  Opting out for now until certs can be regenerated before 1.14
# https://golang.org/doc/go1.12#tls_1_3
export GODEBUG := tls13=0

GOLANGLINTCI_VERSION := latest
GOLANGLINTCI := $(TMP_VERSIONS)/golangci-lint/$(GOLANGLINTCI_VERSION)
$(GOLANGLINTCI):
	$(eval GOLANGLINTCI_TMP := $(shell mktemp -d))
	cd $(GOLANGLINTCI_TMP); go get github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGLINTCI_VERSION)
	@rm -rf $(GOLANGLINTCI_TMP)
	@rm -rf $(dir $(GOLANGLINTCI))
	@mkdir -p $(dir $(GOLANGLINTCI))
	@touch $(GOLANGLINTCI)


GOMOCK_VERSION := v1.3.1
GOMOCK := $(TMP_VERSIONS)/mockgen/$(GOMOCK_VERSION)
$(GOMOCK):
	$(eval GOMOCK_TMP := $(shell mktemp -d))
	cd $(GOMOCK_TMP); go get github.com/golang/mock/mockgen@$(GOMOCK_VERSION)
	@rm -rf $(GOMOCK_TMP)
	@rm -rf $(dir $(GOMOCK))
	@mkdir -p $(dir $(GOMOCK))
	@touch $(GOMOCK)

.PHONY: deps
deps: $(GOLANGLINTCI) $(GOMOCK)

.PHONY: generate-code
generate-code: ## Generate mocks for the project
	@echo "Generate mocks for the project"
	@go generate ./...

.PHONY: golint
golint: $(GOLANGLINTCI)
	golangci-lint run -v ./...

.PHONY: lint
lint: golint

.PHONY: test
test: ## Run package test
	go test -race ./...

.PHONY: fmt
fmt: ## format code
	gofmt -w `find -name "*.go" -type f`

.PHONY: tidy
tidy: ## sanitize/update modules
	go mod tidy

.PHONY: build
build: test
	@mkdir -p .build
	@rm -rf .build/eventstream
	GOOS=${BUILD_GOOS} GOARCH=${BUILD_GOARCH} CGO_ENABLED=${BUILD_CGO_ENABLED} \
		go build --tags all \
			-ldflags "-s -w -X internal.appVersion=`date -u +%Y%m%d.%H%M%S` -X internal.commit=${COMMIT_NUMBER}" \
			-o ".build/eventstream" cmd/eventstream/main.go

.PHONY: run
run: build
	docker-compose -p ${MAIN} -f deploy/develop/docker-compose.yml build service
	docker-compose -p ${MAIN} -f deploy/develop/docker-compose.yml run --service-ports service

.PHONY: stop
stop:
	docker-compose -p ${MAIN} -f deploy/develop/docker-compose.yml stop

.PHONY: destroy
destroy: stop
	docker-compose -p ${MAIN} -f deploy/develop/docker-compose.yml down

.PHONY: image
image: build ## Build docker image
	@echo "Build docker image"
	docker build -t eventstream:${COMMIT_NUMBER} -f deploy/production/Dockerfile .

.PHONY: image_push
image_push: image ## Build docker image and push to the hub.docker.com
	@echo "Build docker image and push to the hub.docker.com"
	docker tag eventstream:${COMMIT_NUMBER} geniusrabbit/eventstream
	docker push geniusrabbit/eventstream

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help