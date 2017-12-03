FROM ubuntu:16.04

# Installing ppa
RUN apt-get update && apt-get install -y software-properties-common apt-transport-https

# Adding oracle java repository
RUN add-apt-repository ppa:webupd8team/java -y
# Add source for sbt (needed for building marathon)
RUN echo "deb https://dl.bintray.com/sbt/debian /" | tee -a /etc/apt/sources.list.d/sbt.list
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 642AC823
RUN apt-get update

# Add license agreement into debconf
RUN echo oracle-java7-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections

RUN apt-get -y install oracle-java7-installer

# Adding gradle 2.3 and maven
RUN add-apt-repository -y ppa:cwchien/gradle && \
    apt-get update && \
    apt-get -y install gradle maven sbt


# The following is added to support zookeeper and cassandra build
RUN apt-get -y install wamerican autoconf libcppunit-dev libtool
