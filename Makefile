GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
BIN_NAME=terraform-provider-appgatesdp
TEST?=./appgate
GORELEASER_VERSION = 0.143.0
ACCTEST_PARALLELISM?=20
TEST_COUNT?=1
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
HOSTNAME=appgate.com
NAMESPACE=appgate
NAME=appgatesdp
VERSION=0.6.11


build:
	go build -o $(BIN_NAME)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go test $(TEST)

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v -count $(TEST_COUNT) -parallel $(ACCTEST_PARALLELISM) $(TESTARGS) -timeout 120m

dev: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}
	mv ${BIN_NAME} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}


bin/goreleaser: bin/goreleaser-${GORELEASER_VERSION}
	@ln -sf goreleaser-${GORELEASER_VERSION} bin/goreleaser
bin/goreleaser-${GORELEASER_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | bash -s -- -b ./bin/ v${GORELEASER_VERSION}
	@mv bin/goreleaser $@

.PHONY: release
release: bin/goreleaser # Publish a release
	bin/goreleaser release
