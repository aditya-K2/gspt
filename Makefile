COUNT := $(shell git rev-list --count HEAD)
SHORT := $(shell git rev-parse --short HEAD)
VERSION := $(shell printf "r%s.%s\n" ${COUNT} ${SHORT})
DATE := $(shell date +%s)
REPO := aditya-K2/gspt
GC := go
VERBOSE =
LDFLAGS := -ldflags="-X github.com/${REPO}/config.Version=${VERSION} -X github.com/${REPO}/config.BuildDate=${DATE}"
BUILD := ${GC} build ${LDFLAGS} ${VERBOSE}

.PHONY: gspt install linux-arm64 linux-amd64 darwin-amd64 darwin-arm64 windows-amd64

all: linux-amd64 linux-arm64 linux-arm darwin-amd64 darwin-arm64 windows-amd64.exe

install: gspt
	install bin/gspt -t "/usr/bin/"

gspt:
	${BUILD} -o bin/gspt

linux-amd64:
	GOOS=linux GOARCH=amd64 \
	${BUILD} -o bin/gspt-linux-amd64

linux-arm64:
	GOOS=linux GOARCH=arm64 \
	${BUILD} -o bin/gspt-linux-arm64

linux-arm:
	GOOS=linux GOARCH=arm \
	${BUILD} -o bin/gspt-linux-arm

darwin-amd64:
	GOOS=darwin GOARCH=amd64 \
	${BUILD} -o bin/gspt-darwin-amd64

darwin-arm64:
	GOOS=darwin GOARCH=arm64 \
	${BUILD} -o bin/gspt-darwin-arm64

windows-amd64.exe:
	GOOS=windows GOARCH=amd64 \
	${BUILD} -o bin/gspt-windows-amd64.exe
