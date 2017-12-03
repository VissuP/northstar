FROM BASE_REGISTRY_URI/DAKOTA_UBUNTU_LOGGER_IMAGE:LOGGER_TAG
MAINTAINER Eugen Feller <eugen.feller@verizon.com>

RUN apt-get update
RUN apt-get -y install luarocks
RUN luarocks install luasocket

ADD start.sh /usr/local/bin/
ADD worker.sh /usr/local/bin/
ADD rte-lua /usr/local/bin/

CMD ["sh", "-c",  "{ exec /usr/local/bin/start.sh 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 |  /usr/local/bin/logger -st=tcp"]

EXPOSE 8082
