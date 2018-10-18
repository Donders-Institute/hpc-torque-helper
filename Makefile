GOPATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PREFIX ?= "/opt/project"

all: build

$(GOPATH)/bin/dep:
	mkdir -p $(GOPATH)/bin
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | GOPATH=$(GOPATH) GOOS=linux sh

build_dep: $(GOPATH)/bin/dep
	cd src/dccn.nl; GOPATH=$(GOPATH) GOOS=linux $(GOPATH)/bin/dep ensure

update_dep: $(GOPATH)/bin/dep
	cd src/dccn.nl; GOPATH=$(GOPATH) GOOS=linux $(GOPATH)/bin/dep ensure --update

build: build_dep
	GOPATH=$(GOPATH) GOOS=linux go install dccn.nl/...

doc:
	@GOPATH=$(GOPATH) GOOS=linux godoc -http=:6060

test: build_dep
	@GOPATH=$(GOPATH) GOOS=linux GOCACHE=off go test -v dccn.nl/...

install: build
	@install -D $(GOPATH)/bin/* $(PREFIX)/bin

clean:
	@rm -rf bin
	@rm -rf pkg
