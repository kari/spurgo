NAME ?= spurgo
VERSION ?= $(shell git describe --tags --always --dirty)
# Supported platforms for cross-compilation
TARGETS ?= linux/amd64 darwin/arm64 windows/amd64

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"
CMD_PATH := .

.PHONY: all clean test lint build dist dev

all: test lint build

lint:
	go vet ./...

test:
	go test -v ./...

build: clean
	mkdir -p build/
	go build -o build/$(NAME) $(LDFLAGS) $(CMD_PATH)

dist: clean lint test
	mkdir -p dist/
	for target in $(TARGETS); do \
		GOOS=$$(echo $$target | cut -d"/" -f1); \
		GOARCH=$$(echo $$target | cut -d"/" -f2); \
		EXT=""; \
    if [ "$$GOOS" = "windows" ]; then EXT=".exe"; fi; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -o dist/$(NAME)-$$GOOS-$$GOARCH$$EXT \
    $(LDFLAGS) \
    $(CMD_PATH); \
	done

clean:
	rm -rf build/ dist/
	go clean

dev: build
	./build/$(NAME)
