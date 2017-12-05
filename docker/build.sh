#!/bin/bash
#
# Create the builders
#
DOCKER_DIR=$ROOT/docker

cd $DOCKER_DIR/golang-builder && make build_go && cd $DOCKER_DIR
cd $DOCKER_DIR/jq && make build_jq && cd $DOCKER_DIR
cd $DOCKER_DIR/npm-builder && make build_npm && cd $DOCKER_DIR
cd $DOCKER_DIR/logger && make compile && make bimage && make uimage && cd $DOCKER_DIR
cd $DOCKER_DIR/java-builder/oracle-java-8 && make jimage && cd $DOCKER_DIR
