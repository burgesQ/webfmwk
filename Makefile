NAME	= webfmwk

# can be go, go1.15.15, docker run -v `pwd`:/app -w app go:1.18 ...
GO_CC	= go

$(NAME): test

.PHONY: all
all: test lint ## Run the test and the linter

.PHONY: tidy
tidy: ; @ $(GO_CC) mod tidy ## Tidy up the deps

.PHONY:		godoc
godoc: ;	@godoc -http=:6060 ## Serve a local godoc

LINT	= golangci-lint run
LINT_ARGS	= --max-same-issues 50

.PHONY: lint
lint: ; @$(LINT) $(LINT_ARGS) ## Run the go linter

.PHONY: lint-fix
lint-fix: ; @$(MAKE) lint LINT_ARGS="$(LINT_ARGS) --fix" ## Run the go linter

#
# Testing
#

# go unit test
TEST_FILES	= ./...
TEST_ARGS		= # -v -short -failfast
MAKE_ARGS = # -cpu=1 -parallel=4 -run <pattern>
COVER_FILE=cover.cov
COVER_HTML=cover.html
TEST_OUT=out.txt
GOFMT	= $(GOPATH)/bin/gotestfmt


.PHONY: test
test: ${GOFMT} ## Run the go unit test
	@go test $(TEST_ARGS) $(MAKE_ARGS) $(TEST_FILES) | tee ${TEST_OUT} | $<

.PHONY: test-verbose
test-verbose: ## Run the go unit test with more verbosity
	@ $(MAKE) test -e MAKE_ARGS="-v"

.PHONY: test-race
test-race: ## Run the go unit test with data race detector
	@ $(MAKE) test -e MAKE_ARGS="-race"

.PHONY: test-msan
test-msan: ## Run the go unit test with memory sanitaizer
	@ CC=clang CCX=clang++ $(MAKE) test -e MAKE_ARGS="-msan"

test-cover: $(COVER_FILE) ## Run the go unit test with code coverage
	 @ go tool cover -func=$(COVER_FILE) | tail -n1

test-cover-html:  $(COVER_HTML) ## Build the coverage report in html format
$(COVER_HTML): $(COVER_FILE)
	 @ go tool cover -html=$(COVER_FILE) -o $(COVER_HTML)

test-cover-gen: $(COVER_FILE) ## Run the test and generate the coverage file per package
$(COVER_FILE):
	@ $(MAKE) test \
		MAKE_ARGS="${MAKE_ARGS} -cover -covermode=atomic -coverprofile $(COVER_FILE)"

.PHONY: test-cover-clean
test-cover-clean: ## Clean the test coverage artifacts
		@ rm -rf $(COVER_FILE) $(COVER_HTML)

test-clean: test-cover-clean ## Proxy test-cover-clean
clean-test: test-clean ## PRoxt test-clean

test-cover-re: test-cover-clean test-cover ## Re-run the test coverage

install-fmt: ${GOFMT} ## Install gotestfmt
${GOFMT}:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest

#
# license generator
#

# license generator binary
LGEN	=  license/cmd/license_generator

lgen: $(LGEN) ## Build the license generator
$(LGEN):
	@cd license/cmd && $(MAKE) re

.PHONY: clgen
clgen: ## Clean the license generator
	@cd license/cmd && $(MAKE) clean

.PHONY: licenseTest
licenseTest: ## Build the test docker image and run it
	@echo " -- Updating sub module " ; git submodule update --init --recursive
	@echo " -- Builder license test docker image" ; cd license && docker build -t license-test --rm .
	@echo " -- Running license test docker image" ; docker run --rm --name license-test license-test || echo " -- Test failed, please check the logs"
	@echo " -- All test succesfuly passed"

.PHONY:		help
help: ## Display this help screen
	@echo "Usage: \n"
	@grep -h -E '^[a-zA-Z\.\-_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
