include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

.PHONY: $(PKG) test
SHELL := /bin/bash
PKGS := github.com/Clever/mgotail
$(eval $(call golang-version-check,1.6))

export MONGO_URL ?= mongodb://localhost:27017/test

test: $(PKGS)
$(PKGS): golang-test-all-deps
	go get -d -t $@
	$(call golang-test-all,$@)
