# NorthStar

<hr>
NorthStar is an advanced data processing and visualization platform that follows the serverless computing paradigm. For data processing, NorthStar users develop code snippets that can be used by 
NorthStar to transform data of arbitrary size and form. Code snippets can be run manually, attached to events, and executed in a periodic fashion. For interactive code snippet development and data 
visualization, NorthStar notebooks are used.
</hr>

## Build Prerequisites

* direnv
* go 1.7
* Docker for Mac

## Operations Prerequisites

* Kafka deployed on Mesos
* ZK deployed on Mesos
* Casssandra deployed on Mesos
* Redis deployed on Mesos (optional)

## Getting the Source and setting your GOPATH

Create a working directory and set your GOPATH to that directory
```
$ mkdir ns
$ export GOPATH=$PWD
$ cd $GOPATH
```
Use go tool to get the northstar source tree
```
$ go get -d github.com/verizonlabs/northstar
```

## Setup build environment
Copy the ".envrc" file located under build directory and make edits to suit your environment.  Example:
```
$ cp src/github.com/verizonlabs/northstar/build/.envrc .
$ direnv allow
$ cd $GOPATH/src/github.com/verizonlabs/northstar

$ export CONTACT=<your email>
$ export ENV=<your env name> (e.g., example)
$ export DC=<your dc name> (e.g., dc1)
$ export TAG=<your release tag> (e.g., release-1.0.0)
```

Open your DC file (e.g., make/env/dc1.mk) and add your Docker Hub username to DOCKER_USER.

## Build and push base Docker images

NorthStar relies on a number of *builders* (docker images that contain various compilers and tools) in order to be built.  These *builders* can be created for your environment and tagged with your identifier by running the following script.
```
$ ./docker/build.sh
```
If you wish to push the builders to a repository you can do so by setting the DOCKER_REGISTRY_HOST_PORT env variable (which defaults to docker.io)
```
$ ./docker/push.sh
```

## Build services

```
$ make build && make push
```

# Deploy/undeploy services

```
$ make deploy
$ make undeploy
```

## Authors

* **Eugen Feller** - <eugen.feller@verizon.com>
* **Brian Avery** - <brian.avery@verizon.com>
* **Yagiz Onat Yazir** - <yagiz.yazir@one.verizon.com>
* **Larry Rau** - <lawrence.rau@verizonwireless.com>

Other contributors: Delvis Gomez, Judy Gao, Tirth Shah, Kevin Tabb, Umang Singh, Safique Ahemad, Sandarsh Devappa, Lalit Kumar, Chandar Natarajan, Vaneeta Singh, Atul Gupta, Bader Aljishi.

## License

This project is licensed under the Apache 2.0 license - see the [LICENSE](LICENSE) file for details
