  # Go parameters
BINDIR      := $(CURDIR)/bin
#DIST_DIRS   := find * -type d -exec
#TARGETS     := darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x windows/amd64
#TARGET_OBJS ?= darwin-amd64.tar.gz darwin-amd64.tar.gz.sha256 darwin-amd64.tar.gz.sha256sum linux-amd64.tar.gz linux-amd64.tar.gz.sha256 linux-amd64.tar.gz.sha256sum linux-386.tar.gz linux-386.tar.gz.sha256 linux-386.tar.gz.sha256sum linux-arm.tar.gz linux-arm.tar.gz.sha256 linux-arm.tar.gz.sha256sum linux-arm64.tar.gz linux-arm64.tar.gz.sha256 linux-arm64.tar.gz.sha256sum linux-ppc64le.tar.gz linux-ppc64le.tar.gz.sha256 linux-ppc64le.tar.gz.sha256sum linux-s390x.tar.gz linux-s390x.tar.gz.sha256 linux-s390x.tar.gz.sha256sum windows-amd64.zip windows-amd64.zip.sha256 windows-amd64.zip.sha256sum
BINNAME     ?= feedback-generator
BINNAME_CLIENT ?= feedback-client

GOCMD=go
FBSERVER=./cmd/server
FBCLIENT=./cmd/client
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_UNIX=$(BINNAME)_unix

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)

# Required for globs to work correctly
SHELL      = /usr/bin/env bash
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

all: test build run
build:
		$(GOBUILD) -o $(BINDIR)/$(BINNAME) -v $(FBSERVER)
build-client:
		$(GOBUILD) -o $(BINDIR)/$(BINNAME_CLIENT) -v $(FBCLIENT)	
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINDIR)/$(BINNAME)
		rm -f $(BINDIR)/$(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINDIR)/$(BINNAME) -v $(FBSERVER)
		cd $(BINDIR) & ./$(BINNAME)
run-client: 
		$(GOBUILD) -o $(BINDIR)/$(BINNAME_CLIENT) -v $(FBCLIENT)
		cd $(BINDIR) & ./$(BINNAME_CLIENT)

# Cross compilation
build-linux:
		CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(FBSERVER)
# Creating Docker build and run the code in docker container
docker-build:
		docker build -t feedback-generator .
docker-run:
		docker run -it  -d feedback-generator 