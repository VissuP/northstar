FROM ubuntu:17.04

ENV JRE_TAR jre-8u144-linux-x64.tar.gz
ENV JRE jre1.8.0_144 
ENV JAVA_HOME /opt/$JRE
ENV PATH $PATH:$JAVA_HOME/bin

RUN \
    apt-get update && \
    apt-get install -y wget && \
    wget -q --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u144-b01/090f390dda5b47b9b721c7dfaa008135/$JRE_TAR && \
    tar xzf $JRE_TAR && \
    mv $JRE /opt && \
    rm -rf $JRE_TAR && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get purge -y --auto-remove wget