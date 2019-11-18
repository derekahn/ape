# Project related variables
PROJECTNAME=$(shell basename "$(PWD)")
M = $(shell printf "\033[34;1m▶\033[0m")
DONE="\n  $(M)  done ✨"

# Go related variables
GOBASE=$(shell pwd)
GOPATH=$(GOBASE)/vendor:$(GOBASE)
GOBIN=$(GOBASE)/bin
GO111MODULE=on
export GO111MODULE

.SILENT: help
help: Makefile
	@echo "\n Choose a command to run in "$(PROJECTNAME)":\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## clean: Clean build files. Runs `go clean` internally; removes cp.out
.SILENT: clean
clean:
	@echo "  $(M)  Cleaning build cache"
	go clean ./...
	find . -name \*cp.out -type f -delete
	sudo rm -rf ./vendor
	rm -rf ./bin
	@echo $(DONE)

## coverage: Checks code coverage
.SILENT: coverage
coverage:
	@echo "  $(M)  👀 Testing code coverage"
	cd "$$GOBASE"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test ./... -coverprofile cp.out
	@echo $(DONE)

## fmt: Runs gofmt on all source files
.SILENT: fmt
fmt:
	@echo "  $(M) 🏃 gofmt"
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./...); do \
		gofmt -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret
	@echo $(DONE)

## install: Install missing dependencies. Builds binary in ./bin
.PHONY: install
install:
	@echo "  $(M)  Checking if there is any missing dependencies"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get) ./...
	@echo $(DONE)

## test: Runs all the tests
.SILENT: test
test:
	@echo "  $(M)  🏃 all the tests"
	cd "$$GOBASE"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test ./...
	@echo $(DONE)
