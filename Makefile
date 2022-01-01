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

.PHONY: test
test:
	ginkgo ./...

.PHONY: e2e
e2e:
	e2e/run.sh

.PHONY: run
run:
	go run cmd/spexec/main.go $(ARGS)

.PHONY: build
build:
	goreleaser build --single-target --snapshot --rm-dist

.PHONY: dedebugo
dedebugo:
	dedebugo --exclude build --exclude dist .

.PHONY: deps
deps:
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
	go install github.com/autopp/dedebugo/cmd/dedebugo@latest
	go install github.com/goreleaser/goreleaser@latest
