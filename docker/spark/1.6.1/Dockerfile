FROM docker.io/efeller/dakota-ubuntu-logger:0.3-2.4
MAINTAINER Eugen Feller <eugen.feller@verizon.com>

# Specify the Java version. Change to update the JDK.
ENV VERSION 8
ENV UPDATE 152

# Update PATH.
ENV JAVA_HOME /usr/lib/jvm/java-${VERSION}-oracle

RUN dpkg --add-architecture i386 && \
    apt-get update && apt-get install ca-certificates wget curl libncurses5:i386 libstdc++6:i386 zlib1g:i386 git-core -y --no-install-recommends

RUN curl -s http://d3kbcqa49mib13.cloudfront.net/spark-1.6.1-bin-hadoop2.6.tgz | tar -xz -C /opt/
RUN cd /opt && ln -s spark-1.6.1-bin-hadoop2.6 spark

# Install JDK
RUN mkdir -p /usr/lib/jvm/ && \
    cd /usr/lib/jvm && \
    wget --no-cookies --no-check-certificate --header "Cookie: oraclelicense=accept-securebackup-cookie" "http://download.oracle.com/otn-pub/java/jdk/8u152-b16/aa0333dd3019491ca4f6ddbe78cdb6d0/jdk-8u152-linux-x64.tar.gz" && \
    tar xzvf jdk-"${VERSION}"u"${UPDATE}"-linux-x64.tar.gz && \
    mv jdk1.${VERSION}.0_${UPDATE} "${JAVA_HOME}" && \
    apt-get autoclean && apt-get --purge -y autoremove && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN update-alternatives --install "/usr/bin/java" "java" "${JAVA_HOME}/bin/java" 1 && \
	update-alternatives --install "/usr/bin/javaws" "javaws" "${JAVA_HOME}/bin/javaws" 1 && \
	update-alternatives --install "/usr/bin/javac" "javac" "${JAVA_HOME}/bin/javac" 1 && \
	update-alternatives --set java "${JAVA_HOME}/bin/java" && \
	update-alternatives --set javaws "${JAVA_HOME}/bin/javaws" && \
	update-alternatives --set javac "${JAVA_HOME}/bin/javac"

# Install Mesos
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv E56151BF && \
  echo deb http://repos.mesosphere.io/ubuntu trusty main > /etc/apt/sources.list.d/mesosphere.list && \
  apt-get update && \
  apt-get -y install mesos=0.28.0-2.0.16.ubuntu1404

RUN chmod -R 755 /opt/spark/bin/
ADD entrypoint.sh /
ADD conf /opt/spark/conf/
ADD lib /opt/spark/lib/

ENV SPARK_HOME /opt/spark
ENV PATH $PATH:/opt/spark/bin
ENV MESOS_NATIVE_JAVA_LIBRARY /usr/local/lib/libmesos.so

WORKDIR /opt/spark
ENTRYPOINT [ "/entrypoint.sh" ]
