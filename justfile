# Variables
NAME := env("NAME", "spurgo")
VERSION := `git describe --tags --always --dirty`
TARGETS := "linux/amd64 darwin/arm64 windows/amd64"
CMD_PATH := "."

# Display available recipes
@default:
    just --list

# Run all checks and build
all: test lint build

# Lint code
lint:
    go vet ./...

# Run tests
test:
    go test -v ./...

# Build the application
build: clean
    mkdir -p build/
    go build -ldflags "-X main.Version={{VERSION}}" -o build/{{NAME}} {{CMD_PATH}}

# Build for multiple platforms
dist: clean lint test
    mkdir -p dist/
    for target in {{TARGETS}}; do \
        GOOS="$(echo $target | cut -d'/' -f1)"; \
        GOARCH="$(echo $target | cut -d'/' -f2)"; \
        EXT=""; \
        if [ "$GOOS" = "windows" ]; then EXT=".exe"; fi; \
        GOOS=$GOOS GOARCH=$GOARCH go build \
            -ldflags "-X main.Version={{VERSION}}" \
            -o dist/{{NAME}}-$GOOS-$GOARCH$EXT \
            {{CMD_PATH}}; \
    done

# Clean build artifacts
clean:
    rm -rf build/ dist/
    go clean

# Build and run
dev: build
    ./build/{{NAME}}
