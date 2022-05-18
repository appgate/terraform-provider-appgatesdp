GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
BIN_NAME=terraform-provider-appgatesdp
TEST?=./appgate
ACCTEST_PARALLELISM?=20
TEST_COUNT?=1
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
HOSTNAME=appgate.com
NAMESPACE=appgate
NAME=appgatesdp
VERSION=0.8.9


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


.PHONY: release
release: # Publish a release
	goreleaser release
