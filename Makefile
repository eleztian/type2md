GO_BIN_PATH := $(if $(GOBIN),$(GOBIN),$(GOPATH)/bin)

VERSION ?= dev
BUILD_TIME ?= `date "+%Y-%m-%d %H:%M:%S"`
COMMIT_ID ?= `git rev-parse --short HEAD`

GO_BUILD := go build --trimpath --ldflags "-w -s \
-X 'main.Version=$(VERSION)' \
-X 'main.CommitID=$(COMMIT_ID)' \
-X 'main.BuildTime=$(BUILD_TIME)'"

all:
	$(GO_BUILD) -o type2md ./*.go
	mv type2md $(GO_BIN_PATH)/type2md