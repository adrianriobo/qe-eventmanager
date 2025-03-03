PROJECT?=github.com/devtools-qe-incubator/eventmanager
VERSION ?= 0.0.4
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
CONTAINER_MANAGER ?= podman
# Image URL to use all building/pushing image targets
IMG ?= quay.io/devtools-qe-incubator/eventmanager:${VERSION}

# Go and compilation related variables
GOPATH ?= $(shell go env GOPATH)
BUILD_DIR ?= out
SOURCE_DIRS = cmd pkg test
# https://golang.org/cmd/link/
LDFLAGS := $(VERSION_VARIABLES) -extldflags='-static' ${GO_EXTRA_LDFLAGS}
GCFLAGS := all=-N -l

# Add default target
.PHONY: default
default: install

# Create and update the vendor directory
.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: check
check: build test lint

# Start of the actual build targets

.PHONY: install
install: $(SOURCES)
	go install -ldflags="$(LDFLAGS)" $(GO_EXTRA_BUILDFLAGS) ./cmd

$(BUILD_DIR)/eventmanager: $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS) \
	-X ${PROJECT}/version.Version=${VERSION} \
	-X ${PROJECT}/version.Commit=${COMMIT} \
	-X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
	-o $(BUILD_DIR)/eventmanager $(GO_EXTRA_BUILDFLAGS) ./cmd


 
.PHONY: build 
build: $(BUILD_DIR)/eventmanager

.PHONY: test
test:
	go test -race --tags build -v -ldflags="$(VERSION_VARIABLES)" ./pkg/... ./cmd/...

.PHONY: clean ## Remove all build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(GOPATH)/bin/eventmanager

.PHONY: fmt
fmt:
	@gofmt -l -w $(SOURCE_DIRS)

$(GOPATH)/bin/golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2

# Run golangci-lint against code
.PHONY: lint
lint: $(GOPATH)/bin/golangci-lint
	$(GOPATH)/bin/golangci-lint run

# Build the container image
.PHONY: container-build
container-build: test
	${CONTAINER_MANAGER} build -t ${IMG} -f images/builder/Dockerfile .

# Push the docker image
.PHONY: container-push
container-push:
	${CONTAINER_MANAGER} push ${IMG}
