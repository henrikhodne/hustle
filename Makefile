HUSTLE_PACKAGE := github.com/joshk/hustle
TARGETS := $(HUSTLE_PACKAGE) $(HUSTLE_PACKAGE)/hustle-server

VERSION_VAR := $(HUSTLE_PACKAGE).VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)

REV_VAR := $(HUSTLE_PACKAGE).RevisionString
REPO_REV := $(shell git rev-parse --sq HEAD)

GO ?= go
GODEP ?= godep
GO_TAG_ARGS ?= -tags full
TAGS_VAR := $(HUSTLE_PACKAGE).BuildTags
GOBUILD_LDFLAGS := -ldflags "-X $(VERSION_VAR) $(REPO_VERSION) -X $(REV_VAR) $(REPO_REV) -X $(TAGS_VAR) '$(GO_TAG_ARGS)' "

HUSTLE_HTTPADDR ?= :8661
HUSTLE_WSADDR ?= :8663
HUSTLE_STATSADDR ?= :8665

all: clean test save

test: build fmtpolice
	$(GO) test -race $(GOBUILD_LDFLAGS) $(GO_TAG_ARGS) -x -v $(TARGETS)

build: deps
	$(GO) install $(GOBUILD_LDFLAGS) $(GO_TAG_ARGS) -x $(TARGETS)

deps: public/pusher.js public/pusher.min.js
	if [ ! -e $${GOPATH%%:*}/src/$(HUSTLE_PACKAGE) ] ; then \
		mkdir -p $${GOPATH%%:*}/src/github.com/joshk ; \
		ln -sv $(PWD) $${GOPATH%%:*}/src/$(HUSTLE_PACKAGE) ; \
	fi
	$(GODEP) restore

clean:
	$(GO) clean -x $(TARGETS) || true
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
		find $${GOPATH%%:*}/pkg -name '*hustle*' -exec rm -v {} \; ; \
	fi

save:
	$(GODEP) save -copy=false $(HUSTLE_PACKAGE)

fmtpolice:
	set -e; for f in $(shell git ls-files '*.go'); do gofmt $$f | diff -u $$f - ; done

serve:
	exec $${GOPATH%%:*}/bin/hustle-server \
	  -http-addr=$(HUSTLE_HTTPADDR) \
	  -ws-addr=$(HUSTLE_WSADDR) \
	  -stats-addr=$(HUSTLE_STATSADDR)

public/pusher.js:
	curl -s -L -o $@ http://js.pusher.com/2.1/pusher.js

public/pusher.min.js:
	curl -s -L -o $@ http://js.pusher.com/2.1/pusher.min.js

.PHONY: all build clean deps serve test fmtpolice
