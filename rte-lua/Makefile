APP = rte-lua
include $(MAKEPATH)/Makefile

SUBDIRS = botc

marathon:
	for subdir in ${SUBDIRS}; do $(MAKE) -C $(BASE_COMPILE_DIR)/$$subdir $@; done;

deploy:
	for subdir in ${SUBDIRS}; do $(MAKE) -C $(BASE_COMPILE_DIR)/$$subdir $@; done;

undeploy:
	curl $(CURL_OPTS) -X DELETE $(MARATHON_PARAMS) $(MARATHON_URL)/$(APPID) $(MARATHON_USER_AUTH)

upgrade_cmd:
	for subdir in ${SUBDIRS}; do $(MAKE) -C $(BASE_COMPILE_DIR)/$$subdir $@; done;
