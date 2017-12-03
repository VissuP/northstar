FROM ubuntu:16.04
MAINTAINER     Alok Asok <alok.asok@verizon.com>
USER root

RUN apt-get update
RUN apt-get install -y jq

ADD jq.sh /bin/
ENTRYPOINT ["/bin/jq.sh"]
