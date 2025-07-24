all: linux windows

GO_BUILD_OPTS := -trimpath

GIT_TAG := $(shell git tag --points-at HEAD)
BINARY_FILENAME := educabiz-scraper$(if $(GIT_TAG),-$(GIT_TAG))
LINUX_FILENAME := $(BINARY_FILENAME)-linux
WINDOWS_FILENAME := $(BINARY_FILENAME)-windows.exe
MACOS_FILENAME := $(BINARY_FILENAME)-macos

linux:
	go build $(GO_BUILD_OPTS) -o $(LINUX_FILENAME)

windows:
	CGO_ENABLED=1 GOOS=windows CC=x86_64-w64-mingw32-gcc go build $(GO_BUILD_OPTS) -ldflags -H=windowsgui -o $(WINDOWS_FILENAME)

macos:
	go build $(GO_BUILD_OPTS) -o $(MACOS_FILENAME)
