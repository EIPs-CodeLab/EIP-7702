SHELL := /usr/bin/env bash

GO ?= go
GO_FILES := $(shell git ls-files '*.go')

.PHONY: fmt fmt-check mod-tidy-check vet test test-race build examples ci

fmt:
	$(GO) fmt ./...

fmt-check:
	@out="$$(gofmt -l $(GO_FILES))"; \
	if [[ -n "$$out" ]]; then \
		echo "These files are not gofmt-formatted:"; \
		echo "$$out"; \
		exit 1; \
	fi

mod-tidy-check:
	$(GO) mod tidy
	@git diff --exit-code -- go.mod go.sum

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

build:
	$(GO) build ./...

examples:
	$(GO) run ./examples/basic-authorization >/dev/null
	$(GO) run ./examples/transaction-batching >/dev/null
	$(GO) run ./examples/send-userop >/dev/null

ci: fmt-check mod-tidy-check vet test-race build examples
