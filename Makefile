#!make

ifdef OS
GOCMD=go
else
GOCMD=/usr/local/go/bin/go
endif

# Go parameters
GOLINTCMD=golangci-lint
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=shitdetector_bot
BINARY_UNIX=$(BINARY_NAME)_unix

PROJECT_PATH=github.com/kettari/shitdetector
PROJECT_CMD ?= bot
BINARY_PATH=bin/

dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

# Assign build version
BUILD_VERSION := $(shell git describe --tags --always --dirty)

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) -v $(PROJECT_NAME)/cmd/$(PROJECT_CMD)

.PHONY: build-docker
build-docker:
	@echo ">> building docker containers"
	docker build \
	    -t $(DOCKER_REGISTRY_PREFIX)$(APP_NAME)-$(APP_ENV):$(BUILD_VERSION) \
	    --build-arg APP_CMD_DIR=bot \
	    .

.PHONY: test
test:
	docker build -f Dockerfile.test \
	    -t $(DOCKER_REGISTRY_PREFIX)test-$(APP_ENV):$(BUILD_VERSION) \
	    . && \
	docker run -v ${PWD}:/go/testdir $(DOCKER_REGISTRY_PREFIX)test-$(APP_ENV):$(BUILD_VERSION)

.PHONY: lint
lint:
	$(GOLINTCMD) run ./...
