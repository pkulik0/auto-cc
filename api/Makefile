PKG = github.com/pkulik0/autocc/api
VERSION = $(shell git describe --tags --always --dirty)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
FLAGS = -ldflags "-X $(PKG)/internal/version.BuildTime=$(BUILD_TIME)"

version:
	@echo $(VERSION) > ./internal/version/version.txt

run: version
	@go run $(FLAGS) ./cmd

build:
	@mkdir -p ./build
	@go build $(FLAGS) -o ./build/api ./cmd

test:
	@go test -v ./...

deps:
	@go install go.uber.org/mock/mockgen@latest

generate: deps
	@go generate ./...
