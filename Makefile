.PHONY: init test build

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
PROVIDER_DIR := $(HOME)/.terraform.d/plugins
PROVIDER := terraform-provider-uaa_v0.0.1
PROVIDER_PATH := $(PROVIDER_DIR)/registry.terraform.io/orange-cloudfoundry/uaa/0.0.1/$(GOOS)_$(GOARCH)
ARTIFACT := terraform-provider-uaa

init:
	go mod tidy
	go mod download

test:
	# go clean -testcache
	go test -v -timeout 10m ./test --tags=containerized

build:
	go build -o $(ARTIFACT)
	mkdir -p $(PROVIDER_PATH)
	cp $(ARTIFACT) $(PROVIDER_PATH)/$(PROVIDER)
