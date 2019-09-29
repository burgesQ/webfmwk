NAME = webfmwk

CC = go

COVER_FILE = .coverage.out

VET = $(CC) vet
TEST_ARGS = -cover -v -short -coverprofile=$(COVER_FILE) -covermode=atomic
TEST = $(CC) test $(TEST_ARGS)

all: $(NAME)

$(NAME): test
	@git status

vet:
	$(VET) .

lint:
	golint .

test:
	$(TEST) .

.PHONY: vet lint test
