include $(MAKEPATH)/Makefile
APP = java-builder
TAG = 8

jimage:
	docker build -t  $(APP):$(TAG) .
	docker tag  $(APP):$(TAG) $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)

jpush:
	docker push $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)

