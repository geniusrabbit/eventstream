SHELL := /bin/bash -o pipefail

APP_TAGS  ?= $(or ${APP_BUILD_TAGS},all)

include deploy/build.mk

PROJECT_WORKSPACE := geniusrabbit
PROJECT_NAME ?= eventstream
DOCKER_COMPOSE := docker compose -p $(PROJECT_WORKSPACE) -f deploy/develop/docker-compose.yml
DOCKER_CONTAINER_IMAGE := ghcr.io/geniusrabbit/eventstream


.PHONY: generate-code
generate-code: ## Generate mocks for the project
	@echo "Generate mocks for the project"
	@go generate ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run -v --build-tags="$(APP_TAGS)" ./...

.PHONY: test
test: ## Run package test
	go test -race -tags "$(APP_TAGS)" ./...

.PHONY: fmt
fmt: ## Run formatting code
	@echo "Fix formatting"
	@gofmt -w ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi

.PHONY: tidy
tidy: ## sanitize/update modules
	go mod tidy

.PHONY: build
build: ## Build application
	@echo "Build application"
	@rm -rf .build
	@$(call do_build,"cmd/eventstream/main.go",eventstream)

.PHONY: build-docker-dev
build-docker-dev: build ## Build develop docker image
	@echo "Build develop docker image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker build -t ${DOCKER_CONTAINER_IMAGE}:latest -f deploy/develop/Dockerfile .

.PHONY: build-lsg
build-lsg: ## Build logstreamgen application
	@echo "Build logstreamgen application"
	@$(call do_build,"examples/logstreamgen/main.go",logstreamgen)

.PHONY: build-docker-lsg-dev
build-docker-lsg-dev: build-lsg ## Build develop docker LSG image
	@echo "Build develop docker LSG image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker build -t ${DOCKER_CONTAINER_IMAGE}-lsg:latest -f deploy/develop/logstreamgen.Dockerfile .

.PHONY: build-dev
build-dev: build-docker-lsg-dev build-docker-dev ## Build develop docker images (service + LSG)

.PHONY: build-docker
build-docker: build ## Build production docker image and push to the hub.docker.com
	@echo "Build docker image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker buildx build \
		--push --platform linux/amd64,linux/arm64,linux/arm,darwin/amd64,darwin/arm64 \
		-t ${DOCKER_CONTAINER_IMAGE}:${TAG_VALUE} -t ${DOCKER_CONTAINER_IMAGE}:latest -f deploy/production/Dockerfile .

.PHONY: run
run: ## Run development environment services
	${DOCKER_COMPOSE} run --rm --service-ports service

.PHONY: ch
ch: ## Run clickhouse client
	${DOCKER_COMPOSE} exec clickhouse clickhouse-client

.PHONY: stop
stop: ## Stop development environment services
	${DOCKER_COMPOSE} stop

.PHONY: destroy
destroy: ## Destroy development environment services
	${DOCKER_COMPOSE} down -v

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
