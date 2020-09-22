VERSION=$$(cat VERSION)
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
BIN_NAME=terraform-provider-appgate_${VERSION}
TEST?=./...
GORELEASER_VERSION = 0.143.0

build:
	go build -o $(BIN_NAME)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go test $(TEST)

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v -count=1 $(TESTARGS) -timeout 120m

example: build
	@mv $(BIN_NAME) examples/
	@cp examples/$(BIN_NAME) examples/aws/appgate-resources



bin/goreleaser: bin/goreleaser-${GORELEASER_VERSION}
	@ln -sf goreleaser-${GORELEASER_VERSION} bin/goreleaser
bin/goreleaser-${GORELEASER_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | bash -s -- -b ./bin/ v${GORELEASER_VERSION}
	@mv bin/goreleaser $@

.PHONY: release
release: bin/goreleaser # Publish a release
	bin/goreleaser release
