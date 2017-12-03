APP = nssim

include $(MAKEPATH)/Makefile

BUILDER_NPM = $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/npm-builder:0.1-7.9.0
compile: npm
npm:
	docker run --rm \
		-v $(ROOT):$(ROOT) \
		-v $(HOME)/.ssh/:/root/.ssh/ \
		-w $(BASE_COMPILE_DIR)/web \
		$(BUILDER_NPM) \
		sh -c "npm install; ng build --base-href=."
