LDFLAGS += -X "dev.sigpipe.me/dashie/git.txt/setting.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
LDFLAGS += -X "dev.sigpipe.me/dashie/git.txt/setting.BuildGitHash=$(shell git rev-parse HEAD)"

OS := $(shell uname)

ifeq ($(OS), Windows_NT)
	EXECUTABLE := git_txt.exe
else
	EXECUTABLE := git_txt
endif

DATA_FILES := $(shell find conf | sed 's/ /\\ /g')
DIST := dist

BUILD_FLAGS:=-o $(EXECUTABLE) -v
TAGS=sqlite
NOW=$(shell date -u '+%Y%m%d%I%M%S')

GOVET=go vet
GOLINT=golint -set_exit_status
GO ?= go

GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*" ! -path "*/bindata.go")
PACKAGES ?= $(filter-out dev.sigpipe.me/dashie/git.txt/integrations,$(shell go list ./... | grep -v /vendor/))
XGO_DEPS = "--deps=http://download.openpkg.org/components/cache/file/file-5.32.tar.gz"
#XGO_DEPS += "--deps=https://github.com/libgit2/libgit2/archive/maint/v0.25.zip"

ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
else
	ifneq ($(DRONE_BRANCH),)
		VERSION ?= $(subst release/v,,$(DRONE_BRANCH))
	else
		VERSION ?= master
	endif
endif

### Targets

.PHONY: build clean

all: build

check: test

web: build
	./$(EXECUTABLE) web

vet:
	$(GOVET) git.txt.go

lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

build:
	$(GO) build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -tags '$(TAGS)'

build-dev: govet
	$(GO) build $(BUILD_FLAGS) -tags '$(TAGS)'

build-dev-race: govet
	$(GO) build $(BUILD_FLAGS) -race -tags '$(TAGS)'

clean:
	$(GO) clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -delete

test:
	$(GO) test -cover -v $(PACKAGES)

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error -i unknwon $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w -i unknwon $(GOFILES)

.PHONY: release
release: release-dirs release-windows release-linux release-copy release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

.PHONY: release-windows
release-windows:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/karalabe/xgo; \
	fi
	xgo $(XGO_DEPS) --image=xgo-git2go-windows -dest $(DIST)/binaries -tags 'netgo $(TAGS)' -ldflags '-linkmode external -extldflags "-static" $(LDFLAGS)' -targets 'windows/*' -out git.txt-$(VERSION) .
ifeq ($(CI),drone)
	mv /build/* $(DIST)/binaries
endif

.PHONY: release-linux
release-linux:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/karalabe/xgo; \
	fi
	xgo $(XGO_DEPS) --image=xgo-git2go-linux -dest $(DIST)/binaries -tags 'netgo $(TAGS)' -ldflags '-linkmode external -extldflags "-static" $(LDFLAGS)' -targets 'linux/*' -out git.txt-$(VERSION) .
ifeq ($(CI),drone)
	mv /build/* $(DIST)/binaries
endif

# No git2go image available for the moment
.PHONY: release-darwin
release-darwin:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/karalabe/xgo; \
	fi
	xgo $(XGO_DEPS) -dest $(DIST)/binaries -tags 'netgo $(TAGS)' -ldflags '$(LDFLAGS)' -targets 'darwin/*' -out git.txt-$(VERSION) .
ifeq ($(CI),drone)
	mv /build/* $(DIST)/binaries
endif

.PHONY: release-copy
release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

.PHONY: release-check
release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)
