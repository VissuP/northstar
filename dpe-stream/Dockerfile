FROM BASE_REGISTRY_URI/DAKOTA_LOGGER_IMAGE:LOGGER_TAG
MAINTAINER Eugen Feller <eugen.feller@verizon.com>

ADD dpe-stream /usr/local/bin/
ADD entrypoint.sh /

EXPOSE 8083

#Config-Auto-Generated ENV
ENV NOR_BUILD_ENV=dev NOR_TWHEEL_CHANNEL_SIZE=100
ENTRYPOINT [ "/entrypoint.sh" ]
