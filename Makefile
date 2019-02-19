PREFIX ?= "/opt/project"

all: build

$(GOPATH)/bin/dep:
	mkdir -p $(GOPATH)/bin
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | GOPATH=$(GOPATH) GOOS=linux sh

build_dep: $(GOPATH)/bin/dep
	GOPATH=$(GOPATH) GOOS=linux $(GOPATH)/bin/dep ensure

update_dep: $(GOPATH)/bin/dep
	GOPATH=$(GOPATH) GOOS=linux $(GOPATH)/bin/dep ensure --update

build: build_dep
	GOPATH=$(GOPATH) GOOS=linux go install Donders-Institute/hpc-torque-helper/...

doc:
	@GOPATH=$(GOPATH) GOOS=linux godoc -http=:6060

test: build_dep
	@GOPATH=$(GOPATH) GOOS=linux GOCACHE=off go test -v Donders-Institute/hpc-torque-helper/test/...

install: build
	@install -D $(GOPATH)/bin/* $(PREFIX)/bin

clean:
	@rm -rf bin
