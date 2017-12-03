APP = dpe-stream
include ${MAKEPATH}/Makefile

SERVICE_PATH=go/src/github.com/northstar/dpe-stream
SERVICE_MARATHON_FILE=marathon.dpestream.worker.json
TMP_FILE=tmp.json

marathon: marathon_app
	cd $(ROOT)/$(SERVICE_PATH) && cp marathon.json $(CURDIR)/$(SERVICE_MARATHON_FILE) ; cd $(CURDIR)
	echo  "{ \"env\" : { \"DPE_STREAM_WORKER_MARATHON_JSON\" : $$($(JQ_BIN) --slurp --raw-input '.' $(SERVICE_MARATHON_FILE)) } }" | sed 's#\\n##g' > $(TMP_FILE)
	$(JQ_BIN) -s '.[0] * .[1] * .[2]' $(MARATHON_BASE_TARGET) $(MARATHON_APP_TARGET) $(TMP_FILE) > marathon.json
	rm -Rf $(TMP_FILE) $(SERVICE_MARATHON_FILE)
