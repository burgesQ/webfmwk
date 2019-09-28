NAME = webfmwk

CC = go

COVER_FILE = .coverage.out

VET = $(CC) vet
TEST_ARGS = -cover -v -short -coverprofile=$(COVER_FILE)
TEST = $(CC) test $(TEST_ARGS)

all: $(NAME)

$(NAME): test
	@git status

vet:
	$(VET) .

test:
	$(TEST) .

.PHONY: vet test
