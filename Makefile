#!/usr/bin/env make
.DEFAULT_GOAL := help
.SHELLFLAGS   := -eou pipefail


.PHONY: modules
modules: ## Tidy and vendor Go modules for local development.
	@go mod tidy
	@go mod vendor

.PHONY: test
test: ## Run tests.
	@go clean -testcache && go test -race ./... -coverprofile=coverage.out

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests.
	@go clean -testcache && go test -bench Benchmark* ./... -benchmem

.PHONY: help
help:
	@grep -E '^[%a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
