OUT_DIR ?= _output
BIN_DIR := $(OUT_DIR)/bin

.DEFAULT_GOAL := help

.PHONY: help
help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps:  ## Install dependencies
	go get -d -v -t ./...

.PHONY: test
test: deps  ## Run unit tests
	go test ./... -cover

.PHONY: clean
clean:  ## Clean build artifacts
	rm -rf $(OUT_DIR)