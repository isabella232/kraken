
NAME := jirarelease
ARCH := amd64
VERSION := 1.0
DATE := $(shell date)
COMMIT_ID := $(shell git rev-parse --short HEAD)
SDK_INFO := $(shell go version)
LD_FLAGS := -X main.version $(VERSION) -X main.commit $(COMMIT_ID) -X main.buildTime '$(DATE)' -X main.sdkInfo '$(SDK_INFO)'

all: clean binaries 

test:
	godep go test

binaries: tools deps test 
	GOOS=darwin GOARCH=$(ARCH) godep go build -ldflags "$(LD_FLAGS)" -o $(NAME)-darwin-$(ARCH)
	GOOS=linux GOARCH=$(ARCH) godep go build -ldflags "$(LD_FLAGS)" -o $(NAME)-linux-$(ARCH)
	GOOS=windows GOARCH=$(ARCH) godep go build -ldflags "$(LD_FLAGS)" -o $(NAME)-windows-$(ARCH).exe

deps:
	go get -v github.com/xoom/jira

tools:
	type godep > /dev/null 2>&1 || go get -v github.com/tools/godep

clean: 
	go clean
	rm -f *.deb
