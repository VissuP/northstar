FROM java:8-jdk
MAINTAINER Nicholas Iaquinto <nickiaq@gmail.com>

# Gradle
ENV GRADLE_VERSION 4.2.1
ENV GRADLE_SHA b551cc04f2ca51c78dd14edb060621f0e5439bdfafa6fd167032a09ac708fbc0

ENV MAVEN_VERSION 3.5.0
ENV MAVEN_SHA a6866b0a0dd5abe3a51124a758a0a842a24edd0bea42ebc74dbf93ad44eb1960

#WORKDIR /usr/bin/app

RUN cd /usr/lib \
 && curl -fl https://downloads.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip -o gradle-bin.zip \
 && echo "$GRADLE_SHA gradle-bin.zip" | sha256sum -c - \
 && unzip "gradle-bin.zip" \
 && rm "gradle-bin.zip"

RUN cd /usr/lib \
    && curl -fl http://apache.claz.org/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.zip -o maven-bin.zip \
    && echo "$MAVEN_SHA maven-bin.zip" | sha256sum -c - \
    && unzip "maven-bin.zip" \
    && rm "maven-bin.zip"

# Set Appropriate Environmental Variables
ENV GRADLE_HOME /usr/lib/gradle-${GRADLE_VERSION}
ENV MAVEN_HOME //usr/lib/apache-maven-${MAVEN_VERSION}
ENV PATH $PATH:$GRADLE_HOME/bin:$MAVEN_HOME/bin



# Caches
VOLUME ["/root/.gradle/caches", "/usr/bin/app"]

# Default command is "/usr/bin/gradle -version" on /usr/bin/app dir
# (ie. Mount project at /usr/bin/app "docker --rm -v /path/to/app:/usr/bin/app gradle <command>")
#ENTRYPOINT ["gradle"]
#CMD ["-version"]
