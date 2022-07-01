NAME				= webfmwk

GO_CC				= go

COVER_FILE  = .coverage.out
COVER_VALUE = -cover -coverprofile=$(COVER_FILE) -covermode=atomic
TEST_COVER  = $(COVER_VALUE)
TEST_ARGS   = -v
TEST_FILTER = #-run <pattern>
app_name		= github.com/burgesQ/$(NAME)
TEST_FILES  = `go list ./...`

TEST        = $(GO_CC) test $(TEST_COVER) $(TEST_ARGS) $(TEST_FILES) $(TEST_FILTER)

LINT				= golangci-lint run

VET					= $(GO_CC) vet .

TIDY				= $(GO_CC) mod tidy

GODOC				= godoc -http=:6060

.PHONY: all
all: $(NAME)

$(NAME): test

.PHONY: vet
vet:
	@ $(VET)

.PHONY: lint
lint:
	@ $(LINT)

.PHONY: test
test:
	 $(TEST)

.PHONY: tidy
tidy:
	@	$(TIDY)

.PHONY: godoc
godoc:
	@ $(GODOC)
