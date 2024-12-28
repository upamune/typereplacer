BINARY_NAME := typereplacer
PKG := github.com/upamune/typereplacer

.PHONY: all build test fmt lint clean

all: build

## Build the CLI tool
build:
	@echo "===> Building $(BINARY_NAME) ..."
	go build -o bin/$(BINARY_NAME) ./cmd/typereplacer

## Run all tests (including golden tests)
test:
	@echo "===> Running tests ..."
	go test -v ./...

## Format Go code with gofmt (in-place)
fmt:
	@echo "===> Formatting code ..."
	gofmt -w .
	go mod tidy

## Lint code (using basic 'go vet')
lint:
	@echo "===> Linting code ..."
	go vet ./...

## Clean up built files
clean:
	@echo "===> Cleaning up ..."
	rm -rf bin/$(BINARY_NAME)
