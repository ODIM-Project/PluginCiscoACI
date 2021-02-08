#!/bin/bash
cp -r /var/plugin_config/ /etc/ && rm -rf /var/plugin_config/* && /aci-plugin/start_plugin.sh
