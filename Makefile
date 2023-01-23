SHELL := /bin/bash -o pipefail
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

BUILD_GOOS ?= $(or ${DOCKER_DEFAULT_GOOS},linux)
BUILD_GOARCH ?= $(or ${DOCKER_DEFAULT_GOARCH},amd64)
BUILD_GOARM ?= 7
BUILD_CGO_ENABLED ?= 0
DOCKER_BUILDKIT ?= 1

LOCAL_TARGETPLATFORM=${BUILD_GOOS}/${BUILD_GOARCH}
ifeq (${BUILD_GOARCH},arm)
	LOCAL_TARGETPLATFORM=${BUILD_GOOS}/${BUILD_GOARCH}/v${BUILD_GOARM}
endif

COMMIT_NUMBER ?= $(or ${DEPLOY_COMMIT_NUMBER},)
ifeq (${COMMIT_NUMBER},)
	COMMIT_NUMBER = $(shell git log -1 --pretty=format:%h)
endif

TAG_VALUE ?= $(or ${DEPLOY_TAG_VALUE},)
ifeq (${TAG_VALUE},)
	TAG_VALUE = $(shell git describe --exact-match --tags `git log -n1 --pretty='%h'`)
endif
ifeq (${TAG_VALUE},)
	TAG_VALUE = commit-${COMMIT_NUMBER}
endif


PROJECT_NAME ?= eventstream

export GO111MODULE := on
# Go 1.13 defaults to TLS 1.3 and requires an opt-out.  Opting out for now until certs can be regenerated before 1.14
# https://golang.org/doc/go1.12#tls_1_3
export GODEBUG := tls13=0

OS_LIST   ?= $(or ${DEPLOY_OS_LIST},linux darwin)
ARCH_LIST ?= $(or ${DEPLOY_ARCH_LIST},amd64 arm64 arm)
APP_TAGS  ?= $(or ${APP_BUILD_TAGS},all)

CONTAINER_IMAGE ?= geniusrabbit/eventstream


.PHONY: generate-code
generate-code: ## Generate mocks for the project
	@echo "Generate mocks for the project"
	@go generate ./...


.PHONY: lint
lint:
	golangci-lint run -v ./...


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


define do_build
	@for os in $(OS_LIST); do \
		for arch in $(ARCH_LIST); do \
			if [ "$$os/$$arch" != "darwin/arm" ]; then \
				echo "Build $$os/$$arch"; \
				GOOS=$$os GOARCH=$$arch CGO_ENABLED=${BUILD_CGO_ENABLED} GOARM=${BUILD_GOARM} \
					go build \
						-ldflags "-s -w -X main.appVersion=`date -u +%Y%m%d` -X main.buildCommit=${COMMIT_NUMBER} -X main.buildVersion=${TAG_VALUE} -X main.buildDate=`date -u +%Y%m%d.%H%M%S`"  \
						-tags ${APP_TAGS} -o .build/$$os/$$arch/$(2) $(1); \
				if [ "$$arch" = "arm" ]; then \
					mkdir -p .build/$$os/$$arch/v${BUILD_GOARM}; \
					mv .build/$$os/$$arch/$(2) .build/$$os/$$arch/v${BUILD_GOARM}/$(2); \
				fi \
			fi \
		done \
	done
endef


.PHONY: build
build: ## Build application
	@echo "Build application"
	@rm -rf .build
	@$(call do_build,"cmd/eventstream/main.go",eventstream)


.PHONY: build-docker-dev
build-docker-dev: build
	echo "Build develop docker image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker build -t ${CONTAINER_IMAGE}:latest -f deploy/develop/Dockerfile .


.PHONY: build-lsg
build-lsg: ## Build logstreamgen application
	@echo "Build logstreamgen application"
	@$(call do_build,"examples/logstreamgen/main.go",logstreamgen)


.PHONY: build-docker-lsg-dev
build-docker-lsg-dev: build-lsg
	echo "Build develop docker LSG image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker build -t ${CONTAINER_IMAGE}-lsg:latest -f deploy/develop/logstreamgen.Dockerfile .


.PHONY: build-docker
build-docker: build ## Build production docker image and push to the hub.docker.com
	@echo "Build docker image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker buildx build \
		--push --platform linux/amd64,linux/arm64,linux/arm,darwin/amd64,darwin/arm64 \
		-t ${CONTAINER_IMAGE}:${TAG_VALUE} -t ${CONTAINER_IMAGE}:latest -f deploy/production/Dockerfile .


.PHONY: run
run:
	docker-compose -p ${PROJECT_NAME} -f deploy/develop/docker-compose.yml run --rm --service-ports service



.PHONY: ch
ch: ## Run clickhouse client
	docker exec -it $(PROJECT_NAME)-clickhouse-1 clickhouse-client


.PHONY: stop
stop:
	docker-compose -p ${PROJECT_NAME} -f deploy/develop/docker-compose.yml stop


.PHONY: destroy
destroy: stop
	docker-compose -p ${PROJECT_NAME} -f deploy/develop/docker-compose.yml down


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.DEFAULT_GOAL := help