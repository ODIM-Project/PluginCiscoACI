.PHONY: dep copy down-containers

COPY = caphandler capmessagebus capmiddleware capmodel capresponse caputilities config go.mod go.sum main.go


copy: 
	$(foreach var,$(COPY),cp -a $(var) build/;)

dep: copy
	build/makedep.sh

build-containers: dep
	cd build && docker-compose build

standup-containers: build-containers
	cd build && docker-compose up -d && docker exec -d build_plugin_1 /aci-plugin/command.sh && docker restart build_plugin_1

down-containers:
	cd build && docker-compose down

all: standup-containers

clean: 
	build/cleanup.sh