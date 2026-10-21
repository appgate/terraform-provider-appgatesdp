BIN_NAME        := terraform-provider-appgatesdp
VERSION         ?= 9.9.9
LDFLAGS         := -X github.com/appgate/terraform-provider-appgatesdp/version.ProviderVersion=$(VERSION)
GOLANGCI_LINT   ?= golangci-lint
# API_VERSION limits generation and pinning to one version, e.g. API_VERSION=v25.
API_VERSION     ?=
SDPGEN_FLAGS     = $(if $(API_VERSION),-version=$(API_VERSION))

GOOS            := $(shell go env GOOS)
GOARCH          := $(shell go env GOARCH)
PLUGIN_DIR      := $(HOME)/.terraform.d/plugins/appgate.com/appgate/appgatesdp/$(VERSION)/$(GOOS)_$(GOARCH)

.PHONY: all generate pin fmt fmt-check vet lint build test test-race testacc tidy-check generate-check docs docs-check check dev release

all: check

## Code generation ----------------------------------------------------------

# Regenerate all API versions, or only API_VERSION.
generate:
	go run ./cmd/sdpgen generate $(SDPGEN_FLAGS)

# Pin versions to the current head of their specification branch.
pin:
	go run ./cmd/sdpgen pin $(SDPGEN_FLAGS)

# Fail if regeneration changes committed files.
generate-check: generate
	@git diff --exit-code || (echo "generated files are out of date: run 'make generate' and commit the result" && exit 1)
	@test -z "$$(git status --porcelain --untracked-files=all -- codegen internal templates)" || (git status --short -- codegen internal templates && echo "generation produced untracked files" && exit 1)

## Formatting and static analysis -------------------------------------------

fmt:
	gofmt -w .

# Verify formatting without modifying files.
fmt-check:
	@files="$$(gofmt -l .)"; if [ -n "$$files" ]; then echo "not gofmt formatted:"; echo "$$files"; exit 1; fi

vet:
	go vet ./...

lint:
	$(GOLANGCI_LINT) run ./...

# Verify go.mod/go.sum are tidy without modifying them.
tidy-check:
	go mod tidy -diff

## Build and test -------------------------------------------------------------

build:
	go build ./...
	go build -o $(BIN_NAME) -ldflags "$(LDFLAGS)" .

# Unit tests. Provider lifecycle tests drive a Terraform CLI against an
# in-memory SDP controller; no live environment is needed.
test:
	go test ./...

test-race:
	go test -race ./...

# Acceptance tests against live controllers, see internal/acctest.
# Read-only checks run with TF_ACC=1; set SDP_ACC_ALLOW_DESTRUCTIVE=1 to
# create and delete objects.
testacc:
	TF_ACC=1 go test ./internal/acctest/... -run '^TestAcc' -v -count=1 -timeout 120m $(TESTARGS)

## Documentation --------------------------------------------------------------

docs:
	go tool tfplugindocs generate --provider-name appgatesdp

docs-check: docs
	@git diff --exit-code -- docs || (echo "docs are out of date: run 'make docs'" && exit 1)

## Aggregate ------------------------------------------------------------------

# Every quality gate; used by CI.
check: fmt-check vet lint build test-race tidy-check generate-check

dev: build
	mkdir -p $(PLUGIN_DIR)
	mv $(BIN_NAME) $(PLUGIN_DIR)

release:
	goreleaser release
