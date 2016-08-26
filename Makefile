## simple makefile to log workflow
.PHONY: all test clean build install release

GOFLAGS ?= $(GOFLAGS:)

all: install fmt test build

build:
	@go generate
	@go build $(GOFLAGS) -o gopoliceapi

setup:
	go get github.com/Masterminds/glide
	go get -u github.com/jteeuwen/go-bindata/...

install:
	@glide install

update:
	@glide update

test: build
	@go test -v $(glide novendor)

cover: build
	@go test -coverprofile=coverage.txt $(glide novendor)
	@go tool cover -html=coverage.txt

travis-ci: build
	go test -v -coverprofile=coverage.txt -covermode=atomic $(glide novendor)

clean:
	@go clean $(GOFLAGS) -i ./...

fmt:
	gofmt -w .

run: build
	@./gopoliceapi

release:
	mkdir -p release
	GOOS=linux   GOARCH=amd64 go build -o release/gopoliceapi-linux-amd64
	GOOS=linux   GOARCH=386   go build -o release/gopoliceapi-linux-386
	GOOS=linux   GOARCH=arm   go build -o release/gopoliceapi-linux-arm
	GOOS=windows GOARCH=386   go build -o release/gopoliceapi-windows-386.exe
	GOOS=windows GOARCH=amd64 go build -o release/gopoliceapi-windows-amd64.exe


define increment_version
	# read the current version from the latest git tag
	$(eval CURRENT_VERSION = $(shell git describe --abbrev=0 | cut -c 2-))

	@echo Incrementing $(1) of $(CURRENT_VERSION)

	# increment the version, using provided component type (major|minor|patch)
	$(eval NEW_VERSION = $(shell go run scripts/semver.go $(1) $(CURRENT_VERSION)))
	@echo Incrementing to $(NEW_VERSION)
	git tag -a v$(NEW_VERSION) -m 'Release $(NEW_VERSION)'
	git push --tags
endef

version.patch:
	$(call increment_version, patch)

version.minor:
	$(call increment_version, minor)

version.major:
	$(call increment_version, major)

## EOF
