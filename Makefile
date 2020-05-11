# Directory to place `go install`ed binaries into.
export GOBIN ?= $(shell pwd)/bin

GOLINT = $(GOBIN)/golint
GEN_ATOMICINT = $(GOBIN)/gen-atomicint
GEN_VALUEWRAPPER = $(GOBIN)/gen-valuewrapper

GO_FILES ?= $(shell find . '(' -path .git -o -path vendor ')' -prune -o -name '*.go' -print)

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -race ./...

.PHONY: gofmt
gofmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	gofmt -e -s -l $(GO_FILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || (echo "gofmt failed:" && cat $(FMT_LOG) && false)

$(GOLINT):
	go install golang.org/x/lint/golint

$(GEN_VALUEWRAPPER): $(wildcard ./internal/gen-valuewrapper/*)
	go build -o $@ ./internal/gen-valuewrapper

$(GEN_ATOMICINT): $(wildcard ./internal/gen-atomicint/*)
	go build -o $@ ./internal/gen-atomicint

.PHONY: golint
golint: $(GOLINT)
	$(GOLINT) ./...

.PHONY: lint
lint: gofmt golint generatenodirty

.PHONY: cover
cover:
	go test -coverprofile=cover.out -coverpkg ./... -v ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: generate
generate: $(GEN_ATOMICINT) $(GEN_VALUEWRAPPER)
	go generate ./...

.PHONY: generatenodirty
generatenodirty:
	@[ -z "$$(git status --porcelain)" ] || ( \
		echo "Working tree is dirty. Commit your changes first."; \
		exit 1 )
	@make generate
	@status=$$(git status --porcelain); \
		[ -z "$$status" ] || ( \
		echo "Working tree is dirty after `make generate`:"; \
		echo "$$status"; \
		echo "Please ensure that the generated code is up-to-date." )
