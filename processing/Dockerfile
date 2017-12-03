FROM BASE_REGISTRY_URI/DAKOTA_LOGGER_IMAGE:LOGGER_TAG
MAINTAINER Eugen Feller <eugen.feller@verizon.com>

ADD processing /usr/local/bin/
CMD ["sh", "-c",  "{ exec processing 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 |  /usr/local/bin/logger -st=tcp"]

EXPOSE 8080

#Config-Auto-Generated ENV
ENV NOR_BUILD_ENV=dev NOR_TWHEEL_CHANNEL_SIZE=100
