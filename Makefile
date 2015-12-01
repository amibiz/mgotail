SHELL := /bin/bash
PKG := github.com/Clever/mgotail
GOLINT := $(GOPATH)/bin/golint
.PHONY: $(PKG) test
GOVERSION := $(shell go version | grep 1.5)
ifeq "$(GOVERSION)" ""
  $(error must be running Go version 1.5)
endif

test: $(PKG)

$(GOLINT):
	go get github.com/golang/lint/golint

$(PKG): $(GOLINT)
	gofmt -w=true $(GOPATH)/src/$@/*.go
	$(GOLINT) $(GOPATH)/src/$@/*.go
	go get -d -t $(PKG)
	go vet $@
	go test -v $@


SHELL := /bin/bash
PKGS := $(shell go list ./... | grep -v /vendor)
GODEP := $(GOPATH)/bin/godep

$(GODEP):
	go get -u github.com/tools/godep

vendor: $(GODEP)
	$(GODEP) save $(PKGS)
	find vendor/ -path '*/vendor' -type d | xargs -IX rm -r X # remove any nested vendor directories
