.PHONY: build test lint fmt clean

VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BRANCH  := $(shell git symbolic-ref --short -q HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w \
	-X main.Version=$(VERSION) \
	-X main.GitCommit=$(COMMIT) \
	-X main.BuildTime=$(DATE) \
	-X main.GitBranch=$(BRANCH)

# The main package is ./src, so `go install` would name the binary "src":
# always build with -o raci.
build:
	go build -ldflags "$(LDFLAGS)" -o raci ./src

test:
	go test -race ./...

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

clean:
	rm -f raci coverage.out
	rm -rf dist/
