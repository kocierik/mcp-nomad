.PHONY: install build test clean release

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=nomad-mcp-server

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf dist/

install:
	$(GOINSTALL) ./...

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME) -transport=stdio

release:
	./scripts/build.sh

# Development tools
tools:
	$(GOGET) github.com/golangci/lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securego/gosec/v2/cmd/gosec@latest

lint:
	golangci-lint run

security:
	gosec ./... 