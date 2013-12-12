TARGETS := github.com/joshk/hustle

VERSION_VAR := github.com/joshk/hustle.VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)

REV_VAR := github.com/joshk/hustle.RevisionString
REPO_REV := $(shell git rev-parse --sq HEAD)

GO_TAG_ARGS ?= -tags full
TAGS_VAR := github.com/joshk/hustle.BuildTags
GOBUILD_LDFLAGS := -ldflags "-X $(VERSION_VAR) $(REPO_VERSION) -X $(REV_VAR) $(REPO_REV) -X $(TAGS_VAR) '$(GO_TAG_ARGS)' "

ADDR := :8661

all: clean test

test: build fmtpolice
	go test -race $(GOBUILD_LDFLAGS) $(GO_TAG_ARGS) -x -v $(TARGETS)

build: deps
	go install $(GOBUILD_LDFLAGS) $(GO_TAG_ARGS) -x $(TARGETS)
	go build -o $${GOPATH%%:*}/bin/hustle-server $(GOBUILD_LDFLAGS) $(GO_TAG_ARGS) ./hustle-server

deps:
	if [ ! -e $${GOPATH%%:*}/src/github.com/joshk/hustle ] ; then \
		mkdir -p $${GOPATH%%:*}/src/github.com/joshk ; \
		ln -sv $(PWD) $${GOPATH%%:*}/src/github.com/joshk/hustle ; \
	fi
	go get -x $(TARGETS)

clean:
	go clean -x $(TARGETS) || true
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
		find $${GOPATH%%:*}/pkg -name '*hustle*' -exec rm -v {} \; ; \
	fi

fmtpolice:
	set -e; for f in $(shell git ls-files '*.go'); do gofmt $$f | diff -u $$f - ; done

serve:
	$${GOPATH%%:*}/bin/hustle-server -a $(ADDR)

.PHONY: all build clean deps serve test fmtpolice
