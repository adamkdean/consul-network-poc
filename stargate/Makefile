#    ___                      _     ___  ___  ___
#   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
#  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
# / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
# \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
#
# Consul Network proof of concept
# (c) 2018 Adam K Dean

#
# Variables
#
OWNER ?= consul-poc
PROJECT ?= stargate

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

TAG_PREFIX ?= local
TAG := $(TAG_PREFIX)/$(OWNER)/$(PROJECT)

SRC ?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PKGS = $(shell go list ./... | grep -v /vendor)
TARGET := $(shell pwd)/bin/$(PROJECT)
ENTRYPOINT = cmd/main.go

#
# Rules
#
.DEFAULT_GOAL: $(TARGET)
.PHONY: build clean run hooks fmt test lint

$(TARGET): $(SRC)
	$(GOBUILD) -o $(TARGET) $(ENTRYPOINT)

build: $(TARGET)
	@true

clean:
	@$(GOCLEAN)

run: build
	@$(TARGET)

%:
	@true

install:
	dep ensure -v

add:
	echo $(MAKECMDGOALS) | cut -d ' ' -f2- | xargs dep ensure -add

update:
	echo $(MAKECMDGOALS) | cut -d ' ' -f2- | xargs dep ensure -update

hooks: .git/hooks/pre-commit

.git/hooks/pre-commit: scripts/hooks/pre-commit.sh
	@cp -f scripts/hooks/pre-commit.sh .git/hooks/pre-commit

test: lint
	$(GOTEST) $(PKGS)

lint: $(GOMETALINTER)
	$(GOMETALINTER) ./... --vendor --fast --disable=maligned

$(GOMETALINTER):
	$(GOGET) -u github.com/alecthomas/gometalinter
	$(GOMETALINTER) --install 1>/dev/null

fmt:
	gofmt -l -w $(SRC)

docker_build:
	docker build \
		--tag $(TAG) \
		--file build/Dockerfile \
		.
