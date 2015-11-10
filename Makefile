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
