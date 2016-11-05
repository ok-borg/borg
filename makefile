VERSION := $(shell git describe --tags)
# This is just used in all, so as to have something in the arch and OS
OS := $(shell uname -s)
ARCH := $(shell uname -m)

all:
	go get -d
	go build -ldflags "-X main.versionNumber=${VERSION} -X main.operatingSystem=${OS} -X main.architecture=${ARCH}"
# TODO learn for loop in makefile
release:
	go get -d
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=darwin -X main.architecture=amd64" -o borg_darwin_amd64
	GOOS=darwin GOARCH=386 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=darwin -X main.architecture=386" -o borg_darwin_386
	GOOS=freebsd GOARCH=386 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=freebsd -X main.architecture=386" -o borg_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=freebsd -X main.architecture=amd64" -o borg_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=freebsd -X main.architecture=arm" -o borg_freebsd_arm
	GOOS=linux GOARCH=386 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=linux -X main.architecture=386" -o borg_linux_386
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=linux -X main.architecture=amd64" -o borg_linux_amd64
	GOOS=linux GOARCH=arm go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=linux -X main.architecture=arm" -o borg_linux_arm
	GOOS=netbsd GOARCH=386 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=netbsd -X main.architecture=386" -o borg_netbsd_386
	GOOS=netbsd GOARCH=amd64 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=netbsd -X main.architecture=amd64" -o borg_netbsd_amd64
	GOOS=netbsd GOARCH=arm go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=netbsd -X main.architecture=arm" -o borg_netbsd_arm
	GOOS=openbsd GOARCH=386 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=openbsd -X main.architecture=386" -o borg_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=openbsd -X main.architecture=amd64" -o borg_openbsd_amd64
	GOOS=windows GOARCH=386 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=windows -X main.architecture=386" -o borg_windows_386
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.versionNumber=${VERSION} -X main.operatingSystem=windows -X main.architecture=amd64" -o borg_windows_amd64
	# upx does not work on some arch/OS combos
	upx borg_*
clean:
	rm borg*
	