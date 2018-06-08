.PHONY: all dep build build-darwin build-linux

all: dep build

dep:
	dep ensure

prep-build:
	mkdir -p build

build-darwin: main.go dep prep-build
	env GOOS=darwin GOARCH=amd64 go build -o build/corgi_darwin_amd64

build-linux: main.go dep prep-build
	env GOOS=linux GOARCH=amd64 go build -o build/corgi_linux_amd64

build: build-darwin build-linux
