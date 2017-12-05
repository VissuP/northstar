#!/bin/bash
#
# Push the builders to the designated repository (e.g. Docker.io)
#
DOCKER_DIR=$ROOT/docker

cd $DOCKER_DIR/golang-builder && make push_go && cd $DOCKER_DIR
cd $DOCKER_DIR/jq && make push_jq && cd $DOCKER_DIR
cd $DOCKER_DIR/npm-builder && make push_npm && cd $DOCKER_DIR
cd $DOCKER_DIR/logger && make bpush && make upush && cd $DOCKER_DIR
cd $DOCKER_DIR/java-builder/oracle-java-8 && make jpush && cd $DOCKER_DIR
