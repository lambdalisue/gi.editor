VERSION := $$(git describe --tags --always --dirty)

# http://postd.cc/auto-documented-makefile/
.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

deps: ## Install dev dependencies
	GO111MODULE=off \
	    go get -u -v github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=off \
	    go get -u -v golang.org/x/tools/cmd/goimports

lint: ## Run lint
	goimports -l -w .
	golangci-lint run ./...

test: ## Run test
	go test ./...

build: ## Build binary
	go build -ldflags "-X main.appVersion=${VERSION}" ./...
