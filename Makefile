.PHONY: all clean dep build build-darwin build-linux

VERSION := $(shell git describe --tags --abbrev=0)

all: dep clean build sha256-checksum

dep:
	dep ensure

prep-build:
	mkdir -p build

build-darwin: main.go dep prep-build
	cd build; mkdir -p corgi_$(VERSION)_macOS_64-bit
	env GOOS=darwin GOARCH=amd64 go build -o build/corgi_$(VERSION)_macOS_64-bit/corgi
	cd build; tar -czf corgi_$(VERSION)_macOS_64-bit.tar.gz corgi_$(VERSION)_macOS_64-bit
	cd build; rm -rf corgi_$(VERSION)_macOS_64-bit

build-linux: main.go dep prep-build
	cd build; mkdir -p corgi_$(VERSION)_linux_64-bit
	env GOOS=linux GOARCH=amd64 go build -o build/corgi_$(VERSION)_linux_64-bit/corgi
	cd build; tar -czf corgi_$(VERSION)_linux_64-bit.tar.gz corgi_$(VERSION)_linux_64-bit
	cd build; rm -rf corgi_$(VERSION)_linux_64-bit

sha256-checksum: build-darwin
	shasum -a256 build/corgi_$(VERSION)_macOS_64-bit.tar.gz

clean:
	rm -rf build

build: build-darwin build-linux
