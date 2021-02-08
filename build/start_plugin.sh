#!/bin/bash
mkdir -p /var/log/plugin_logs
export logFolder="/var/log/plugin_logs"
ip=`echo $HOSTIP`
RootServiceUUID=$(uuidgen)
sed -i "s#.*RootServiceUUID\": \"\",# \"RootServiceUUID\": \"${RootServiceUUID}\",#" /etc/plugin_config/config.json
sed -i 's/"Host":\s*".*",/"Host": "plugin",/g' /etc/plugin_config/config.json
sed -i 's/"ListenerHost":\s*".*",/"ListenerHost": "plugin",/g' /etc/plugin_config/config.json
sed -i 's@"RootCACertificatePath":\s*".*",@"RootCACertificatePath": "/etc/plugin_certs/rootCA.crt",@g' /etc/plugin_config/config.json
sed -i 's@"PrivateKeyPath":\s*".*",@"PrivateKeyPath": "/etc/plugin_certs/odimra_server.key",@g' /etc/plugin_config/config.json
sed -i 's@"CertificatePath":\s*".*"@"CertificatePath": "/etc/plugin_certs/odimra_server.crt"@g' /etc/plugin_config/config.json
sed -i "s#.*LBHost.*# \"LBHost\": \"${ip}\",#" /etc/plugin_config/config.json
sed -i "s#.*LBPort.*# \"LBPort\": \"45021\"#" /etc/plugin_config/config.json
sed -i 's@"MessageQueueConfigFilePath":\s*".*",@"MessageQueueConfigFilePath": "/etc/plugin_config/platformconfig.toml",@g' /etc/plugin_config/config.json

sed -i "s#.*KServersInfo.*#KServersInfo      = [\"kafka:9092\"]#" /etc/plugin_config/platformconfig.toml
sed -i "s#.*KAFKACertFile.*#KAFKACertFile      = \"/etc/plugin_certs/odimra_kafka_client.crt\"#" /etc/plugin_config/platformconfig.toml
sed -i "s#.*KAFKAKeyFile.*#KAFKAKeyFile      = \"/etc/plugin_certs/odimra_kafka_client.key\"#" /etc/plugin_config/platformconfig.toml
sed -i "s#.*KAFKACAFile.*#KAFKACAFile      = \"/etc/plugin_certs/rootCA.crt\"#" /etc/plugin_config/platformconfig.toml

systemctl enable plugin

while true; do
        sleep 5s
done

