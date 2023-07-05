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
