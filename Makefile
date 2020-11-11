.PHONY: all
all: native native_clean clean

# The binary to build, path to dir under cmd
BIN :=
# Enable CGO to build lib C
CGO := 0
# This repo's root import path (under GOPATH).
PKG := video5pm-api

SHELL ::= /bin/bash

# Check if current dir is git repo?
ifeq ($(wildcard $(if $(shell git rev-parse --show-toplevel 2>/dev/null),$(shell git rev-parse --show-toplevel)/,).git),)
	# This version-strategy uses a manual value to set the version string
	VERSION := v1.0-snapshot
else
	# This version-strategy uses git tags to set the version string
	VERSION := $(shell git describe --tags --always --dirty="-dev")
endif

# check os and arch
GOOS	= linux
GOARCH	= amd64

UNAME	:= $(shell uname)
OS		:= $(shell echo $(UNAME) |tr '[:upper:]' '[:lower:]')
ARCH	:= $(shell uname -m)
ifeq ($(OS), darwin)
	GOOS = darwin
endif

ifneq ($(ARCH), x86_64)
	GOARCH = 386
endif

# Replace backslashes with forward slashes for use on Windows.
# Make is !@#$ing weird.
E		:=
BSLASH 	:= \$E
FSLASH 	:= /

# Directories
WD			:= $(subst $(BSLASH),$(FSLASH),$(shell pwd))
MD			:= $(subst $(BSLASH),$(FSLASH),$(shell dirname "$(realpath $(lastword $(MAKEFILE_LIST)))"))
PKGDIR		= $(MD)
CMDDIR		= $(PKGDIR)/cmd
DISTDIR		:= $(WD)/dist
CONFDIR		:= $(WD)/conf

# Parameters
CMDPKG		= $(PKG)/cmd
CMDS		:= $(shell find "$(CMDDIR)/" -mindepth 1 -maxdepth 2 -type d | sed 's/ /\\ /g' | xargs -n1 basename)
LDFLAGS		:= -X main.version=$(VERSION)

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

# build functions
_pre_check:
	@echo "[check] ..."
	@if [ -z $(BIN) ]; then echo "Must set BIN"; exit 1; fi
	@echo "[check] done"

_build:
	$(eval OUTPUTDIR := $(DISTDIR)/$(GOOS)-$(GOARCH)/$(BIN))
	$(eval BINNAME := $(shell basename $(BIN)))
	@mkdir -p "$(DISTDIR)"
	@echo "[build] ..."
	@echo "app: $(BIN), version: $(VERSION), os: $(GOOS) $(GOARCH)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO) \
	go build -ldflags "$(LDFLAGS)" -o "$(OUTPUTDIR)/bin/$(BINNAME)" "$(CMDPKG)/$(BIN)"
	@echo "[build] done"

_post_buid: $(GOPATH)/src/video5pm-api/templates/run.sh.in
	$(eval OUTPUTDIR := $(DISTDIR)/$(GOOS)-$(GOARCH)/$(BIN))
	$(eval BINNAME := $(shell basename $(BIN)))
	@echo "[post_build] ..."
	@echo $(VERSION) > "$(OUTPUTDIR)/VERSION"
	@echo "copied version file"
	@sed \
		-e 's|ARG_APP_PID_NAME|$(BINNAME)-$(VERSION)|g' \
		-e 's|ARG_APP_BIN_NAME|$(BINNAME)|g' \
		$(GOPATH)/src/video5pm-api/templates/run.sh.in > "$(OUTPUTDIR)/run.sh"
	@chmod +x "$(OUTPUTDIR)/run.sh"
	@echo "copied run script file"
	@mkdir -p "$(OUTPUTDIR)/conf"
	@if [ -f $(OUTPUTDIR)/conf/conf.toml ]; then \
		echo "found config file: $(OUTPUTDIR)/conf/conf.toml"; \
	else \
		touch $(OUTPUTDIR)/conf/conf.toml; \
		echo "created empty conf.toml file"; \
	fi
	@echo "[post_build] done"

_clean_build:
	$(eval OUTPUTDIR := $(DISTDIR)/$(GOOS)-$(GOARCH)/$(BIN))
	@echo "[clean_build] $(OUTPUTDIR) ..."
	@rm -rf "$(OUTPUTDIR)"
	@echo "[clean_buil] done"

# build native (base on machine env)
native: _pre_check
native: _build
native: _post_buid

native_clean: _pre_check
native_clean: _clean_build

test_app: _pre_check
test_app:
	$(eval APP_PKGS := $(shell cd cmd/$(BIN) && go list ./...))
	@go test -cover $(APP_PKGS)

clean:
	@echo "[clean up] all ..."
	@rm -rf "$(DISTDIR)"
	@echo "[clean up] done"

# cd into the GOPATH to workaround ./... not following symlinks
_allpackages = $(shell ( go list ./... 2>&1 1>&3 | \
	grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) 1>&2 ) 3>&1 | \
	grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)))

# memoize allpackages, so that it's executed only once and only if used
allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)