# Image `golang-builder`

Debian Jessie

 + `Go` 1.8.3
This is the Git repo of the Docker official image for golang. See the Docker Hub page for the full readme on how to use this Docker image and for information regarding contributing and issues.

This image is meant to be used in combination with a Makefile to build a go project, i.e. something like this:

     docker run -v $(S):/go/src/gateway thingservice.verizon.com:5000/golang-builder /bin/bash -c "cd gateway; make all"
