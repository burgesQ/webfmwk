LINTER	= golangci-lint
LINT    = ${LINTER} run
LINT_ARGS	= --max-same-issues 50

.PHONY: lint
lint: ## Run the go linter
	@$(LINT) $(LINT_ARGS)

.PHONY: lint-fix
lint-fix: ; @$(MAKE) lint LINT_ARGS="$(LINT_ARGS) --fix" ## Run the go linter
