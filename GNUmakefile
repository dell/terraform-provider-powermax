default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

go-pkgs:
	go mod download

lint:
	golangci-lint run --fix

build:
	go build 

