BINARY    := promptc
CMD       := ./cmd/$(BINARY)
BUILD_DIR := bin

MODEL_URL := https://dl.fbaipublicfiles.com/fasttext/supervised-models/lid.176.ftz

GO       := go
GOFLAGS  ?=
LDFLAGS  ?=

# Suppress macOS linker warning about duplicate -lc++ from go-fasttext CGO
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
export CGO_LDFLAGS += -Wl,-no_warn_duplicate_libraries
endif

VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE     ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  += -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.DEFAULT_GOAL := build

## build: Compile the binary for the current platform
.PHONY: build
build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) $(CMD)

## build-all: Cross-compile for all supported platforms
.PHONY: build-all
build-all:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		output=$(BUILD_DIR)/$(BINARY)-$$os-$$arch$$ext; \
		echo "Building $$output"; \
		GOOS=$$os GOARCH=$$arch $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $$output $(CMD) || exit 1; \
	done

## run: Run promptc with ARGS (e.g. make run ARGS="explain closures")
.PHONY: run
run: build
	./$(BUILD_DIR)/$(BINARY) $(ARGS)

## install: Install promptc to $GOPATH/bin
.PHONY: install
install:
	$(GO) install $(GOFLAGS) -ldflags "$(LDFLAGS)" $(CMD)

## test: Run all tests
.PHONY: test
test:
	$(GO) test ./... -race -count=1

## test-verbose: Run all tests with verbose output
.PHONY: test-verbose
test-verbose:
	$(GO) test ./... -race -count=1 -v

## cover: Run tests with coverage report (excludes cmd/ integration tests)
.PHONY: cover
cover:
	$(GO) test $$($(GO) list ./... | grep -v promptc/cmd/) -race -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML report: make cover-html"

## cover-html: Open coverage report in browser
.PHONY: cover-html
cover-html: cover
	$(GO) tool cover -html=coverage.out -o coverage.html
	open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || echo "Open coverage.html in your browser"

## vet: Run go vet
.PHONY: vet
vet:
	$(GO) vet ./...

## lint: Run golangci-lint
.PHONY: lint
lint:
	golangci-lint run ./...

## fmt: Format all Go source files
.PHONY: fmt
fmt:
	gofmt -s -w .

## fmt-check: Check formatting without modifying files
.PHONY: fmt-check
fmt-check:
	@test -z "$$(gofmt -s -l .)" || (echo "Files need formatting:"; gofmt -s -l .; exit 1)

## check: Run fmt-check, vet, lint, and tests
.PHONY: check
check: fmt-check vet lint test

## tidy: Run go mod tidy
.PHONY: tidy
tidy:
	$(GO) mod tidy

## download-model: Download fastText language detection model
.PHONY: download-model
download-model:
	@mkdir -p data
	curl -L -o data/lid.176.ftz $(MODEL_URL)
	@echo "Model downloaded to data/lid.176.ftz"

## clean: Remove build artifacts and coverage files
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f $(BINARY)

## help: Show this help message
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /' | column -t -s ':'
