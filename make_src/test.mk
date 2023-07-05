#
# Testing
#

# go unit test
TEST_FILES   = ./...
TEST_VERBOSE ?= -v
TEST_ARGS    ?= -shuffle=on # -short -failfast
MAKE_ARGS    = # -cpu=1 -parallel=4 -run <pattern>
COVER_FILE   = cover.cov
COVER_HTML   = cover.html
TEST_JSON    ?= -json
PRETTY_ARGS  ?= | tee out.txt | gotestfmt
GOFMT        = $(GOPATH)/bin/gotestfmt
GO_CC ?= go
.PHONY: test
test: ## Run the go unit test
	${GO_CC} test ${TEST_ARGS} ${MAKE_ARGS} ${TEST_VERBOSE} ${TEST_JSON} ${TEST_FILES} ${PRETTY_ARGS}

test-cover-gen: ${COVER_FILE} ## Run the test and generate the coverage file per package
${COVER_FILE}:
	@$(MAKE) test MAKE_ARGS="${MAKE_ARGS} -coverpkg github.com/burgesQ/${NAME}/... -covermode atomic -coverprofile ${COVER_FILE}"

test-cover: ${COVER_FILE} ## Run the go unit test with code coverage
	@${GO_CC} tool cover -func=${COVER_FILE} | tail -n1

test-cover-html:  ${COVER_HTML} ## Build the coverage report in html format
${COVER_HTML}: ${COVER_FILE}
	@${GO_CC} tool cover -html=${COVER_FILE} -o ${COVER_HTML}

.PHONY: test-cover-clean
test-cover-clean: ## Clean the test coverage artifacts
	@rm -rf ${COVER_FILE} ${COVER_HTML}

test-clean: test-cover-clean ## Proxy test-cover-clean
clean-test: test-clean ## Proxy test-clean
test-cover-re: test-cover-clean test-cover test-cover-html ## Re-run the test coverage

install-fmt: ${GOFMT} ## Install gotestfmt
${GOFMT}:
	${GO_CC} install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest

.PHONY: tidy
tidy: ; @ $(GO_CC) mod tidy ## Tidy up the deps
