# Copyright (C) 2021 Akira Tanimura (@autopp)
#
# Licensed under the Apache License, Version 2.0 (the “License”);
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an “AS IS” BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
VERSION?=$(shell git rev-parse --short HEAD)
ifeq ($(GOOS),windows)
EXT=.exe
else
EXT=
endif

PRODUCT=spexec
BUILD_DIR=$(CURDIR)/build
TARGET_DIR_NAME=$(PRODUCT)-$(GOOS)-$(GOARCH)
TARGET_DIR=$(BUILD_DIR)/$(TARGET_DIR_NAME)
EXEFILE=$(TARGET_DIR)/$(PRODUCT)$(EXT)
ARTIFACT=$(TARGET_DIR).zip

.PHONY: test
test:
	ginkgo ./...

.PHONY: run
run:
	go run cmd/spexec/main.go $(ARGS)

.PHONY: build
build:
	goreleaser build --single-target --snapshot --rm-dist

.PHONY: release
release: $(ARTIFACT)

$(ARTIFACT): build
	cd $(BUILD_DIR) && zip $@ $(TARGET_DIR_NAME)/*

.PHONY: clean
clean:
	rm -fR $(BUILD_DIR)
