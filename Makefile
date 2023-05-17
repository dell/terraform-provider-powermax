lint:
	golangci-lint run --fix

vet:
	echo "running go vet"
	go vet

fmt:
	terraform fmt -recursive examples/
	go fmt

code-check: lint vet fmt

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

generate:
	go generate ./...

download:
	go mod download

build: download
	mkdir -p out
	go build -v -o ./out

sweep:
	go test -v ./powermax -timeout 5h -sweep=all

all: download code-check testacc
	
