# Build-time variables (injected via ldflags).
VERSION   ?= 0.2.0-dev
COMMIT    ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILDTIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BASE_URL  ?= __CONFIG__

LDFLAGS = -s -w \
	-X github.com/carbonstop/carbonstop-cli/internal/version.Version=$(VERSION) \
	-X github.com/carbonstop/carbonstop-cli/internal/version.Commit=$(COMMIT) \
	-X github.com/carbonstop/carbonstop-cli/internal/version.BuildTime=$(BUILDTIME) \
	-X github.com/carbonstop/carbonstop-cli/internal/version.BaseURL=$(BASE_URL)

.PHONY: all build dist clean test install

all: build

build:
	go build -ldflags "$(LDFLAGS)" -o carbonstop ./cmd/carbonstop/

build-release: build

# Cross-compile for distribution.
dist:
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/carbonstop-linux-amd64       ./cmd/carbonstop/
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/carbonstop-darwin-amd64       ./cmd/carbonstop/
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/carbonstop-darwin-arm64       ./cmd/carbonstop/
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/carbonstop-windows-amd64.exe  ./cmd/carbonstop/

clean:
	rm -f carbonstop carbonstop.exe
	rm -rf dist/

test:
	go test ./...

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/carbonstop/
