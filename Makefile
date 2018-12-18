PREFIX					?= $(shell pwd)
GIT_REF 				:= $(shell git log --pretty=format:'%h' -n 1)
CURRENT_USER 			:= $(shell whoami)
VERSION 				:= $(shell cat ./VERSION)

# Exports
export GIT_REF
export VERSION
export CONFIG_PATH

all: build-ctl

new-install:
	@go install github.com/blend/go-sdk/cmd/coverage

cover:
	@coverage

test:
	@go test -timeout 5s ./...

build-ctl:
	@echo "$(VERSION)/$(GIT_REF) >> building blogctl"
	@go install -ldflags="-X github.com/wcharczuk/blogctl/pkg/config.Version=${VERSION} -X github.com/wcharczuk/blogctl/pkg/config.GitRef=${GIT_REF}" .
