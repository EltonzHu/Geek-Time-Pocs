GOCMD=GO
#GOGET=$(GOCMD) get
#GOCLEAN=$(GOCMD) clean
GOBUILD=$(GOCMD) build

BINDIR=bin
CMDDIR=./cmd
APP_NAME=httpserv
DOCKER_ACCOUNT=xxx

GO_IMAGE?=golang:1.19.2-alpine3.16
RUNTIME_IMAGE?=alpine:3.16

# VERBOSE_FLAG =
# ifeq ($(VERBOSE), yes)
# 	VERBOSE_FLAG := -v
# endif

# Getting Release tag
# VERSION = $(shell git ls-remote --refs --tags git@https://github.com/EltonzHu/Geek-Time-Pocs.git \
# 			| cut -d '/' -f 3 \
# 			| sort --version-sort \
# 			| tail -n1)
BUILD_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT_SHA = $(shell git rev-parse HEAD)
#BUILD_SHA_SHORT = $(shell git rev-parse --short HEAD)
BUILD_DATETIME = $(shell date +"%Y-%m-%d_%H:%M")

LD_FLAGS = "-s \
			-w \
			-X main.commitHash=$(BUILD_COMMIT_SHA) \
			-X main.scmBranch=$(BUILD_BRANCH) \
			-X main.buildDate=$(BUILD_DATETIME)"

define announce
	@echo "# $(1)"; echo
endef

all: fmt vet test testrace
fmt:
	@echo "fmt code"
	@go fmt ./...
vet:
	@echo "vet code"
	@go vet ./...
test:
	@echo "run go testing"
	@go test ./...
testrace:
	@echo "run test race conditions"
	@go test -race -cpu 1,4 -timeout 7m ./...

local-build: all
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags $(LD_FLAGS) -o $(BINDIR)/$(APP_NAME) -v $(CMDDIR)

build-image: all
	$(call announce, $@)
	@echo "build docker image"
	@docker build -t $(APP_NAME) -f ./Dockerfile . \
		--build-arg GO_IMAGE=$(GO_IMAGE) \
		--build-arg RUNTIME_IMAGE=$(RUNTIME_IMAGE) \
		--build-arg LD_FLAGS=$(LD_FLAGS)

upload-image: build-image
	@echo "upload image $(APP_NAME)"
	@docker tag $(APP_NAME) $(DOCKER_ACCOUNT)/$(APP_NAME):$(BUILD_DATETIME)
	@docker push $(DOCKER_ACCOUNT)/$(APP_NAME):$(BUILD_DATETIME)

clean:
	@rm -rf ./$(BINDIR)

.PHONY: \
	all \
	fmt \
	vet \
	test \
	testrace \
	local-build \
	build-image \
	clean 
