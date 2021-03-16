# PluginCiscoACI
PluginCiscoACI
This repository hosts the code for HPE ACI Plugin

## Docker Container Install
This makes it easier to stand up aci-plugin. You do with a simple "make all" and you will get a aci-plugin container running with a port open on 45003
### Required Sofware to Install

To be able to install with docker you need:
* docker-ce
* docker-ce-cli
* docker-compose
* golang-docker-credential-helpers
* python-docker
* python-dockerpty
* python-dockerpycreds

## Makefile operations

**Build Containers**
This will create all the docker images
```bash
$ make build-containers
```

**Stand up Containers**
This will bring up all containers, which includes fetching all online images
```bash
$ make standup-containers
```

**Bring down Containers**
This will just be like docker-compose down, which removes all containers but keeps the images
```bash
$ make down-containers
```

**Make All**
Will build all images and bring up the continaers
```bash
$ make all
```

**Clean everything up except the dependency image**
```bash
$ make clean

**Clean everything up including the dependency image & log files**
```bash
$ make deepclean