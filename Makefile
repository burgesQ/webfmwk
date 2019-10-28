NAME				= webfmwk

CC					= go

COVER_FILE	= .coverage.out
TEST_FILE		= `go list ./... | grep -v example`

TEST_ARGS		= -cover -v -short -coverprofile=$(COVER_FILE) -covermode=atomic
TEST				= $(CC) test $(TEST_ARGS) $(TEST_FILE)

LINT				= golangci-lint run --enable-all -D gochecknoglobals -D lll -D errcheck -D gofmt -D golint -D stylecheck

VET					= $(CC) vet .

all: $(NAME)

$(NAME): test
	@git status

vet:
	$(VET)

lint:
	$(LINT)

test:
	$(TEST)

.PHONY: vet lint test
