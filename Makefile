DEFENVFILE?=$(ROOT)/make/env/$(DC)-default.mk
ENVFILE?=$(ROOT)/make/env/$(DC)-$(ENV).mk

ifneq ("$(wildcard $(ENVFILE))","")
include $(ENVFILE)
else
include $(DEFENVFILE)
endif

NS_SERVICE_DIRS = \
	kafkamngr \
	processing \
	data \
	cron \
	object \
	portal \
	auth \
	northstarapi \
	rte-lua \
	dpe-stream \
	nssim

ifeq ($(USE_DPE_SPARK), true)
NS_SERVICE_DIRS += dpe-spark
endif

marathon: marathon_ns
marathon_ns:
	for d in $(NS_SERVICE_DIRS); do \
		$(MAKE) -C $$d marathon || exit 1; \
	done

compile: compile_ns
compile_ns:
	for d in $(NS_SERVICE_DIRS); do \
		$(MAKE) -C $$d compile || exit 1; \
	done

build: build_ns
build_ns:
	@echo "Building NS services"
	for d in $(NS_SERVICE_DIRS); do \
		echo $$d; \
		$(MAKE) -C $$d build || exit 1; \
	done

push: push_ns
push_ns:
	@echo "Pushing NS services"
	for d in $(NS_SERVICE_DIRS); do \
		echo $$d; \
		$(MAKE) -C $$d push || exit 1; \
	done

deploy: deploy_ns
deploy_ns:
	@echo "Deploying on northstar nodes in $(DC) $(ENV) $(TAG) $(BUILD_ENV)"
	for d in $(NS_SERVICE_DIRS); do \
		echo $$d; \
		$(MAKE) -C $$d deploy || exit 1; \
	done

undeploy: undeploy_ns
undeploy_ns:
	@echo "Undeploying on northstar nodes in $(DC) $(ENV) $(TAG) $(BUILD_ENV)"
	for d in $(NS_SERVICE_DIRS); do \
		echo $$d; \
		$(MAKE) -C $$d undeploy || exit 1; \
	done

upgrade: upgrade_ns
upgrade_ns:
	@echo "upgradeing on northstar nodes in $(DC) $(ENV) $(TAG) $(BUILD_ENV)"
	for d in $(NS_SERVICE_DIRS); do \
		echo $$d; \
		$(MAKE) -C $$d upgrade || exit 1; \
	done

rollback: rollbacke_ns
rollback_ns:
	@echo "rollbacking on northstar nodes in $(DC) $(ENV) $(TAG) $(BUILD_ENV)"
	for d in $(NS_SERVICE_DIRS); do \
		echo $$d; \
		$(MAKE) -C $$d rollback || exit 1; \
	done

rollbacktags: rollbacktagse_ns
rollbacktags_ns:
	@echo "rollbacktagsing on northstar nodes in $(DC) $(ENV) $(TAG) $(BUILD_ENV)"
	for d in $(NS_SERVICE_DIRS); do \
		echo $$d; \
		$(MAKE) -C $$d rollbacktags || exit 1; \
	done

print-% : ; @echo $* = $($*)
