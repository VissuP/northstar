include $(MAKEPATH)/Makefile
APP = npm-builder
TAG = 0.1-7.9.0

build_npm:
	docker build -t  $(APP):$(TAG) .

push_npm:
	docker tag  $(APP):$(TAG) $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)
	docker push $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)

