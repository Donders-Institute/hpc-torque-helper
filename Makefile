GOOS ?= "linux"

SECRET ?= "my-secret"

VERSION ?= "master"

GOLDFLAGS = "-X github.com/Donders-Institute/hpc-torque-helper/internal/grpc.secret=$(SECRET)"

all: build

$(GOPATH)/bin/dep:
	mkdir -p $(GOPATH)/bin
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | GOPATH=$(GOPATH) GOOS=$(GOOS) sh

build_dep: $(GOPATH)/bin/dep
	GOPATH=$(GOPATH) GOOS=$(GOOS) $(GOPATH)/bin/dep ensure

update_dep: $(GOPATH)/bin/dep
	GOPATH=$(GOPATH) GOOS=$(GOOS) $(GOPATH)/bin/dep ensure --update

build: build_dep
	GOPATH=$(GOPATH) GOOS=$(GOOS) go install $(GOLDFLAG) github.com/Donders-Institute/hpc-torque-helper/...

doc:
	@GOPATH=$(GOPATH) GOOS=$(GOOS) godoc -http=:6060

test: build_dep
	@GOPATH=$(GOPATH) GOOS=$(GOOS) GOCACHE=off go test $(GOLDFLAG) -v github.com/Donders-Institute/hpc-torque-helper/...
release:
	VERSION=$(VERSION) rpmbuild --undefine=_disable_source_fetch -bb build/rpm/centos7.spec

github_release:
	scripts/gh-release.sh $(VERSION) false

clean:
	@rm -rf $(GOPATH)/bin/cluster-*
	@rm -rf $(GOPATH)/bin/trqhelpd
	@rm -rf $(GOPATH)/pkg/*/Donders-Institute/hpc-torque-helper
