
## Detect OS
OSS=$(shell uname)
ifeq ($(OSS),Windows_NT)
	OSS=windows
endif
ifeq ($(OSS),Linux)
	OSS=linux
endif
ifeq ($(OSS),Darwin)
	OSS=darwin
endif

## Detect processor ARCH
ARCH=$(shell uname -m)
ARMV=7
ifeq ($(ARCH),x86_64)
	ARCH=amd64
endif
ifeq ($(ARCH),i386)
	ARCH=386
endif
ifeq ($(ARCH),i686)
	ARCH=386
endif
ifeq ($(ARCH),aarch64)
	ARCH=arm64
endif
ifeq ($(ARCH),armv7l)
	ARCH=arm
	ARMV=7
endif
ifeq ($(ARCH),armv6l)
	ARCH=arm
	ARMV=6
endif
ifeq ($(ARCH),armv5l)
	ARCH=arm
endif

SHELL := /bin/bash -o pipefail

BUILD_GOOS ?= $(or ${DOCKER_DEFAULT_GOOS},${OSS})
BUILD_GOARCH ?= $(or ${DOCKER_DEFAULT_GOARCH},${ARCH})
# https://github.com/golang/go/wiki/MinimumRequirements#amd64
BUILD_GOAMD64 ?= $(or ${DOCKER_DEFAULT_GOAMD64},1)
BUILD_GOAMD64_LIST ?= $(or ${DOCKER_DEFAULT_GOAMD64_LIST},1)
BUILD_GOARM ?= $(or ${DOCKER_DEFAULT_GOARM},7)
BUILD_GOARM_LIST ?= $(or ${DOCKER_DEFAULT_BUILD_GOARM_LIST},${BUILD_GOARM})
BUILD_CGO_ENABLED ?= 0
DOCKER_BUILDKIT ?= 1

LOCAL_TARGETPLATFORM=${BUILD_GOOS}/${BUILD_GOARCH}
ifeq (${BUILD_GOARCH},arm)
	LOCAL_TARGETPLATFORM=${BUILD_GOOS}/${BUILD_GOARCH}/v${BUILD_GOARM}
endif
ifeq (${BUILD_GOARCH},arm64)
	LOCAL_TARGETPLATFORM=${BUILD_GOOS}/${BUILD_GOARCH}/v8
endif

COMMIT_NUMBER ?= $(or ${DEPLOY_COMMIT_NUMBER},)
ifeq (${COMMIT_NUMBER},)
	COMMIT_NUMBER = $(shell git log -1 --pretty=format:%h)
endif

TAG_VALUE ?= $(or ${DEPLOY_TAG_VALUE},)
ifeq (${TAG_VALUE},)
	TAG_VALUE = $(shell git describe --exact-match --tags `git log -n1 --pretty='%h'` 2>/dev/null)
endif
ifeq (${TAG_VALUE},)
	TAG_VALUE = commit-${COMMIT_NUMBER}
endif

OS_LIST   ?= $(or ${DEPLOY_OS_LIST},linux darwin)
ARCH_LIST ?= $(or ${DEPLOY_ARCH_LIST},amd64 arm64 arm)
APP_TAGS  ?= $(or ${APP_BUILD_TAGS},all)

# Prepare the list of platforms to build
define build_platform_list
	for os in $(OS_LIST); do \
		for arch in $(ARCH_LIST); do \
			if [ "$$os/$$arch" != "darwin/arm" ]; then \
				if [ "$$arch" = "arm" ]; then \
					for armv in $(BUILD_GOARM_LIST); do \
						i="$${os}/$${arch}/v$${armv}"; \
						echo -n "$${i},"; \
					done; \
				else \
					if [ "$$arch" = "amd64" ]; then \
						for amd64v in $(BUILD_GOAMD64_LIST); do \
							if [ "$$amd64v" == "1" ]; then \
								i="$${os}/$${arch}"; \
							else \
								i="$${os}/$${arch}/v$${amd64v}"; \
							fi; \
							echo -n "$${i},"; \
						done; \
					else \
						i="$${os}/$${arch}"; \
						echo -n "$${i},"; \
					fi; \
				fi; \
			fi; \
		done; \
	done;
endef

# Extract the list of platforms to build
DOCKER_PLATFORM_LIST := $(shell $(call build_platform_list))
DOCKER_PLATFORM_LIST := $(shell echo $(DOCKER_PLATFORM_LIST) | sed 's/.$$//')

# Build for all platforms from build_platform_list
define do_build
	@for os in $(OS_LIST); do \
		for arch in $(ARCH_LIST); do \
			if [ "$$os/$$arch" != "darwin/arm" ]; then \
				if [ "$$arch" = "arm64" ]; then \
					echo "Build $$os/$$arch/v8"; \
					GOOS=$$os GOARCH=$$arch CGO_ENABLED=${BUILD_CGO_ENABLED} \
						go build \
							-ldflags "-s -w -X main.appVersion=`date -u +%Y%m%d` -X main.buildCommit=${COMMIT_NUMBER} -X main.buildVersion=${TAG_VALUE} -X main.buildDate=`date -u +%Y%m%d.%H%M%S`"  \
							-tags "${APP_TAGS}" -o .build/$$os/$$arch/v8/$(2) $(1); \
						cp .build/$$os/$$arch/v8/$(2) .build/$$os/$$arch/$(2); \
				else \
					if [ "$$arch" = "arm" ]; then \
						for armv in $(BUILD_GOARM_LIST); do \
							echo "Build $$os/$$arch/v$$armv"; \
							GOOS=$$os GOARCH=$$arch CGO_ENABLED=${BUILD_CGO_ENABLED} GOARM=$$armv \
								go build \
									-ldflags "-s -w -X main.appVersion=`date -u +%Y%m%d` -X main.buildCommit=${COMMIT_NUMBER} -X main.buildVersion=${TAG_VALUE} -X main.buildDate=`date -u +%Y%m%d.%H%M%S`"  \
									-tags "${APP_TAGS}" -o .build/$$os/$$arch/v$$armv/$(2) $(1); \
						done; \
					else \
						if [ "$$arch" = "amd64" ]; then \
							for amd64v in $(BUILD_GOAMD64_LIST); do \
								if [ "$$amd64v" == "1" ]; then \
									echo "Build $$os/$$arch -> v1"; \
									GOOS=$$os GOARCH=$$arch CGO_ENABLED=${BUILD_CGO_ENABLED} GOAMD64=v$$amd64v \
										go build \
											-ldflags "-s -w -X main.appVersion=`date -u +%Y%m%d` -X main.buildCommit=${COMMIT_NUMBER} -X main.buildVersion=${TAG_VALUE} -X main.buildDate=`date -u +%Y%m%d.%H%M%S`"  \
											-tags "${APP_TAGS}" -o .build/$$os/$$arch/$(2) $(1); \
								else \
									echo "Build $$os/$$arch/v$$amd64v"; \
									GOOS=$$os GOARCH=$$arch CGO_ENABLED=${BUILD_CGO_ENABLED} GOAMD64=v$$amd64v \
										go build \
											-ldflags "-s -w -X main.appVersion=`date -u +%Y%m%d` -X main.buildCommit=${COMMIT_NUMBER} -X main.buildVersion=${TAG_VALUE} -X main.buildDate=`date -u +%Y%m%d.%H%M%S`"  \
											-tags "${APP_TAGS}" -o .build/$$os/$$arch/v$$amd64v/$(2) $(1); \
								fi; \
							done; \
						else \
							echo "Build $$os/$$arch"; \
							GOOS=$$os GOARCH=$$arch CGO_ENABLED=${BUILD_CGO_ENABLED} \
								go build \
									-ldflags "-s -w -X main.appVersion=`date -u +%Y%m%d` -X main.buildCommit=${COMMIT_NUMBER} -X main.buildVersion=${TAG_VALUE} -X main.buildDate=`date -u +%Y%m%d.%H%M%S`"  \
									-tags "${APP_TAGS}" -o .build/$$os/$$arch/$(2) $(1); \
						fi; \
					fi; \
				fi; \
			fi; \
		done; \
	done
endef
