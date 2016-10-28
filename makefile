VERSION := $(shell git describe --tags)
# There are the default architecture values can be changed via 'make OS=linux ARCH=amd64'
OS = darwin
ARCH = amd64
all:
	GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "-X main.versionNumber=${VERSION} -X main.operatingSystem=${OS} -X main.architecture=${ARCH}"
clean:
	rm borg