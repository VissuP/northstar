APP=logger
TAG=$(LOGGER_TAG)

include $(MAKEPATH)/Makefile

GROUP = dakota
BIMAGE ?= $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(GROUP)-$(APP):$(TAG)
JIMAGE ?= $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(GROUP)-java-$(APP):$(TAG)
UIMAGE ?= $(DOCKER_REGISTRY_HOST_PORT)/$(DOCKER_USER)/$(GROUP)-ubuntu-$(APP):$(TAG)

bimage:
		docker build -t $(BIMAGE) .

jimage:
		sed \
			-e 's/LOGGER_TAG/'"$(LOGGER_TAG)"'/' \
			-e 's/BASE_REGISTRY_URI/$(BASE_REGISTRY_URI)/' \
			Dockerfile.java  > Dockerfile.java.tmp; \
		if [ -a $(FILEPATH) ];\
		then \
			IFS=$$'\n';\
			for line in `cat $(FILEPATH)| egrep -v "^[ ]*$|^[ ]*#"`; do echo $${line} ; echo "ENV $${line}" >> Dockerfile.java.tmp; done;\
		else \
			echo "$(FILEPATH) not found"; \
		fi
		docker build -t $(JIMAGE) -f ./Dockerfile.java.tmp . || rm Dockerfile.java.tmp
		mv Dockerfile.java.tmp Dockerfile.java.generated

uimage:
		sed \
			-e 's/LOGGER_TAG/'"$(LOGGER_TAG)"'/' \
			-e 's/BASE_REGISTRY_URI/$(BASE_REGISTRY_URI)/' \
			Dockerfile.ubuntu_16_04  > Dockerfile.ubuntu_16_04.tmp; \
		if [ -a $(FILEPATH) ];\
		then \
			IFS=$$'\n';\
			for line in `cat $(FILEPATH)| egrep -v "^[ ]*$|^[ ]*#"`; do echo $${line} ; echo "ENV $${line}" >> Dockerfile.ubuntu_16_04.tmp; done;\
		else \
			echo "$(FILEPATH) not found"; \
		fi
		docker build -t $(UIMAGE) -f ./Dockerfile.ubuntu_16_04.tmp . || rm Dockerfile.ubuntu_16_04.tmp
		mv Dockerfile.ubuntu_16_04.tmp Dockerfile.ubuntu_16_04.generated

bpush:
	docker push $(BIMAGE)

jpush:
	docker push $(JIMAGE)

upush:
	docker push $(UIMAGE)

image:
		@echo no-op
push:
		@echo no-op

