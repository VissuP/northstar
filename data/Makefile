APP = data
include $(MAKEPATH)/Makefile

CQL_FILE_LIST=$(shell ls ./schemas/*.cql)

comma:= ,
empty:=
space:= $(empty) $(empty)
str :=$(foreach var, $(SEEDS_ENV),'$(var)':$(CASSANDRA_CLUSTER_RF))
str :=$(subst $(space),$(comma),$(str))

marathon:update_rf

update_rf:
	if [ "" = "$(SEEDS_ENV)" ]	; then \
		echo ==============================================================; \
		echo CQL_FILE_LIST= $(CQL_FILE_LIST); \
		echo ==============================================================; \
		for file in $(CQL_FILE_LIST)	; do \
			sed  "s/'replication_factor': '[[:digit:]]'/'replication_factor': '$(CASSANDRA_CLUSTER_RF)'/" ./$$file > ./$$file.bk ;\
			mv ./$$file.bk ./$$file ;\
			sed  "s/'replication_factor' : [[:digit:]]/'replication_factor': $(CASSANDRA_CLUSTER_RF)/" ./$$file > ./$$file.bk ;\
			mv ./$$file.bk ./$$file ;\
			echo Updated $$file with cassandra replication_factor $(CASSANDRA_CLUSTER_RF); \
			cat ./$$file |grep replication_factor || true; \
		done;\
	else\
		echo ==============================================================; \
		echo CQL_FILE_LIST= $(CQL_FILE_LIST); \
		echo ==============================================================; \
		for file in $(CQL_FILE_LIST)	; do \
			sed  "s/replication.*=.*.{ *'class' *:.*.}/replication = {'class':'NetworkTopologyStrategy',$(str)}/g" ./$$file > ./$$file.bk ;\
			mv ./$$file.bk ./$$file ;\
			echo Updated $$file with cassandra cluster $(str); \
			cat ./$$file |grep NetworkTopologyStrategy || true; \
		done;\
	fi;
