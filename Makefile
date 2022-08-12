SHELL:=/bin/bash

GO=go
GO_LINT=golangci-lint
GO_VULN=govulncheck

GO_VET_OPTS=-v
GO_TEST_OPTS=-v -race
RM_OPTS=-f

CMD_DIRS:=$(wildcard cmd/*)
CMDS:=$(subst cmd,bin,$(CMD_DIRS))

.SECONDEXPANSION:
bin/%:
	$(GO) build $(GO_BUILD_OPT) -o $@ ./cmd/$*

.SECONDEXPANSION:
docker/%:
	$(DOCKER) image build -f $@ -t $(DOCKER_REPO)/$@:$(DOCKER_TAG) $(DOCKER_CONTEXT)

.PHONY: build
build: $(CMDS)

.PHONY: vuln
vuln:
	$(GO_VULN) ./...

.PHONY: vet
vet:
	$(GO) vet $(GO_VET_OPTS) ./...

.PHONY: test
test:
	$(GO) test $(GO_TEST_OPTS) ./...

.PHONY: lint
lint:
	$(GO_LINT) run

.PHONY: check
check: vet test lint vuln

.PHONY: mod
mod:
	$(GO) mod tidy

.PHONY: clean
clean:
	-$(GO) clean
	-rm $(RM_OPTS) bin/*

.PHONY: all
all: mod check build

.DEFAULT_GOAL=all