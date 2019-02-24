PREFIX ?= "/opt/project"

GOOS ?= "linux"

SECRET ?= "my-build-secret"

all: build

$(GOPATH)/bin/dep:
	mkdir -p $(GOPATH)/bin
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | GOPATH=$(GOPATH) GOOS=$(GOOS) sh

build_dep: $(GOPATH)/bin/dep
	GOPATH=$(GOPATH) GOOS=$(GOOS) $(GOPATH)/bin/dep ensure

update_dep: $(GOPATH)/bin/dep
	GOPATH=$(GOPATH) GOOS=$(GOOS) $(GOPATH)/bin/dep ensure --update

build: build_dep
	GOPATH=$(GOPATH) GOOS=$(GOOS) go install \
	-ldflags "-X github.com/Donders-Institute/hpc-torque-helper/internal/grpc.secret=$(SECRET)" \
	github.com/Donders-Institute/hpc-torque-helper/...

doc:
	@GOPATH=$(GOPATH) GOOS=$(GOOS) godoc -http=:6060

test: build_dep
	@GOPATH=$(GOPATH) GOOS=$(GOOS) GOCACHE=off go test \
	-ldflags "-X github.com/Donders-Institute/hpc-torque-helper/internal/grpc.secret=my-test-secret" \
	-v github.com/Donders-Institute/hpc-torque-helper/...

install: build
	@install -D $(GOPATH)/bin/* $(PREFIX)/bin

clean:
	@rm -rf $(GOPATH)/bin/cluster-*
	@rm -rf $(GOPATH)/bin/trqhelpd
	@rm -rf $(GOPATH)/pkg/*/Donders-Institute/hpc-torque-helper
