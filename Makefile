
# non-phony targets
docker-build:
	docker build -t port:1.0 .
	#docker tag port:1.0 sample-app

docker-run:
	docker run -it -p 127.0.0.1:8090:8090/tcp port:1.0

docker-push:
	#docker login
	#docker tag port:1.0 sample-app
	#docker push port:1.0

run: docker-build docker-run

.PHONY: clean
clean:
	@echo Removing the following dangling Docker images ...
	@sleep 1
	@docker images -f dangling=true
	@sleep 1
	@docker image prune --force
	@sleep 1
	@echo Dangling images deleted.

#TODO:
#support testing
#support other go commands
#makefile tpl: https://github.com/TheNetAdmin/Makefile-Templates/blob/master/SmallProject/Template/Makefile