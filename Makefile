VERSION=$(shell git describe --tags)
BUILD=$(shell date +%FT%T%z)

GO_CMD=go
GO_BUILD=$(GO_CMD) build -v
GO_BUILD_RACE=$(GO_CMD) build -race
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_INSTALL=$(GO_CMD) install -v
GO_CLEAN=$(GO_CMD) clean
GO_DEPS=glide install
GO_DEPS_UPDATE=glide update
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
GO_LINT=golint

# Packages
TOP_PACKAGE_DIR := github.com/bjwschaap
PACKAGE_LIST := logcombiner

LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

.PHONY: all build build-race test test-verbose deps update-deps install clean fmt vet lint release

all: build

build:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Build $$p ..."; \
		$(GO_BUILD) ${LDFLAGS} $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

release:
	@mkdir -p release
	@echo "==> Version: ${VERSION} Build: ${BUILD}"
	@for p in $(PACKAGE_LIST); do \
		echo "==> Making linux AMD64 release for $$p ..."; \
		GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o release/$$p-linux-x64 $(TOP_PACKAGE_DIR)/$$p || exit 1; \
		echo "==> Making linux 386 release for $$p ..."; \
		GOOS=linux GOARCH=386 go build ${LDFLAGS} -o release/$$p-linux-386 $(TOP_PACKAGE_DIR)/$$p || exit 1; \
		echo "==> Making linux ARM release for $$p ..."; \
		GOOS=linux GOARCH=arm go build ${LDFLAGS} -o release/$$p-linux-arm $(TOP_PACKAGE_DIR)/$$p || exit 1; \
		echo "==> Making OSX x64 release for $$p ..."; \
		GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o release/$$p-osx-x64 $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

build-race: vet
	@for p in $(PACKAGE_LIST); do \
		echo "==> Build $$p ..."; \
		$(GO_BUILD_RACE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

test: deps
	@for p in $(PACKAGE_LIST); do \
		echo "==> Unit Testing $$p ..."; \
		$(GO_TEST) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

test-verbose: deps
	@for p in $(PACKAGE_LIST); do \
		echo "==> Unit Testing $$p ..."; \
		$(GO_TEST_VERBOSE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

deps:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Install dependencies for $$p ..."; \
		$(GO_DEPS) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

update-deps:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Update dependencies for $$p ..."; \
		$(GO_DEPS_UPDATE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

install:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Install $$p ..."; \
		$(GO_INSTALL) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

clean:
	@rm -rf release
	@for p in $(PACKAGE_LIST); do \
		echo "==> Clean $$p ..."; \
		$(GO_CLEAN) $(TOP_PACKAGE_DIR)/$$p; \
	done

fmt:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Formatting $$p ..."; \
		$(GO_FMT) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done
vet:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Vet $$p ..."; \
		$(GO_VET) $(TOP_PACKAGE_DIR)/$$p; \
	done

lint:
	@for p in $(PACKAGE_LIST); do \
		echo "==> Lint $$p ..."; \
		$(GO_LINT) src/$(TOP_PACKAGE_DIR)/$$p; \
	done

# vim: set noexpandtab shiftwidth=8 softtabstop=0:
