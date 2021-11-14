GOCMD			:=$(shell which go)
GOBUILD			:=$(GOCMD) build


IMPORT_PATH		:=subtitle/cmd
BUILD_TIME		:=$(shell date "+%F %T")
COMMIT_ID       :=$(shell git rev-parse HEAD)
GO_VERSION      :=$(shell $(GOCMD) version)
VERSION			:=$(shell git describe --tags)
#VERSION			:=v0.0.2
BUILD_USER		:=$(shell whoami)
FLAG			:="-X '${IMPORT_PATH}.buildTime=${BUILD_TIME}' -X '${IMPORT_PATH}.commitID=${COMMIT_ID}' -X '${IMPORT_PATH}.goVersion=${GO_VERSION}' -X '${IMPORT_PATH}.goVersion=${GO_VERSION}' -X '${IMPORT_PATH}.Version=${VERSION}' -X '${IMPORT_PATH}.buildUser=${BUILD_USER}'"

BINARY_DIR=bin/subtitle
BINARY_NAME:=subtitle


# linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-linux

#mac
build-darwin:
	CGO_ENABLED=0 GOOS=darwin $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin
# windows
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG)  -o $(BINARY_DIR)/$(BINARY_NAME)-win.exe
# common
build:
	CGO_ENABLED=0  $(GOBUILD) -ldflags $(FLAG)  -o $(BINARY_DIR)/$(BINARY_NAME)

# 全平台
build-all:
	make build-linux
	make build-darwin
	make build-win
	upx $(BINARY_DIR)/*
	cd bin&&tar zcf ${BINARY_NAME}.tgz ${BINARY_NAME}