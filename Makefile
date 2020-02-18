GOOS ?= "linux"

SECRET ?= "my-secret"

VERSION ?= "master"

GOLDFLAGS = "-X github.com/Donders-Institute/hpc-torque-helper/internal/grpc.secret=$(SECRET)"

ifndef GOPATH
	GOPATH := $(HOME)/go
endif

.PHONY: build

all: build

build:
	GOOS=$(GOOS) go build --ldflags $(GOLDFLAGS) -o $(GOPATH)/bin/trqhelpd github.com/Donders-Institute/hpc-torque-helper/cmd/trqhelpd

doc:
	@GOOS=$(GOOS) godoc -http=:6060

test: build
	GOOS=$(GOOS) go test --ldflags $(GOLDFLAGS) -v github.com/Donders-Institute/hpc-torque-helper/...

release:
	VERSION=$(VERSION) rpmbuild --undefine=_disable_source_fetch -bb build/rpm/centos7.spec

github_release:
	@scripts/gh-release.sh $(VERSION) false

clean:
	@rm -rf $(GOPATH)/bin/trqhelpd
	@rm -rf $(GOPATH)/pkg/*/Donders-Institute/hpc-torque-helper
