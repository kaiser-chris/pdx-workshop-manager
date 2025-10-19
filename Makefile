SOURCES := $(shell find . -name '*.go')
MAIN_PACKAGE_PATH := .

ZIG := $(CURDIR)/scripts/zig.sh
ZIG_CC := $(ZIG) cc -w
ZIG_CXX := $(ZIG) c++ -w

LINUXGNU_GOFLAGS := --ldflags '-linkmode external -w' $(COMMON_GOFLAGS)
LINUXGNU_GLIBC_VERSION := 2.17

LINUXMUSL_GOFLAGS := --ldflags '-linkmode external -w -extldflags -static' $(COMMON_GOFLAGS)

WINDOWS_GOFLAGS := $(COMMON_GOFLAGS)

# Always use build with cgo enabled
export CGO_ENABLED = 1

# Not running in a docker container
export  ALLOW_OUTSIDE_DOCKER = 1


.PHONY: all
all: linux windows

.PHONY: clean
clean:
	rm -rf dist

.PHONY: linux
linux: dist/linux-amd64/pdx-workshop-manager

dist/linux-amd64/pdx-workshop-manager: $(SOURCES)
	$(eval export CC = $(ZIG_CC) --target=x86_64-linux-gnu.$(LINUXGNU_GLIBC_VERSION))
	$(eval export CXX = $(ZIG_CXX) --target=x86_64-linux-gnu.$(LINUXGNU_GLIBC_VERSION))
	$(eval export GOOS = linux)
	$(eval export GOARCH = amd64)
	@echo CC="$(CC)" CXX="$(CXX)" GOOS="$(GOOS)" GOARCH="$(GOARCH)"
	go build $(LINUXGNU_GOFLAGS) -o $@ $(MAIN_PACKAGE_PATH)
	cp sdk/redistributable_bin/linux64/libsteam_api.so dist/linux-amd64/

.PHONY: windows
windows: dist/windows-amd64/pdx-workshop-manager.exe

dist/windows-amd64/pdx-workshop-manager.exe: $(SOURCES)
	$(eval export CC = $(ZIG_CC) --target=x86_64-windows-gnu)
	$(eval export CXX = $(ZIG_CXX) --target=x86_64-windows-gnu)
	$(eval export GOOS = windows)
	$(eval export GOARCH = amd64)
	@echo CC="$(CC)" CXX="$(CXX)" GOOS="$(GOOS)" GOARCH="$(GOARCH)"
	go build $(WINDOWS_GOFLAGS) -o $@ $(MAIN_PACKAGE_PATH)
	cp sdk/redistributable_bin/win64/steam_api64.dll dist/windows-amd64/