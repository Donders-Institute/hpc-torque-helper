GOPATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PREFIX ?= "/opt/project"

all: build

external:
	@GOPATH=$(GOPATH) GOOS=linux go get ./...

build: external
	@GOPATH=$(GOPATH) GOOS=linux go install ./...

doc:
	@GOPATH=$(GOPATH) GOOS=linux godoc -http=:6060

test: external
	@GOPATH=$(GOPATH) GOOS=linux GOCACHE=off go test -v dccn.nl/cmd/...

install: build
	@install -D $(GOPATH)/bin/* $(PREFIX)/bin

clean:
	@rm -rf bin
	@rm -rf pkg
