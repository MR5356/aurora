NAME ?= aurora
OUT_DIR ?= _output
BIN_DIR := $(OUT_DIR)/bin
PLUGIN_DIR := $(OUT_DIR)/plugin
MODULE_NAME = github.com/MR5356/aurora

IMAGE_REGISTRY ?= toodo/aurora
TARGET_PLATFORM ?= linux/arm64,linux/amd64

VERSION ?= $(shell git describe --tags 2>/dev/null)

# if git describe error
ifneq ($(VERSION),)
  # VERSION has already been set
else
  BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
  COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null)
  VERSION = $(BRANCH)_$(COMMIT)
endif

GO_FLAGS ?= "-s -w -X '$(MODULE_NAME)/pkg/version.Version=$(VERSION)'"

.DEFAULT_GOAL := help

.PHONY: help
help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init:  ## Initialize the project
	@go install github.com/swaggo/swag/cmd/swag@latest

version:  ## Print the version
	@echo $(VERSION)

.PHONY: doc
doc:  ## Generate documentation
	swag init -d cmd/aurora -g aurora.go
	swag fmt -d cmd/aurora -g aurora.go

.PHONY: deps
deps: doc  ## Install dependencies
	go get -d -v -t ./...

.PHONY: static
static: clean  ## Build frontend
	cd frontend && yarn && yarn build-only && cd ..
	cp -r ./frontend/dist/* pkg/server/static

.PHONY: build
build: deps  ## Build the binary
	go build -ldflags $(GO_FLAGS) -o $(BIN_DIR)/aurora ./cmd/aurora

.PHONY: release
release: clean deps static  ## Build and release the binary
	chmod +x hack/release.sh
	./hack/release.sh $(NAME) $(OUT_DIR)

.PHONY: test
test: deps  ## Run unit tests
	go test $(shell go list ./... | grep -v /docs) -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: proto
proto:  ## Generate proto
	@protoc --proto_path=. --go-grpc_out=. --go_out=paths=source_relative:. --go-grpc_opt=paths=source_relative ./pkg/domain/runner/proto/task.proto

plugin:  ## Build builtin plugins
	@go build -o $(PLUGIN_DIR)/checkout ./pkg/domain/runner/builtin/checkout

.PHONY: docker
docker:  ## Build docker image
	docker buildx build --platform $(TARGET_PLATFORM) -t $(IMAGE_REGISTRY):$(VERSION) . --push
	docker buildx build --platform $(TARGET_PLATFORM) -t $(IMAGE_REGISTRY):latest . --push

.PHONY: docker-release
docker-release: clean  ## Build and release the binary by using docker
	docker buildx build -f bin.Dockerfile --output type=local,dest=_output .

.PHONY: clean
clean:  ## Clean build artifacts
	find ./pkg/server/static/* | grep -v robots.txt | xargs rm -rf
	rm -rf $(OUT_DIR)