SHELL := /bin/bash

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)
TARGET_MAX_CHAR_NUM=20

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`
BINARY := anchorctl

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type d -name '*.go' -not -path "./vendor/*")

.PHONY: fmt lint build run docker

fmt:
	@gofmt -s -w .

lint:
	@golint -set_exit_status ./pkg/cmd ./pkg/logging ./pkg/kubetest ./cmd

test:
	@go test -short ./pkg/cmd ./pkg/logging ./pkg/kubetest ./cmd

test-coverage:
	@go test -short -coverprofile cover.out -covermode=atomic ./pkg/cmd ./pkg/logging ./pkg/kubetest ./cmd
	@cat cover.out >> coverage.txt

build:
	@go build $(LDFLAGS) -o ./anchorctl -v ./cmd/main.go

run: fmt lint build
	./anchorctl test -f ./samples/kube-test.yaml -k kubetest -v 5

docker:
	@docker build -t "covarity/$(BINARY):$(VERSION)" \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile .

help:
	@echo ''
	@echo 'Usage:'
	@echo '  $(YELLOW)make$(RESET) $(GREEN)<target>$(RESET)'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\.\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)