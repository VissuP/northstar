include $(MAKEPATH)/Makefile
APP = jq-runner
TAG = 1.0

build_jq:
	docker build -t  $(APP):$(TAG) .
	docker tag  $(APP):$(TAG) $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)

push_jq:
	docker push $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(APP):$(TAG)

