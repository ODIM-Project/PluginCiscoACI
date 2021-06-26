[![build_deploy_test Actions Status](https://github.com/ODIM-Project/ODIM/workflows/build_deploy_test/badge.svg)](https://github.com/ODIM-Project/ODIM/actions)
[![build_unittest Actions Status](https://github.com/ODIM-Project/ODIM/workflows/build_unittest/badge.svg)](https://github.com/ODIM-Project/ODIM/actions)

# Table of contents

- [Overview of Cisco ACI](#overview-of-cisco-aci)
	- [Deploying the Cisco ACI plugin](#deploying-the-cisco-aci-plugin)
	- [Adding a plugin into the Resource Aggregator for ODIM framework](#adding-a-plugin-into-the-resource-aggregator-for-odim-framework)
	- [Cisco ACI fabric APIs](#Cisco-ACI-fabric-APIs)
	- [Configuring proxy server for a plugin version](#configuring-proxy-server-for-a-plugin-version)
	- [Plugin configuration parameters](#plugin-configuration-parameters)

# Overview of Cisco ACI

Cisco ACI (Application Centric Infrastructure) is an open ecosystem model that uses a holistic, systems-based approach to integrate hardware and software, and the physical and virtual elements, to enable unique business value for modern data centers.

Resource Aggregator for ODIM supports Cisco ACI plugin that can abstract, translate, and expose southbound resource information to the resource aggregator through
RESTful APIs.

## Deploying the Cisco ACI plugin

**Prerequisites**

Kubernetes cluster is set up and the resource aggregator is successfully deployed.

1. Create a directory called `plugins` on the deployment node.

   ```
   $ mkdir plugins
   ```

3. In the `plugins` directory, create a directory called `aciplugin`.

   ```
   $ mkdir ~/plugins/aciplugin
   ```

4. Log in to each cluster node and run the following commands: 

   ```
   $ sudo mkdir -p /var/log/aciplugin_logs/
   ```

   ```
   $ sudo chown odimra:odimra /var/log/aciplugin_logs
   ```

5. On the deployment node, copy the Cisco ACI plugin configuration file to `~/plugins/aciplugin`.

   ```
   $ cp ~/ODIM/odim-controller/helmcharts/aciplugin/aciplugin-config.yaml ~/plugins/aciplugin
   ```

5. Open the Cisco ACI plugin configuration YAML file.

   ```
   $ vi ~/plugins/aciplugin/aciplugin-config.yaml
   ```

   **Sample aciplugin-config.yaml file:**

   ```
   odimra:
    namespace: odim
    groupID: 2021
    haDeploymentEnabled: false
   aciplugin:
    eventListenerNodePort: 30086
    aciPluginRootServiceUUID: a127eedc-c29b-416c-8c82-413153a3c351
    username: admin
    password: sTfTyTZFvNj5zU5Tt0TfyDYU-ye3_ZqTMnMIj-
    LAeXaa8vCnBqq8Ga7zV6ZdfqQCdSAzmaO5AJxccD99UHLVlQ==
    lbHost: 10.24.1.232
    lbPort: 30086
    logPath: /var/log/aciplugin_logs
    apicHost: 10.42.0.50
    apicUserName: admin
    apicPassword: Apic123
    odimURL: https://api:45000
    odimUserName: admin
    odimPassword: Od!m12$4
   ```

8. Update the following mandatory parameters in the plugin configuration file:

   - **hostname**: Hostname of the cluster node where the Cisco ACI plugin will be installed.
   - **lbHost**: IP address of the cluster node where the Cisco ACI plugin will be installed.
   - **lbPort**: Default port is 30084.
   - **aciPluginRootServiceUUID**: RootServiceUUID to be used by the Cisco ACI plugin service.
   
   Other parameters can either be empty or have default values. Optionally, you can update them with values based on your requirements. For more information on each parameter, see [Plugin configuration parameters](#plugin-configuration-parameters).

7. Generate the Helm package for the Cisco ACI plugin on the deployment node:

   1. Navigate to `odim-controller/helmcharts/aciplugin`.

      ```
      $ cd ~/ODIM/odim-controller/helmcharts/aciplugin
      ```

   2. Run the following command to create `aciplugin` Helm package at `~/plugins/aciplugin`:

      ```
      $ helm package aciplugin -d ~/plugins/aciplugin
      ```

      The Helm package for the Cisco ACI plugin is created in the tgz format.

8. Save the Cisco ACI plugin Docker image on the deployment node at `~/plugins/aciplugin`.

   ```
   $ sudo docker save aciplugin:1.0 -o ~/plugins/aciplugin/aciplugin.tar
   ```

9. If it is a three-node cluster configuration, log in to each cluster node and [configure proxy server for the plugin](#configuring-proxy-server-for-a-plugin-version). 

   Skip this step if it is a one-node cluster configuration.

10. Navigate to the `/ODIM/odim-controller/scripts` directory on the deployment node.

    ```
    $ cd ~/ODIM/odim-controller/scripts
    ```

11. Run the following command to install the Cisco ACI plugin: 

    ```
    $ python3 odim-controller.py --config \
     /home/${USER}/ODIM/odim-controller/scripts\
    /kube_deploy_nodes.yaml --add plugin --plugin aciplugin
    ```

12. Run the following command on the cluster nodes to verify the Cisco ACI plugin pod is up and running: 

    `$ kubectl get pods -n odim`
    
    Example output showing the Cisco ACI plugin pod details:
    
    | NAME                      | READY | STATUS  | RESTARTS | AGE   |
    | ------------------------- | ----- | ------- | -------- | ----- |
    | aciplugin-5fc4b6788-2xx97 | 1/1   | Running | 0        | 4d22h |

13. [Add the Cisco ACI plugin into the Resource Aggregator for ODIM framework](#adding-a-plugin-into-the-resource-aggregator-for-odim-framework). 

## Adding a plugin into the Resource Aggregator for ODIM framework

After a plugin is successfully deployed, you must add it into the Resource Aggregator for ODIM framework to access the plugin service.

**Prerequisites**

The plugin you want to add is successfully deployed.

1. To add a plugin, perform HTTP `POST` on the following URI: 

   `https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources` 

   -   `{odim_host}` is the virtual IP address of the Kubernetes cluster.

   -   `{port}` is the API server port configured in Nginx. The default port is `30080`. If you have changed the default port, use that as the port.

   Provide a JSON request payload specifying:

   -   The plugin address \(the plugin name or hostname and the plugin port\)

   -   The username and password of the plugin user account

   -   A link to the connection method having the details of the plugin

   **Sample request payload for adding Cisco ACI Plugin:** 

   ```
   {
      "HostName":"aciplugin:45007",
      "UserName":"admin",
      "Password":"Plug!n12$4",
      "Links":{
              "ConnectionMethod": {
                "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/d172e66c-b4a8-437c-981b-1c07ddfeacaa"
            }
      }
   }
   ```

   **Request payload parameters** 

   | Parameter        | Type                    | Description                                                  |
   | ---------------- | ----------------------- | ------------------------------------------------------------ |
   | HostName         | String \(required\)<br> | It is the plugin service name and the port specified in the Kubernetes environment. For default plugin ports, see [Resource Aggregator for ODIM default ports](#resource-aggregator-for-odim-default-ports).<br><blockquote>NOTE:<br>If you are using a different port for a plugin, ensure that the port is greater than `45000`.<br></blockquote> |
   | UserName         | String \(required\)<br> | The plugin username. See default administrator account usernames of all the plugins in "Default plugin credentials".<br> |
   | Password         | String \(required\)<br> | The plugin password. See default administrator account passwords of all the plugins in "Default plugin credentials".<br> |
   | ConnectionMethod | Array \(required\)<br>  | Links to the connection methods that are used to communicate with this endpoint: `/redfish/v1/AggregationService/AggregationSources`.<br><blockquote>NOTE: Ensure that the connection method information for the plugin you want to add is updated in the odim-controller configuration file.<br></blockquote>To know which connection method to use, do the following:<br>    1.  Perform HTTP `GET` on: `/redfish/v1/AggregationService/ConnectionMethods`.<br>You will receive a list of links to available connection methods.<br>    2.  Perform HTTP `GET` on each link. Check the value of the `ConnectionMethodVariant` property in the JSON response. It displays the details of a plugin. Choose a connection method having the details of the plugin of your choice. For available connection method variants, see "Connection method variants" table.<br> |
   
   | Plugin           | Default username | Default password | Connection method variant    |
   | ---------------- | ---------------- | ---------------- | ---------------------------- |
   | Cisco ACI plugin | admin            | Plug!n12$4       | Fabric:XAuthToken:ACI_v1.0.0 |
   
   Use the following curl command to add the plugin:

   ```
curl -i POST \
      -H 'Authorization:Basic {base64_encoded_string_of_[odim_username:odim_password]}' \
      -H "Content-Type:application/json" \
      -d \
   '{"HostName":"{plugin_host}:{port}",
     "UserName":"{plugin_userName}",
     "Password":"{plugin_password}", 
    "Links":{
         "ConnectionMethod": {
            "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/{ConnectionMethodId}"
         }
   }
   }' \
 'https://{odim_host}:30080/redfish/v1/AggregationService/AggregationSources'
   ```
   
   <blockquote>
    NOTE: To generate a base64 encoded string of `{odim_username:odim_password}`, run the following command:</blockquote>
   
   ```
   $ echo -n '{odim_username}:{odim_password}' | base64 -w0
   ```
   
   Replace `{base64_encoded_string_of_[odim_username:odim_password]}` with the generated base64 encoded string in the curl command. You will receive:
   
   -   An HTTP `202 Accepted` status code.
   -   A link to the task monitor associated with this operation in the response header.
   
   To know the status of this task, perform HTTP `GET` on the `taskmon` URI until the task is complete. If the plugin is added successfully, you will receive an HTTP `200 OK` status code.
   
   After the plugin is successfully added, it will also be available as a manager resource at:
   
   `/redfish/v1/Managers`.
   
   For more information, refer to "Adding a plugin" in the [Resource Aggregator for Open Distributed Infrastructure Management™ API Reference and User Guide](https://github.com/ODIM-Project/ODIM/tree/development/docs).

2. To verify that the added plugin is active and running, do the following: 

   1. To get the list of all available managers, perform HTTP `GET` on: 

      `/redfish/v1/Managers` 

      You will receive JSON response having a collection of links to the manager resources. You will see the following links in the collection:

      -   A link to the resource aggregator manager.

      -   Links to all the added plugin managers.

   2. To identify the plugin id of the added plugin, perform HTTP `GET` on each manager link in the response. 

      The JSON response body for a plugin manager has `Name` as the plugin name.
      Example:
      The JSON response body for the URP plugin manager has `Name` as `URP`.

      **Sample response \(URP manager\)** 

      ```
      {
         "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
         "@odata.etag":"W/\"AA6D42B0\"",
         "@odata.id":"/redfish/v1/Managers/536cee48-84b2-43dd-b6e2-2459ac0eeac6",
         "@odata.type":"#Manager.v1_3_3.Manager",
         "FirmwareVersion":"1.0",
         "Id":"a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b",
         "ManagerType":"Service",
         "Name":"URP",
         "Status":{
            "Health":"OK",
            "State":"Enabled"
         },
         "UUID":"a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b"
      }
      ```

   3. Check in the JSON response of the plugin manager, if: 

      -    `State` is `Enabled` 

      -    `Health` is `Ok` 

      For more information, refer to "Managers" in [Resource Aggregator for Open Distributed Infrastructure Management™ API Reference and User Guide](https://github.com/ODIM-Project/ODIM/tree/development/docs).

## Cisco ACI fabric APIs

Resource Aggregator for ODIM exposes Redfish APIs to view and manage simple fabrics. A fabric is a network topology consisting of entities such as interconnecting switches, zones, endpoints, and address pools. The Redfish `Fabrics` APIs allow you to create and remove these entities in a fabric.

When creating fabric entities, ensure to create them in the following order:

1.  Zone-specific address pools

2.  Address pools for zone of zones

3.  Zone of zones

4.  Endpoints

5.  Zone of endpoints


When deleting fabric entities, ensure to delete them in the following order:

1.  Zone of endpoints

2.  Endpoints

3.  Zone of zones

4.  Address pools for zone of zones

5.  Zone-specific address pools


<blockquote>
IMPORTANT:


Before using the `Fabrics` APIs, ensure that the fabric manager is installed, its plugin is deployed, and added into the Resource Aggregator for ODIM framework.

| API URI                                                      | Operation Applicable | Required privileges            |
| ------------------------------------------------------------ | -------------------- | ------------------------------ |
| /redfish/v1/Fabrics/\{fabricId\}/AddressPools                | GET, POST            | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/AddressPools/\{addresspoolid\} | GET, DELETE          | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/Zones                       | GET, POST            | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/Zones/\{zoneId\}            | GET, PATCH, DELETE   | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/Endpoints                   | GET, POST            | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/Endpoints/\{endpointId\}    | GET, DELETE          | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\} /Switches/\{switchId\}/Ports/\{portid\}<br> | GET                  | `Login`                        |

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/AddressPools`                |
| **Description**    | This operation creates an address pool for a zone of zones in a specific fabric. |
| **Returns**        | - Link to the created address pool in the `Location` header.<br />- JSON schema representing the created address pool. |
| **Response code**  | On success, `201 Created`                                    |
| **Authentication** | Yes                                                          |


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
  "Name": "HPE-AddressPool-ZoneofZone",
"Description": "VLANRange for creating domain",
"Ethernet": {
"IPv4": {
"VLANIdentifierAddressRange": {
"Lower": 100,
"Upper": 200
							  }
   		}
			}
}'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools'

```

>**Sample request body**

```
{
  "Name": "HPE-AddressPool-ZoneofZone",
  "Description": "VLANRange for creating domain",
  "Ethernet": {
  "IPv4": {
  "VLANIdentifierAddressRange": {
  "Lower": 100,
  "Upper": 200
 							    }
 		  }
 		 	  }
 }
```

**Request parameters**

| Parameter                    | Type                     | Description                                                  |
| ---------------------------- | ------------------------ | ------------------------------------------------------------ |
| Name                         | String                   | Name for the address pool.                                   |
| Description                  | String \(optional\)<br>  | Description for the address pool.                            |
| IPv4\{                       | \(required\)<br>         |                                                              |
| VlanIdentifierAddressRange\{ | \(required\)<br>         | A single VLAN \(virtual LAN\) used for creating the IP interface for the user Virtual Routing and Forwarding \(VRF\).<br> |
| Lower                        | Integer \(required\)<br> | VLAN lower address                                           |
| Upper\}                      | Integer \(required\)<br> | VLAN upper address                                           |
| IbgpAddressRange\{           | \(required\)<br>         | IPv4 address used as the Router Id for the VRF per switch.<br> |
| Lower                        | String \(required\)<br>  | IPv4 lower address                                           |
| Upper\}                      | String \(required\)<br>  | IPv4 upper address                                           |
| EbgpAddressRange\{           | \(optional\)<br>         | External neighbor IPv4 addresses.                            |
| Lower                        | String \(required\)<br>  | IPv4 lower address                                           |
| Upper\} \}                   | String \(required\)<br>  | IPv4 upper address                                           |
| Ebgp\{                       | \(optional\)<br>         |                                                              |
| AsNumberRange\{              | \(optional\)<br>         | External neighbor ASN.<br>**NOTE:**<br> `EbgpAddressRange` and `AsNumberRange` values should be a matching sequence and should be of same length. |
| Lower                        | Integer \(optional\)<br> |                                                              |
| Upper\} \}                   | Integer \(optional\)<br> |                                                              |
| BgpEvpn\{                    | \(required\)<br>         |                                                              |
| RouteDistinguisherList       | Array \(required\)<br>   | Single route distinguisher value for the VRF.<br>            |
| RouteTargetList              | Array \(optional\)<br>   | Route targets. By default, the route targets will be configured as both import and export.<br> |
| GatewayIPAddressList\}       | Array \(required\)<br>   | IP pool to assign IPv4 address to the IP interface used by the VRF per switch.<br> |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:18:58 GMT
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d",
   "@odata.type":"#AddressPool.vxx.AddressPool",
   "BgpEvpn":{
      "AnycastGatewayIPAddress":"",
      "AnycastGatewayMACAddress":"",
      "GatewayIPAddressList":[
         "192.168.18.122/31",
         "192.168.18.123/31"
      ],
      "RouteDistinguisherList":[
         "65002:102"
      ],
      "RouteTargetList":[
         "65002:102",
         "65002:102"
      ]
   },
   "Description":"",
   "Ebgp":{
      "AsNumberRange":{
         "Lower":65120,
         "Upper":65125
      }
   },
   "IPv4":{
      "EbgpAddressRange":{
         "Lower":"172.12.1.10",
         "Upper":"172.12.1.15"
      },
      "FabricLinkAddressRange":{
         "Lower":"",
         "Upper":""
      },
      "IbgpAddressRange":{
         "Lower":"192.12.1.10",
         "Upper":"192.12.1.15"
      },
      "LoopbackAddressRange":{
         "Lower":"",
         "Upper":""
      },
      "NativeVlan":0,
      "VlanIdentifierAddressRange":{
         "Lower":3002,
         "Upper":3002
      }
   },
   "Id":"84766158-cbac-4f69-8ed5-fa5f2b331b9d",
   "Links":{
      "Zones":[

      ]
   },
   "Name":"AddressPool for ZoneOfZones - Vlan3002"
}
```



## Configuring proxy server for a plugin version

1. Log in to each cluster node and navigate to the following path: 

   ```
   $ cd /opt/nginx/servers
   ```

2. Create a plugin configuration file called `<plugin-name>_nginx_server.conf`: 

   ```
   $ vi <plugin-name>_nginx_server.conf
   ```

   Example:

   ```
   $ vi grfplugin_nginx_server.conf
   ```

3. Copy the following content into the `<plugin-name>_nginx_server.conf` file on each cluster node: 

   ```
   upstream <plugin_name>  {
     server <k8s_self_node_IP>:<plugin_node_port> max_fails=2 fail_timeout=10s;
     server <k8s_node2_IP>:<plugin_node_port> max_fails=2 fail_timeout=10s backup;
     server <k8s_node3_IP>:<plugin_node_port> max_fails=2 fail_timeout=10s backup;
   }
    
   server {
           listen <k8s_self_node_IP>:<nginx_plugin_port> ssl;
           listen <VIP>:<nginx_plugin_port> ssl;
           server_name odim_proxy;
           ssl_session_timeout  5m;
           ssl_prefer_server_ciphers on;
           ssl_protocols TLSv1.2;
           ssl_certificate  /opt/nginx/certs/server.crt;
           ssl_certificate_key /opt/nginx/certs/server.key;
           ssl_trusted_certificate /opt/nginx/certs/rootCA.crt;
    
           location / {
                   proxy_pass https://<plugin_name>;
                   proxy_http_version 1.1;
                   proxy_set_header X-Forwarded-For $remote_addr;
                   proxy_pass_header Server;
                   proxy_ssl_protocols TLSv1.2;
                   proxy_ssl_certificate /opt/nginx/certs/server.crt;
                   proxy_ssl_certificate_key /opt/nginx/certs/server.key;
                   proxy_ssl_trusted_certificate /opt/nginx/certs/rootCA.crt;
           }
   }
   ```

   In this content, replace the following placeholders \(highlighted in bold\) with the actual values:

   | Placeholder                      | Description                                                  |
   | -------------------------------- | ------------------------------------------------------------ |
   | <plugin_name>                    | Name of the plugin. Example: "aciplugin"<br>                 |
   | <k8s_self_node_IP>               | The physical IP address of the cluster node.                 |
   | <k8s_node2_IP><k8s_node3_IP><br> | The physical IP addresses of the other cluster nodes.        |
   | <plugin_node_port>               | The port specified for the eventListenerNodePort configuration parameter in the `<plugin_name>-config.yaml` file. |
   | <VIP>                            | Virtual IP address specified in the keepalived.conf file.    |
   | <nginx_plugin_port>              | Any free port on the cluster node. It must be available on all the other cluster nodes. Preferred port is above 45000.<br>Ensure that this port is not used as any other service port.<br>**NOTE**: You can reach the resource aggregator API server at:<br>`https://<VIP>:<nginx_api_port>`.<br> |

4. Restart Nginx systemd service only on the leader node \(cluster node where Keepalived priority is set to a higher number\): 

   ```
   $ sudo systemctl restart nginx
   ```

   <blockquote>
   NOTE:If you restart Nginx on a follower node \(cluster node having lower Keepalived priority number\), the service fails to start with the following error:</blockquote>
   
   ```
   nginx: [emerg] bind() to <VIP>:<nginx_port> failed (99: Cannot assign requested address)
   ```

## Plugin configuration parameters

The following table lists all the configuration parameters required to deploy a plugin service:

| Parameter             | Description                                                  |
| --------------------- | ------------------------------------------------------------ |
| odimra                | List of configurations required for deploying the services of Resource Aggregator for ODIM and third-party services.<br> **NOTE**: Ensure that the values of the parameters listed under odimra are same as the ones specified in the `kube_deploy_nodes.yaml` file. |
| namespace             | Namespace to be used for creating the service pods of Resource Aggregator for ODIM. The default value is "odim". You can optionally change it to a different value. |
| groupID               | Group ID to be used for creating the odimra group.The default value is 2021. You can optionally change it to a different value.<br>**NOTE**: Ensure that the group id is not already in use on any of the nodes. |
| haDeploymentEnabled   | When set to true, it deploys third-party services as a three-instance cluster. By default, it is set to true. Before setting it to false, ensure that there are at least three nodes in the Kubernetes cluster. |
| eventListenerNodePort | The port used for listening to plugin events. Refer to the sample plugin yaml configuration file to view the sample port information. |
| RootServiceUUID       | RootServiceUUID to be used by the plugin service. To generate an UUID, run the following command:<br> ```$ uuidgen```<br> Copy the output and paste it as the value for rootServiceUUID. |
| username              | Username of the plugin.                                      |
| password              | The encrypted password of the plugin.                        |
| lbHost                | If there is only one cluster node, the lbHost is the IP address of the cluster node. If there is more than one cluster node \( haDeploymentEnabled is true\), lbHost is the virtual IP address configured in Nginx and Keepalived. |
| lbPort                | If it is a one-cluster configuration, the lbPort must be same as eventListenerNodePort.<br>If there is more than one cluster node \(haDeploymentEnabled is true\), lbPort is the Nginx API node port configured in the Nginx plugin configuration file. |
| logPath               | The path where the plugin logs are stored. The default path is `/var/log/<plugin_name>_logs`<br />**Example**: `/var/log/grfplugin\_logs` |

 