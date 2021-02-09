#!/bin/bash
docker image ls | grep aci_plugin_builddep > /dev/null 2>&1
if [ ${?} -eq 0 ]; then
        echo "builddep already exists"
        exit 0
else
        cd build && docker build -t aci_plugin_builddep:tst -f Dockerfile.builddep .
        exit 0
fi

