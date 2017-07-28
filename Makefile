# LDFLAGS += -X "/setting.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
# LDFLAGS += -X "/setting.BuildGitHash=$(shell git rev-parse HEAD)"

OS := $(shell uname)

DATA_FILES := $(shell find conf | sed 's/ /\\ /g')

BUILD_FLAGS:=-o git_txt -v
TAGS=sqlite
NOW=$(shell date -u '+%Y%m%d%I%M%S')
GOVET=go tool vet -composites=false -methods=false -structtags=false

PACKAGES ?= $(filter-out dev.sigpipe.me/dashie/git.txt/integrations,$(shell go list ./... | grep -v /vendor/))

GENERATED := bindata/bindata.go

.PHONY: build clean

all: build

check: test

web: build
	./git_txt web

govet:
	$(GOVET) git.txt.go

build: $(GENERATED)
	go build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -tags '$(TAGS)'

build-dev: $(GENERATED) govet
	go build $(BUILD_FLAGS) -tags '$(TAGS)'

build-dev-race: $(GENERATED) govet
	go build $(BUILD_FLAGS) -race -tags '$(TAGS)'

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -delete

test:
	go test -cover -v $(PACKAGES)

bindata: bindata/bindata.go

bindata/bindata.go: $(DATA_FILES)
	@hash go-bindata > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/jteeuwen/go-bindata/...; \
		fi
	go-bindata -o=$@ -ignore="\\.DS_Store|README.md|TRANSLATORS" -pkg=bindata conf/...
