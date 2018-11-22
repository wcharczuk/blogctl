PREFIX					?= $(shell pwd)
DOCKER_IMAGE_NAME       ?= warden
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
GIT_REF 				:= $(shell git log --pretty=format:'%h' -n 1)
CURRENT_USER 			:= $(shell whoami)
VERSION 				:= $(shell cat ./VERSION)

# Exports
export GIT_REF
export VERSION
export CONFIG_PATH

all: build-ctl

build-ctl:
	@echo "$(VERSION)/$(GIT_REF) >> building blogctl"
	@go install -ldflags="-X github.com/wcharczuk/photoblog/pkg/config.Version=${VERSION} -X github.com/wcharczuk/photoblog/pkg/config.GitRef=${GIT_REF}" ./blogctl
