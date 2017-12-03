FROM BASE_REGISTRY_URI/DAKOTA_LOGGER_IMAGE:LOGGER_TAG
MAINTAINER Eugen Feller <eugen.feller@verizon.com>

ADD kafkamngr /usr/local/bin/

CMD ["sh", "-c",  "{ exec kafkamngr 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 |  /usr/local/bin/logger -st=tcp"]

EXPOSE 8083

#Config-Auto-Generated ENV
ENV NOR_BUILD_ENV=dev NOR_TWHEEL_CHANNEL_SIZE=100
ENV MAX_CPUS=4
ENV MAX_MEMORY=2048
ENV STATS_INTERVAL=30
ENV NUM_INSTANCES=4
ENV HTTPCLIENTLB_RESOLVE_MESOS_DNS=false
ENV HTTPCLIENTLB_MESOS_DNS_QUERY_INTERVAL=30
