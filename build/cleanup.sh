#!/bin/bash
if [ -a build/docker-compose.yml ]; then
	cd build
	docker-compose down
	LIST=`docker image ls | grep -E 'aci-plugin|aci_plugin_builddep' | awk '{print $3}'`
	docker rmi $LIST
	rm -rf caphandler capmessagebus capmiddleware capmodel capresponse caputilities config go.mod go.sum main.go
        echo "Cleanup Done"
        cd ../
        exit 0
else
	echo "docker-compose.yml doesn't exist, are you in the aci-plugin directory?"
	exit 1
fi
