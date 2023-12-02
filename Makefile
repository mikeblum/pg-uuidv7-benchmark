## pg-uuidv7 build tooling

MAKEFLAGS += --silent

GOLANGCI_LINT_VERSION = v1.55.0

SHELL = bash

all: help

## help: Prints a list of available build targets.
help:
	echo "Usage: make <OPTIONS> ... <TARGETS>"
	echo ""
	echo "Available targets are:"
	echo ''
	sed -n 's/^##//p' ${PWD}/Makefile | column -t -s ':' | sed -e 's/^/ /'
	echo
	echo "Targets run by default are: `sed -n 's/^all: //p' ./Makefile | sed -e 's/ /, /g' | sed -e 's/\(.*\), /\1, and /'`"

build: install clean
	sqlc generate

## clean: Cleanup generated artifacts
clean:
	rm -rf internal

## install: Install db cli tools
install:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

## lint: Lint with golangci-lint
lint: tidy
	go mod vendor
	golangci-lint run --color always --verbose --print-resources-usage ./...

## fmt: Format with gofmt
fmt:
	go fmt ./...

# tidy: Tidy with go mod tidy
tidy:
	go clean -modcache
	go mod tidy

## pre-commit: Chain lint + test
pre-commit: lint test

## sqlc: Generate sqlc schema
sqlc: clean
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	sqlc generate

## run: go run main.go
run: sqlc
	go run main.go

## test: Test with go test
test:
	go test -test.v -race -covermode=atomic -coverprofile=coverage.out ./... && go tool cover -html=coverage.out && rm coverage.out

## test-perf: Benchmark tests with go test -bench
test-perf:
	go test -test.v -benchmem -bench=. -coverprofile=coverage-bench.out ./... && go tool cover -html=coverage-bench.out && rm coverage-bench.out

## vuln: Scan against the Go vulnerability database
vuln:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: lint fmt tidy pre-commit test test-perf vuln
