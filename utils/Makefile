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
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

SRC ?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

#
# Rules
#
.PHONY: build clean run hooks fmt test lint up

all: fmt lint
	@true

lint: $(GOMETALINTER)
	$(GOMETALINTER) ./... --vendor --fast --disable=maligned

$(GOMETALINTER):
	$(GOGET) -u github.com/alecthomas/gometalinter
	$(GOMETALINTER) --install 1>/dev/null

fmt:
	gofmt -l -w $(SRC)
