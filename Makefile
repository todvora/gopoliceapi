## simple makefile to log workflow
.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install fmt test build

build:
	@go generate
	@go build $(GOFLAGS) -o policecz

setup:
	go get github.com/Masterminds/glide
	go get -u github.com/jteeuwen/go-bindata/...

install:
	@glide install

update:
	@glide update

test:
	@go test -v $(glide novendor)

cover:
	@go test -coverprofile=coverage.txt $(glide novendor)
	@go tool cover -html=coverage.txt

travis-ci:
	go test -v -coverprofile=coverage.txt -covermode=atomic $(glide novendor)

clean:
	@go clean $(GOFLAGS) -i ./...

fmt:
	gofmt -w .

run: build
	@./policecz

## EOF