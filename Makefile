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

GO ?= go
GOVET=$(GO) vet
GOLINT=golint -set_exit_status
GOFMT ?= gofmt -s

GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*" ! -path "*/bindata.go")
PACKAGES ?= $(filter-out dev.sigpipe.me/dashie/git.txt/integrations,$(shell go list ./... | grep -v /vendor/))
PACKAGES_ALL ?= $(shell go list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -name "*.go" -type f)
XGO_DEPS = "--deps=https://s3.sigpipe.me/tarballs/1-mingw-libgnurx-2.5.1-src.tar.gz https://s3.sigpipe.me/tarballs/2-file-5.32.tar.gz"

ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
else
	ifneq ($(DRONE_BRANCH),)
		VERSION ?= $(subst release/v,,$(DRONE_BRANCH))
	else
		VERSION ?= master
	endif
endif

### Targets build and checks

.PHONY: build clean

all: build

web: build
	./$(EXECUTABLE) web

vet:
	$(GOVET) $(PACKAGES_ALL)

lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES_ALL); do golint -set_exit_status $$PKG || exit 1; done;

build:
	$(GO) build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -tags '$(TAGS)'

build-dev: vet
	$(GO) build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -tags '$(TAGS)'

build-dev-race: vet
	$(GO) build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -race -tags '$(TAGS)'

clean: clean-mac
	$(GO) clean -i ./...
	rm -f $(EXECUTABLE)

clean-mac:
	find . -name ".DS_Store" -delete

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

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@diff=$$($(GOFMT) -d $(GOFILES)); \
		if [ -n "$$diff" ]; then \
			echo "Please run 'make fmt' and commit the result:"; \
			echo "$${diff}"; \
			exit 1; \
		fi;

### Targets for tests
# Use PACKAGES instead of PACKAGES_ALL because the integrations tests are run separately

test: fmt-check
	$(GO) test -cover -v $(PACKAGES)

### Targets for releases
.PHONY: release
release: release-dirs release-linux release-copy release-check

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

# This is an ugly hack, we will not need that when xgo will support cmake and sources-order
release-lx64: release-dirs release-build-lx64 release-copy release-check release-pack-lx64

release-build-lx64:
	cp $(EXECUTABLE) $(DIST)/binaries/$(EXECUTABLE)-linux-x86_64
	cp -r conf $(DIST)/release/
	cp README.md LICENSE* $(DIST)/release/

release-pack-lx64:
	cd $(DIST)/release; tar czvf git.txt_$(VERSION).tgz $(EXECUTABLE)-linux-x86_64 conf README.md LICENSE*
