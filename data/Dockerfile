FROM docker.io/CASSANDRA_TAG
MAINTAINER Eugen Feller <eugen.feller@verizon.com>

ADD schemas /opt/schemas
ADD data /usr/local/bin/
ADD dbsetup.sh /usr/local/bin/
ADD run.sh /usr/local/bin/

EXPOSE 8082

ENTRYPOINT []
CMD ["sh", "-c",  "{ exec /usr/local/bin/run.sh 2>&1 1>&3 3>&- | /usr/local/bin/logger -st=tcp -ost=false; } 3>&1 1>&2 |  /usr/local/bin/logger -st=tcp"]
