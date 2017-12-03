include $(MAKEPATH)/Makefile
APP = go-builder
TAG = 1.7.3

build_go:
	docker build -t  $(APP):$(TAG) .

push_go:
	docker tag  $(APP):$(TAG) $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)
	docker push $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)

