include $(MAKEPATH)/Makefile
APP = java-builder
TAG = 8

jimage:
	docker build -t  $(APP):$(TAG) .

jpush:
	docker tag  $(APP):$(TAG) $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)
	docker push $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)

