# Copyright (C) 2016 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO := GO15VENDOREXPERIMENT=1 go
VERSION ?= $(shell cat version/VERSION)
REVISION=$(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown')
HOST=$(shell hostname -f)
BUILD_DATE=$(shell date +%Y%m%d-%H:%M:%S)
GO_VERSION=$(shell go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
# PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
# FORMATTED := $(shell $(GO) fmt $(PACKAGE_DIRS))

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR ?= ./out
ORG := github.com/fabric8io
REPOPATH ?= $(ORG)/fabric8-dependency-wait-service
# ROOT_PACKAGE := $(shell go list .)
# $(info ROOT_PACKAGE is $(ROOT_PACKAGE))

# ORIGINAL_GOPATH := $(GOPATH)
#GOPATH := $(shell pwd)/_gopath
# $(info GOPATH is $(GOPATH))

BUILDFLAGS := -ldflags \
  " -X $(ROOT_PACKAGE)/version.Version='$(VERSION)'\
    -X $(ROOT_PACKAGE)/version.Revision='$(REVISION)'\
    -X $(ROOT_PACKAGE)/version.Branch='$(BRANCH)'\
    -X $(ROOT_PACKAGE)/version.BuildUser='${USER}@$(HOST)'\
    -X $(ROOT_PACKAGE)/version.BuildDate='$(BUILD_DATE)'\
    -X $(ROOT_PACKAGE)/version.GoVersion='$(GO_VERSION)'\
    -s -w -extldflags '-static'"
# " -s -w -extldflags '-static'"

all:out/fabric8-dependency-wait-service-linux-amd64 out/fabric8-dependency-wait-service-darwin-amd64 out/fabric8-dependency-wait-service-darwin-amd64 out/fabric8-dependency-wait-service-windows-amd64.exe out/fabric8-dependency-wait-service-linux-arm docker

out/fabric8-dependency-wait-service-linux-amd64: 
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build $(BUILDFLAGS) -o $(BUILD_DIR)/fabric8-dependency-wait-service-linux-amd64 $(ROOT_PACKAGE)

out/fabric8-dependency-wait-service-darwin-amd64: 
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build $(BUILDFLAGS) -o $(BUILD_DIR)/fabric8-dependency-wait-service-darwin-amd64 $(ROOT_PACKAGE)

out/fabric8-dependency-wait-service-windows-amd64.exe: 
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build $(BUILDFLAGS) -o $(BUILD_DIR)/fabric8-dependency-wait-service-windows-amd64.exe $(ROOT_PACKAGE)

out/fabric8-dependency-wait-service-linux-arm: 
	CGO_ENABLED=0 GOARCH=arm GOOS=linux go build $(BUILDFLAGS) -o $(BUILD_DIR)/fabric8-dependency-wait-service-linux-arm $(ROOT_PACKAGE)

.PHONY: test
test: 
	go test -v 

.PHONY: release
release: clean test cross
	mkdir -p release
	cp out/fabric8-dependency-wait-service-*-amd64* release
	cp out/fabric8-dependency-wait-service-*-arm* release
	gh-release checksums sha256
	gh-release create fabric8-services/fabric8-dependency-wait-service $(VERSION) master v$(VERSION)

.PHONY: cross
cross: out/fabric8-dependency-wait-service-linux-amd64 out/fabric8-dependency-wait-service-darwin-amd64 out/fabric8-dependency-wait-service-windows-amd64.exe out/fabric8-dependency-wait-service-linux-arm


.PHONY: clean
clean:
	rm -rf out/

.PHONY: docker
docker: out/fabric8-dependency-wait-service-linux-amd64
	docker build -t "fabric8/fabric8-dependency-wait-service:dev" .
