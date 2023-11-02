# non-phony targets
REGISTRY_PASSWORD := $(shell echo $REGISTRY_PASSWORD)
REGISTRY_USERNAME := $(shell echo $REGISTRY_USERNAME)
OS := $(shell make install-bin-deps)
ALL_SH_SCRIPT=$(shell find . -type f -name "*.sh" -exec basename {} \;)

test-all:
#	go test ./...

docker-build:
	docker build --force-rm --tag port:1.0 --label "head=port" .
	#docker tag port:1.0 sample-app

docker-run:
	docker run -it --publish 127.0.0.1:8090:8090/tcp --label "head=port" --restart=unless-stopped --rm port:1.0

docker-push: load_envs
	echo $(REGISTRY_PASSWORD) | docker login --username $(REGISTRY_USERNAME) --password-stdin
	docker tag port:1.0 date:$(date)
	docker push port:1.0 --all-tags

install-bin-deps: install-shellcheck

install-shellcheck:
	INSTALLED=$(which shellcheck >> /dev/null; echo $?)
    ifeq ($(INSTALLED),"0")
    	$(shell @echo shellcheck already installed)
   	else
		OS := $(shell uname)
		ifeq ($(OS),Darwin)
			$(shell brew install shellcheck)
		else ifeq ($(OS),Linux)
			$(shell sudo apt install shellcheck)
		else
			$(shell echo "Cannot install shellcheck. OS not Darwin or Linux.")
		endif
	endif

load_envs:
	[ ! -f .env ] && (export $(cat .env | xargs))
	echo $REGISTRY_USERNAME

run: docker-build docker-run

shellcheck: install-bin-deps
	$(foreach i, $(ALL_SH_SCRIPT), shellcheck --color=auto --format=diff $(i);)


.PHONY: clean
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