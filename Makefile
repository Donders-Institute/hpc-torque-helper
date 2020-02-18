GOOS ?= "linux"

SECRET ?= "my-secret"

VERSION ?= "master"

GOLDFLAGS = "-X github.com/Donders-Institute/hpc-torque-helper/internal/grpc.secret=$(SECRET)"

.PHONY: build

all: build

build:
	GOPATH=$(GOPATH) GOOS=$(GOOS) go build $(GOLDFLAG) -o bin/trqhelpd github.com/Donders-Institute/hpc-torque-helper/cmd/trqhelpd

doc:
	@GOPATH=$(GOPATH) GOOS=$(GOOS) godoc -http=:6060

test: build
	@GOPATH=$(GOPATH) GOOS=$(GOOS) GOCACHE=off go test $(GOLDFLAG) -v github.com/Donders-Institute/hpc-torque-helper/...
release:
	VERSION=$(VERSION) rpmbuild --undefine=_disable_source_fetch -bb build/rpm/centos7.spec

github_release:
	@scripts/gh-release.sh $(VERSION) false

clean:
	@rm -rf $(GOPATH)/bin/cluster-*
	@rm -rf $(GOPATH)/bin/trqhelpd
	@rm -rf $(GOPATH)/pkg/*/Donders-Institute/hpc-torque-helper
