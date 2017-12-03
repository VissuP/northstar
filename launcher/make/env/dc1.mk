DOCKER_REGISTRY_HOST_PORT ?= docker.io
DOCKER_REGISTRY_IP_PORT = docker.io
DCOS_URL = http://127.0.0.1:9090

AUTH_TOK = $(shell curl -s -X POST -H 'Content-type: application/json' -d '{"uid":"test", "password":"password"}' http://127.0.0.1:9090/acs/api/v1/auth/login | grep token | awk '{print $$2}')

NGINX_BASE_IMAGE = mon-nginx-base:release-2.2-1

MESOS_DNS_SUFIX = mon-marathon-service.mesos
MESOS_DCOS_DNS_SUFIX = mesos
CROSS_DC_MESOS_SUFFIX=".burlington"
ZK_CONT_PATH = \/zookeeper

# Use inside data center IP for MARATHON_HOST
MARATHON_HOST = 127.0.0.1
MARATHON_PORT = 19090

MARATHON_URL = http://127.0.0.1:19090/v2/apps
MARATHON_AUTH = test:password
MARATHON_AUTH_TYPE = BASIC_AUTH
MARATHON_AUTH_URL = http://127.0.0.1:9090/acs/api/v1/auth/login

# Use inside data center IP for MARATHON_BASE_URL
MARATHON_BASE_URL = http://127.0.0.1:9090/service/mon-marathon-service
MARATHON_USERNAME = test
MARATHON_PASSWORD = password

MARATHON_BASE_PATH = /v2/apps
MARATHON_IS_HTTPS = false

REGISTRY_USERNAME = test
REGISTRY_PASSWORD = password

URIS = file:///docker.tar.gz
