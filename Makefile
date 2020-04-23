VERSION=$$(cat VERSION)
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
BIN_NAME=terraform-provider-appgate_${VERSION}
TEST?=./...


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
