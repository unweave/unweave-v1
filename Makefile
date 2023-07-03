GOTEST=go test
GOFILES=$(shell go list ./... | grep -v 'wip')
GOBUILD=go build

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build ci echo

all: help

build: ## compile the repo
	$(GOBUILD) -o bin/unweave .

bin:
	mkdir bin

test: ## test the repo
	$(GOTEST) -race -timeout 1m $(GOFILES)

echo:
	echo $(GOFILES)

format:
	go fmt ./...

lint: ## run the linters
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3 run --timeout=10m

ci: lint test build ## run the ci jobs

help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
