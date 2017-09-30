# LDFLAGS += -X "/setting.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
# LDFLAGS += -X "/setting.BuildGitHash=$(shell git rev-parse HEAD)"

OS := $(shell uname)

DATA_FILES := $(shell find conf | sed 's/ /\\ /g')

BUILD_FLAGS:=-o git_txt -v
TAGS=sqlite
NOW=$(shell date -u '+%Y%m%d%I%M%S')
GOVET=go vet
GOLINT=golint -set_exit_status

GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*" ! -path "*/bindata.go")
PACKAGES ?= $(filter-out dev.sigpipe.me/dashie/git.txt/integrations,$(shell go list ./... | grep -v /vendor/))

.PHONY: build clean

all: build

check: test

web: build
	./git_txt web

vet:
	$(GOVET) git.txt.go

lint:
	$(GOLINT) $(PACKAGES)

build:
	go build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -tags '$(TAGS)'

build-dev: govet
	go build $(BUILD_FLAGS) -tags '$(TAGS)'

build-dev-race: govet
	go build $(BUILD_FLAGS) -race -tags '$(TAGS)'

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -delete

test:
	go test -cover -v $(PACKAGES)

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error -i unknwon $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w -i unknwon $(GOFILES)