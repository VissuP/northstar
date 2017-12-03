APP = portal

include $(MAKEPATH)/Makefile

BUILDER_NPM = $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/npm-builder:0.1-7.9.0

compile: npm
npm:
	docker run --rm \
		-v $(HOME)/.ssh/:/root/.ssh/ \
		-v $(ROOT):$(ROOT) \
		-w $(BASE_COMPILE_DIR)/web \
		$(BUILDER_NPM) \
		sh -c "npm install && npm run build:prod"
js:
	@echo "Compiling typescript to javascript"
	cd web; tsc

rmjs:
	@echo "Deleting javascript files"
	find web/app -name "*.js" -type f -delete
	find web/app -name "*.js.map" -type f -delete

clean:
	docker run --rm \
		-v $(ROOT):$(ROOT) \
		-w $(BASE_COMPILE_DIR) \
		$(BUILDER_NPM) \
		rm -Rf $(BASE_COMPILE_DIR)/web/tslint-result rm -Rf $(BASE_COMPILE_DIR)/web/dist
