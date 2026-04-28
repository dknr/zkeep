# zkeep Makefile

BINARY_ZKEEP = zkeep
BINARY_ZKEEPD = zkeepd
BUILD_DIR = build
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

.PHONY: all build clean install test

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o $(BUILD_DIR)/$(BINARY_ZKEEP) ./cmd/$(BINARY_ZKEEP)
	go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o $(BUILD_DIR)/$(BINARY_ZKEEPD) ./cmd/$(BINARY_ZKEEPD)

clean:
	rm -rf $(BUILD_DIR)

install: build
	install -m 755 $(BUILD_DIR)/$(BINARY_ZKEEPD) /usr/local/bin/$(BINARY_ZKEEPD)
	install -m 755 $(BUILD_DIR)/$(BINARY_ZKEEP) /usr/local/bin/$(BINARY_ZKEEP)
	install -m 755 etc/init.d/$(BINARY_ZKEEPD) /etc/init.d/$(BINARY_ZKEEPD)

test:
	go test ./...
