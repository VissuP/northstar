APP = rte-lua-botc

SERVICE_PATH=rte-lua
SERVICE_MARATHON_FILE=marathon.rtelua.json
TMP_FILE=tmp.json

include $(MAKEPATH)/Makefile

marathon: marathon_app
	cd $(ROOT)/$(SERVICE_PATH) && cp marathon.json $(CURDIR)/$(SERVICE_MARATHON_FILE) ; cd $(CURDIR)
	echo  "{ \"env\" : { \"RTE_LUA_MARATHON_JSON\" : $$($(JQ_BIN) --slurp --raw-input '.' $(SERVICE_MARATHON_FILE)) } }" | sed 's#\\n##g' > $(TMP_FILE)
	$(JQ_BIN) -s '.[0] * .[1] * .[2]' $(MARATHON_BASE_TARGET) $(MARATHON_APP_TARGET) $(TMP_FILE) > marathon.json
	rm -Rf $(TMP_FILE) $(SERVICE_MARATHON_FILE)
