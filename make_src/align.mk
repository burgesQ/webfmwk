ALIGNER		?= betteralign
ALIGN_OPT	?= # -apply

.PHONY: align
align: ## Run better align to analyse struct fields memory alignement
	${ALIGNER} -test_files ${ALIGN_OPT} ./...

.PHONY: align-apply
align-apply: ## Fix found struct fields memory miss-alignement
	$(MAKE) align ALIGN_OPT=-apply

install-align: $(GOPATH)/bin/${ALIGNER} ## Install betteralign analyse tool
$(GOPATH)/bin/${ALIGNER}:
	go install github.com/dkorunic/betteralign/cmd/betteralign@latest
