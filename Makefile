GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
MODULE := $(shell awk 'NR==1{print $$2}' go.mod)
NAME=$$(grep TerraformProviderName proxmoxtf/version.go | grep -o -e 'terraform-provider-[a-z]*')
TARGETS=darwin linux windows
TERRAFORM_PLUGIN_EXTENSION=
VERSION=$$(grep TerraformProviderVersion proxmoxtf/version.go | grep -o -e '[0-9]\.[0-9]\.[0-9]')
VERSION_EXAMPLE=9999.0.0

ifeq ($(OS),Windows_NT)
	TERRAFORM_PLATFORM=windows_amd64
	TERRAFORM_PLUGIN_CACHE_DIRECTORY=$$(cygpath -u "$(shell pwd -P)")/cache/plugins
	TERRAFORM_PLUGIN_EXTENSION=.exe
else
	UNAME_S=$(shell uname -s)

	ifeq ($(UNAME_S),Darwin)
		TERRAFORM_PLATFORM=darwin_amd64
	else
		TERRAFORM_PLATFORM=linux_amd64
	endif

	TERRAFORM_PLUGIN_CACHE_DIRECTORY=$(shell pwd -P)/cache/plugins
endif

TERRAFORM_PLUGIN_DIRECTORY=$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)/registry.terraform.io/danitso/proxmox/$(VERSION)/$(TERRAFORM_PLATFORM)
TERRAFORM_PLUGIN_DIRECTORY_EXAMPLE=$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)/registry.terraform.io/danitso/proxmox/$(VERSION_EXAMPLE)/$(TERRAFORM_PLATFORM)
# TODO: don't use home dir
TERRAFORM_PLUGIN_EXECUTABLE=$(TERRAFORM_PLUGIN_DIRECTORY)/$(NAME)_v$(VERSION)_x4$(TERRAFORM_PLUGIN_EXTENSION)
TERRAFORM_PLUGIN_EXECUTABLE_EXAMPLE=$(TERRAFORM_PLUGIN_DIRECTORY_EXAMPLE)/$(NAME)_v$(VERSION_EXAMPLE)_x4$(TERRAFORM_PLUGIN_EXTENSION)

default: build

build:
	mkdir -p "$(TERRAFORM_PLUGIN_DIRECTORY)"
	rm -f "$(TERRAFORM_PLUGIN_EXECUTABLE)"
	// TODO: this is only for test code
	go build -ldflags="-s -w -X $(MODULE)/proxmox.disableHTTPSCheck=true" -o "$(TERRAFORM_PLUGIN_EXECUTABLE)"
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/danitso/proxmox/$(VERSION)/$(TERRAFORM_PLATFORM)
	cp -f "$(TERRAFORM_PLUGIN_EXECUTABLE)" ~/.terraform.d/plugins/registry.terraform.io/danitso/proxmox/$(VERSION)/$(TERRAFORM_PLATFORM)/$(NAME)_v$(VERSION)_x4$(TERRAFORM_PLUGIN_EXTENSION)

example: example-build example-init example-apply example-apply example-destroy

example-apply:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& terraform apply -auto-approve

example-build:
	mkdir -p "$(TERRAFORM_PLUGIN_DIRECTORY_EXAMPLE)"
	rm -f "$(TERRAFORM_PLUGIN_EXECUTABLE_EXAMPLE)"
	go build -o "$(TERRAFORM_PLUGIN_EXECUTABLE_EXAMPLE)"

example-destroy:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& terraform destroy -auto-approve

example-init:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& rm -f .terraform.lock.hcl \
		&& terraform init \
			-verify-plugins=false

example-plan:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& terraform plan

fmt:
	gofmt -s -w $(GOFMT_FILES)

init:
	go get ./...

targets: $(TARGETS)

test:
	go test -v ./...

e2e-box:
	@cd e2e-tests/vagrant-box && make

e2e-test:
	@cd e2e-tests && make test

$(TARGETS):
	GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build \
		-o "dist/$@/$(NAME)_v$(VERSION)-custom_x4" \
		-a -ldflags '-extldflags "-static"'
	zip \
		-j "dist/$(NAME)_v$(VERSION)-custom_$@_amd64.zip" \
		"dist/$@/$(NAME)_v$(VERSION)-custom_x4"

.PHONY: build example example-apply example-destroy example-init example-plan fmt init targets test $(TARGETS)
