#!/bin/bash

DOCKER_DIR=$ROOT/docker

cd $DOCKER_DIR/golang-builder && make build_go && make push_go && cd $DOCKER_DIR
cd $DOCKER_DIR/jq && make build_jq && make push_jq && cd $DOCKER_DIR
cd $DOCKER_DIR/npm-builder && make build_npm && make push_npm && cd $DOCKER_DIR
cd $DOCKER_DIR/logger && make compile && make bimage && make uimage && make bpush && make upush && cd $DOCKER_DIR
cd $DOCKER_DIR/java-builder/oracle-java-8 && make jimage && make jpush && cd $DOCKER_DIR
