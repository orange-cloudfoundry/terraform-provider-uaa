.PHONY: test build

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
PROVIDER_DIR := $(HOME)/.terraform.d/plugins
PROVIDER := terraform-provider-uaa_v0.0.1
PROVIDER_PATH := $(PROVIDER_DIR)/registry.terraform.io/orange-cloudfoundry/uaa/0.0.1/$(GOOS)_$(GOARCH)
ARTIFACT := terraform-provider-uaa

test:
	 go test -v -timeout 10m ./test/...

build:
	 go build -o $(ARTIFACT)
	 mkdir -p $(PROVIDER_PATH)
	 cp $(ARTIFACT) $(PROVIDER_PATH)/$(PROVIDER)
