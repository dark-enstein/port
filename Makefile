# If the .env file exists, include it
-include .env

# Export the variables
export REGISTRY_PASSWORD ?= $(REGISTRY_PASSWORD)
export REGISTRY_USERNAME ?= $(REGISTRY_USERNAME)
ALL_SH_SCRIPT := $(shell find . -type f -name "*.sh" -print0 | xargs -0 -n1 basename)
RED := $(shell tput setaf 1)
GREEN := $(shell tput setaf 2)
ENDCOLOR := $(shell tput sgr0)

.PHONY: test-all docker-build docker-run docker-push install-bin-deps install-shellcheck load_envs run shellcheck clean

test-all:
#	go test ./...

docker-build:
	docker build --force-rm --tag port:1.0 --label "head=port" .
	#docker tag port:1.0 sample-app

docker-run:
	docker run -it --publish 127.0.0.1:8090:8090/tcp --label "head=port" --restart=unless-stopped --rm port:1.0

docker-push:
	@docker login --username $(REGISTRY_USERNAME) --password-stdin <<<$(REGISTRY_PASSWORD); \
	@docker tag port:1.0 date:$(shell date '+%s'); \
	@docker push port:1.0 --all-tags

install-bin-deps: install-shellcheck

install-shellcheck:
	@if ! which shellcheck > /dev/null; then \
		OS=$$(uname); \
		if [ "$$OS" = "Darwin" ]; then \
			brew install shellcheck; \
		elif [ "$$OS" = "Linux" ]; then \
			sudo apt-get update && sudo apt-get install -y shellcheck; \
		else \
			echo "$(RED)Cannot install shellcheck. OS not Darwin or Linux.$(ENDCOLOR)"; \
			exit 1; \
		fi \
	fi

verify:
	echo "$$REGISTRY_PASSWORD"

load_envs:
	$(eval include_env := $(shell cat .env | sed 's/#.*//g' | xargs))
	$(foreach var,$(include_env),$(eval $(var)))

run: docker-build docker-run

shellcheck: install-bin-deps
	@echo "Checking scripts: $(ALL_SH_SCRIPT)"
	@code=0; \
	for i in $(ALL_SH_SCRIPT); do \
		shellcheck --color=auto "$$i" || code=$$?; \
	done; \
	if [ $$code -eq 0 ]; then \
		echo "$(GREEN)Scripts verified$(ENDCOLOR)"; \
	else \
		echo "$(RED)Some scripts have issues$(ENDCOLOR)"; \
		exit $$code; \
	fi


clean:
	@echo Removing the following dangling Docker images ...
	@sleep 1
	@docker images -f dangling=true
	@sleep 1
	@docker image prune --force
	@sleep 1
	@echo Dangling images deleted.

####tod-o:
#support testing
#support other go commands
#makefile tpl: https://github.com/TheNetAdmin/Makefile-Templates/blob/master/SmallProject/Template/Makefile