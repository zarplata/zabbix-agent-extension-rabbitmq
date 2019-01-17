.PHONY: all clean-all build cleand-deps deps ver

DATE := $(shell git log -1 --format="%cd" --date=short | sed s/-//g)
COUNT := $(shell git rev-list --count HEAD)
COMMIT := $(shell git rev-parse --short HEAD)

BINARYNAME := zabbix-agent-extension-rabbitmq
VERSION := "${DATE}.${COUNT}_${COMMIT}"

LDFLAGS := "-X main.version=${VERSION}"


default: all 

all: clean-all deps build

ver:
	@echo ${VERSION}

clean-all: clean-deps
	@echo Clean builded binaries
	rm -rf .out/
	@echo Done

build:
	@echo Build
	go build -v -o .out/${BINARYNAME} -ldflags ${LDFLAGS} *.go
	@echo Done

clean-deps:
	@echo Clean dependencies
	rm -rf vendor/*

deps:
	@echo Fetch dependencies
	dep ensure -v


remove:
	@echo Remove
	rm /usr/bin/${BINARYNAME}
	@echo Done

.PHONY: archlinux-package
archlinux-package:
	@echo Build package
	$(shell ./build-archlinux.sh)
