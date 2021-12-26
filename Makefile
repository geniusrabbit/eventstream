SHELL := /bin/bash -o pipefail
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

BUILD_GOOS ?= $(or ${DOCKER_DEFAULT_GOOS},linux)
BUILD_GOARCH ?= $(or ${DOCKER_DEFAULT_GOARCH},amd64)
BUILD_GOARM ?= 7
BUILD_CGO_ENABLED ?= 0

LOCAL_TARGETPLATFORM=${BUILD_GOOS}/${BUILD_GOARCH}
ifeq (${BUILD_GOARCH},arm)
	LOCAL_TARGETPLATFORM=${BUILD_GOOS}/${BUILD_GOARCH}/v${BUILD_GOARM}
endif

COMMIT_NUMBER ?= $(shell git log -1 --pretty=format:%h)
TAG_VALUE ?= $(shell git describe --exact-match --tags ${COMMIT_NUMBER})

ifeq (${TAG_VALUE},)
	TAG_VALUE = commit-${COMMIT_NUMBER}
endif

PROJDIR ?= $(CURDIR)/../
PROJECT_NAME ?= eventstream

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

OS_LIST = linux darwin
ARCH_LIST = amd64 arm64 arm

APP_TAGS = all

CONTAINER_IMAGE ?= geniusrabbit/eventstream

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
fmt: ## Run formatting code
	@echo "Fix formatting"
	@gofmt -w ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi


.PHONY: tidy
tidy: ## sanitize/update modules
	go mod tidy


.PHONY: build
build: test ## Build application
	@echo "Build application"
	@rm -rf .build/eventstream
	@rm -rf .build
	@for os in $(OS_LIST); do \
		for arch in $(ARCH_LIST); do \
			if [ "$$os/$$arch" != "darwin/arm" ]; then \
				echo "Build $$os/$$arch"; \
				GOOS=$$os GOARCH=$$arch CGO_ENABLED=${BUILD_CGO_ENABLED} GOARM=${BUILD_GOARM} \
					go build -ldflags "-s -w -X internal.appVersion=`date -u +%Y%m%d.%H%M%S` -X internal.commit=${COMMIT_NUMBER}"  \
						-tags ${APP_TAGS} -o .build/$$os/$$arch/eventstream cmd/eventstream/main.go; \
				if [ "$$arch" = "arm" ]; then \
					mkdir -p .build/$$os/$$arch/v${BUILD_GOARM}; \
					mv .build/$$os/$$arch/eventstream .build/$$os/$$arch/v${BUILD_GOARM}/eventstream; \
				fi \
			fi \
		done \
	done


.PHONY: run
run: build
	docker-compose -p ${PROJECT_NAME} -f deploy/develop/docker-compose.yml build service
	docker-compose -p ${PROJECT_NAME} -f deploy/develop/docker-compose.yml run --service-ports service


.PHONY: stop
stop:
	docker-compose -p ${PROJECT_NAME} -f deploy/develop/docker-compose.yml stop


.PHONY: destroy
destroy: stop
	docker-compose -p ${PROJECT_NAME} -f deploy/develop/docker-compose.yml down


.PHONY: build-docker
build-docker: build ## Build production docker image and push to the hub.docker.com
	@echo "Build docker image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker buildx build \
		--push --platform linux/amd64,linux/arm64,linux/arm,darwin/amd64,darwin/arm64 \
		-t  ${CONTAINER_IMAGE}:${TAG_VALUE} -t ${CONTAINER_IMAGE}:latest -f deploy/production/Dockerfile .


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.DEFAULT_GOAL := help