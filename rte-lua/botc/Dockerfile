FROM  BASE_REGISTRY_URI/DAKOTA_UBUNTU_LOGGER_IMAGE:LOGGER_TAG
MAINTAINER  Eugen Feller <eugen.feller@verizon.com>

RUN apt-get update
RUN apt-get install -y ca-certificates

ADD [ "rte-lua-botc", "/usr/local/bin/" ]
CMD ["sh", "-c",  "{ exec rte-lua-botc 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 |  /usr/local/bin/logger -st=tcp"]
