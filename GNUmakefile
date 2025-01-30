default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	script/fmt

test:
	script/test

testacc:
	TF_ACC=1 script/test

.PHONY: fmt lint test testacc build install generate
