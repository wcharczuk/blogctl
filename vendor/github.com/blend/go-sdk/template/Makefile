all: build

VERSION=v1.6.3
GIT_SHA=$(shell git log --pretty=format:'%h' -n 1)
SHASUMCMD := $(shell command -v sha1sum || command -v shasum; 2> /dev/null)
TARCMD := $(shell command -v tar || command -v tar; 2> /dev/null)

.PHONY: build
build: 
	mkdir -p ./build/dist/darwin
	mkdir -p ./build/dist/linux
	GOOS=darwin GOARCH=amd64 go build -o ./build/dist/darwin/template-darwin-amd64 -ldflags "-X main.Version=${VERSION} -X blendlabs.com/template.GitVersion=${GIT_SHA}" template/main.go
	GOOS=linux GOARCH=amd64 go build -o ./build/dist/linux/template-linux-amd64  -ldflags "-X main.Version=${VERSION} -X blendlabs.com/template.GitVersion=${GIT_SHA}" template/main.go
	(${SHASUMCMD} ./build/dist/darwin/template-darwin-amd64 | cut -d' ' -f1) > ./build/dist/darwin/template-darwin-amd64.sha1
	(${SHASUMCMD} ./build/dist/linux/template-linux-amd64 | cut -d' ' -f1) > ./build/dist/linux/template-linux-amd64.sha1
	${TARCMD} -zcvf ./build/dist/template-darwin-amd64.tar.gz ./build/dist/darwin
	${TARCMD} -zcvf ./build/dist/template-linux-amd64.tar.gz ./build/dist/linux
	rm -rf ./build/dist/darwin
	rm -rf ./build/dist/linux

.PHONY: release-tag
release-tag:
	@git tag ${VERSION}
	@git push --tags

.PHONY: test
test:
	@go test 
