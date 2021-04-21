SHELL := /bin/bash

COMMIT_SHA=$(shell git rev-parse --short HEAD)

# Go and compilation related variables
BUILD_DIR ?= out
SOURCE_DIRS = cmd pkg test
RELEASE_DIR ?= release

# Docs build related variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
HOST_BUILD_DIR=$(BUILD_DIR)/$(GOOS)-$(GOARCH)
GOPATH ?= $(shell go env GOPATH)
ORG := github.com/adrianriobo
REPOPATH ?= $(ORG)/qe-eventmanager
PACKAGE_DIR := packaging/$(GOOS)

SOURCES := $(shell git ls-files  *.go ":^vendor")

# https://golang.org/cmd/link/
LDFLAGS := $(VERSION_VARIABLES) -extldflags='-static' ${GO_EXTRA_LDFLAGS}

# Add default target
.PHONY: default
default: install

# Create and update the vendor directory
.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: check
check: cross test cross-lint

# Start of the actual build targets

.PHONY: install
install: $(SOURCES)
	go install -ldflags="$(LDFLAGS)" $(GO_EXTRA_BUILDFLAGS) ./cmd

$(BUILD_DIR)/macos-amd64/qe-eventmanager: $(SOURCES)
	GOARCH=amd64 GOOS=darwin go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/macos-amd64/qe-eventmanager $(GO_EXTRA_BUILDFLAGS) ./cmd

$(BUILD_DIR)/linux-amd64/qe-eventmanager: $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/linux-amd64/qe-eventmanager $(GO_EXTRA_BUILDFLAGS) ./cmd

$(BUILD_DIR)/windows-amd64/qe-eventmanager.exe: $(SOURCES)
	GOARCH=amd64 GOOS=windows go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/windows-amd64/qe-eventmanager.exe $(GO_EXTRA_BUILDFLAGS) ./cmd

.PHONY: cross ## Cross compiles all binaries
cross: $(BUILD_DIR)/macos-amd64/qe-eventmanager $(BUILD_DIR)/linux-amd64/qe-eventmanager $(BUILD_DIR)/windows-amd64/qe-eventmanager.exe

.PHONY: test
test:
	go test -race --tags build -v -ldflags="$(VERSION_VARIABLES)" ./pkg/... ./cmd/...

.PHONY: clean ## Remove all build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(GOPATH)/bin/qe-eventmanager

.PHONY: fmt
fmt:
	@gofmt -l -w $(SOURCE_DIRS)

$(GOPATH)/bin/golangci-lint:
	pushd /tmp && GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.37.1 && popd

# Run golangci-lint against code
.PHONY: lint cross-lint
lint: $(GOPATH)/bin/golangci-lint
	$(GOPATH)/bin/golangci-lint run

cross-lint: $(GOPATH)/bin/golangci-lint
	GOOS=darwin $(GOPATH)/bin/golangci-lint run
	GOOS=linux $(GOPATH)/bin/golangci-lint run
	GOOS=windows $(GOPATH)/bin/golangci-lint run