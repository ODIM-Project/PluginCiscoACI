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

## 7Deploying the Cisco ACI plugin

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

5. Open the Dell plugin configuration YAML file.
5. Open the Cisco ACI plugin configuration YAML file.
   ```
   $ vi ~/plugins/aciplugin/aciplugin-config.yaml
   ```

   **Sample aciplugin-config.yaml file:**

   ```
   odimra:
    namespace: odim
    groupID: 2021
   aciplugin:
    hostname: knode1
    eventListenerNodePort: 30084
    aciPluginRootServiceUUID: 7a38b735-8b9f-48a0-b3e7-e5a180567d37
    username: admin
    password: sTfTyTZFvNj5zU5Tt0TfyDYU-ye3_ZqTMnMIj-LAeXaa8vCnBqq8Ga7zV6ZdfqQCdSAzmaO5AJxccD99UHLVlQ==
    lbHost: 10.24.1.232
    lbPort: 30084
    logPath: /var/log/aciplugin_logs
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
   
9. Generate the Helm package for the Cisco ACI plugin on the deployment node:

   1. Navigate to `odim-controller/helmcharts/aciplugin`.

      ```
      $ cd ~/ODIM/odim-controller/helmcharts/aciplugin
      ```

   2. Run the following command to create `aciplugin` Helm package at `~/plugins/aciplugin`:

      ```
      $ helm package aciplugin -d ~/plugins/aciplugin
      ```

      The Helm package for the Cisco ACI plugin is created in the tgz format.

8. Save the Dell plugin Docker image on the deployment node at `~/plugins/aciplugin`.
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
   | HostName         | String \(required\)<br> | It is the plugin service name and the port specified in the Kubernetes environment. For default plugin ports, see [Resource Aggregator for ODIM default ports](#resource-aggregator-for-odim-default-ports).<br>**NOTE**: If you are using a different port for a plugin, ensure that the port is greater than `45000`. |
   | UserName         | String \(required\)<br> | The plugin username. See default administrator account usernames of all the plugins in "Default plugin credentials". |
   | Password         | String \(required\)<br> | The plugin password. See default administrator account passwords of all the plugins in "Default plugin credentials". |
   | ConnectionMethod | Array \(required\)<br>  | Links to the connection methods that are used to communicate with this endpoint: `/redfish/v1/AggregationService/AggregationSources`.<br>**NOTE**: Ensure that the connection method information for the plugin you want to add is updated in the odim-controller configuration file.<br>To know which connection method to use, do the following:<br>    1.  Perform HTTP `GET` on: `/redfish/v1/AggregationService/ConnectionMethods`.<br>You will receive a list of links to available connection methods.<br>    2.  Perform HTTP `GET` on each link. Check the value of the `ConnectionMethodVariant` property in the JSON response. It displays the details of a plugin. Choose a connection method having the details of the plugin of your choice. For available connection method variants, see "Connection method variants" table.<br> |
   
   | Plugin           | Default username | Default password | Connection method variant    |
   | ---------------- | ---------------- | ---------------- | ---------------------------- |
   | Cisco ACI plugin | admin            | Plug!n12$4       | Fabric:XAuthToken:ACI_v1.0.0 |
   
   Use the following curl command to add the plugin:
   
   ```
   curl -i POST \
       -H 'Authorization:Basic {base64_encoded_string_of_<odim_username:odim_password>}' \
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
     'https://{odim_host}:30080/redfish/v1/AggregationService/AggregationSources' -kNOTE: To generate a base64 encoded string of `{odim_username:odim_password}`, run the following command:
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
   
      `/redfish/v1/Managers`
   
   For more information, refer to "Adding a plugin" in the [Resource Aggregator for Open Distributed Infrastructure Management™ API Reference and User Guide](https://github.com/ODIM-Project/ODIM/tree/development/docs).  
   
2. To verify that the added plugin is active and running, do the following: 

   1. To get the list of all available managers, perform HTTP `GET` on: 

      `/redfish/v1/Managers` 

      You will receive JSON response having a collection of links to the manager resources. You will see the following links in the collection:

      -   A link to the resource aggregator manager.

      -   Links to all the added plugin managers.

   2. To identify the plugin Id of the added plugin, perform HTTP `GET` on each manager link in the response. 
      The JSON response body for a plugin manager has `Name` as the plugin name.
      Example:
      The JSON response body for the URP plugin manager has `Name` as `CiscoACI`.

      **Sample response**

      ```
      {
              "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
              "@odata.etag":"W/\"AA6D42B0\"",
              "@odata.id":"/redfish/v1/Managers/536cee48-84b2-43dd-b6e2-2459ac0eeac6",
              "@odata.type":"#Manager.v1_3_3.Manager",
              "FirmwareVersion":"1.0",
              "Id":"a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b",
              "ManagerType":"Service",
              "Name":"CiscoACI",
              "Status":{
                 "Health":"OK",
                 "State":"Enabled"
              },
              "UUID":"a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b"
           }
      ```

3. Check in the JSON response of the plugin manager, if: 

      - `State` is `Enabled` 

      - `Health` is `Ok` 

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
    IMPORTANT:Before using the `Fabrics` APIs, ensure that the fabric manager is installed, its plugin is deployed, and added into the Resource Aggregator for ODIM framework. </blockquote>

| API URI                                                      | Operation Applicable | Required privileges            |
| ------------------------------------------------------------ | -------------------- | ------------------------------ |
| /redfish/v1/Fabrics/\{fabricId\}/AddressPools                | GET, POST            | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/AddressPools/\{addresspoolid\} | GET, DELETE          | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/Zones                       | GET, POST            | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/Zones/\{zoneId\}            | GET, PATCH, DELETE   | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/Endpoints                   | GET, POST            | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\}/Endpoints/\{endpointId\}    | GET, DELETE          | `Login`, `ConfigureComponents` |
| /redfish/v1/Fabrics/\{fabricId\} /Switches/\{switchId\}/Ports/\{portid\}<br> | GET                  | `Login`                        |

## Creating an addresspool for a zone of zones

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

| Parameter                   | Type                      | Description                                                  |
| --------------------------- | ------------------------- | ------------------------------------------------------------ |
| Name                        | String (optional)         | Name for the address pool.                                   |
| Description                 | String (optional)         | Description for the address pool.                            |
| Ethernet{                   |                           |                                                              |
| IPv4\{                      | \(required\)<br>          |                                                              |
| VlanIdentifierAddressRange{ | (required)                | A single VLAN (virtual LAN) used for creating the IP interface for the user Virtual Routing and Forwarding (VRF). |
| Lower                       | Integer \(required\)<br>  | VLAN lower address                                           |
| Upper\}<br />}}             | Integer \(required\)<br/> | VLAN upper address<br />                                     |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/AddressPools/e08cdf22-3c69-4548-a73b-0532111876de
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#AddressPool.AddressPool",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/AddressPools/e08cdf22-3c69-4548-a73b-0532111876de",
"@odata.type":"#AddressPool.v1_1_0.AddressPool",
"Description":"VLANRange for creating domain",
"Ethernet":{
"IPv4":{
"VLANIdentifierAddressRange":{
"Lower":100,
"Upper":200
}
}
},
"Name":"TEST-AddressPool-ZoneofZone",
"id":"e08cdf22-3c69-4548-a73b-0532111876de"
}
```

## Creating an addresspool for a zone of endpoints

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/AddressPools`                |
| **Description**    | This operation creates an address pool that can be used by a zone of endpoints. |
| **Returns**        | - Link to the created address pool in the `Location` header.<br />- JSON schema representing the created address pool. |
| **Response code**  | On success, `201 Created`                                    |
| **Authentication** | Yes                                                          |


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
  "Name":"Test-AddressPool-1",
"Ethernet":{
"IPv4":{
"GatewayIPAddress":"10.18.100.1/24",
"VLANIdentifierAddressRange":{
"Lower":100,
"Upper":100
							  }
   		}
			}
}'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools'

```

>**Sample request body**

```
{
 "Name":"Test-AddressPool-1",
 "Ethernet":{
 "IPv4":{
 "GatewayIPAddress":"10.18.100.1/24",
 "VLANIdentifierAddressRange":{
 "Lower":100,
 "Upper":100
 							    }
 		  }
 		 	  }
 }
```

**Request parameters**

| Parameter                   | Type                     | Description                                                  |
| --------------------------- | ------------------------ | ------------------------------------------------------------ |
| Name                        | String (optional)        | Name for the address pool.                                   |
| Ethernet{                   |                          |                                                              |
| IPv4\{                      | \(required\)<br>         |                                                              |
| GatewayIPAddressList\{      | Array \(required\)<br>   | IP pool to assign IPv4 address to the IP interface for VLAN per switch. |
| VlanIdentifierAddressRange{ | (required)               | A single VLAN (virtual LAN) used for creating the IP interface for the user Virtual Routing and Forwarding (VRF). |
| Lower                       | Integer \(required\)<br> | VLAN lower address                                           |
| Upper\}<br />}}             |                          | VLAN upper address<br />Lower and Upper must have the same values for the addresspool created for ZoneOfEndpoints. |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/AddressPools/bb2cd119-01e5-499d-8465-c219ad891842
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#AddressPool.AddressPool",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/AddressPools/bb2cd119-01e5-499d-8465-c219ad891842",
"@odata.type":"#AddressPool.v1_1_0.AddressPool",
"Ethernet":{
"IPv4":{
"GatewayIPAddress":"10.18.100.1/24",
"VLANIdentifierAddressRange":{
"Lower":100,
"Upper":100
}
}
},
"Name":"HPE-AddressPool-1",
"id":"bb2cd119-01e5-499d-8465-c219ad891842"
}
```

## Creating a default zone

| **Method**         | `POST`                                                      |
| ------------------ | ----------------------------------------------------------- |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/Zones`                      |
| **Description**    | This operation creates a default zone in a specific fabric. |
| **Returns**        | JSON schema representing the created zone.                  |
| **Response code**  | On success, `201 Created`                                   |
| **Authentication** | Yes                                                         |


>**curl command**


```
curl -i POST \
-H "X-Auth-Token:{X-Auth-Token}" \
-d \
'{
"Name":"HPE-Tenant",
"Description":"Default Zone",
"ZoneType":"Default"
}'
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'

```

>**Sample request body**

```
{
"Name":"HPE-Tenant",
"Description":"Default Zone",
"ZoneType":"Default"
}
```

**Request parameters**

| Parameter   | Type                  | Description                                                  |
| ----------- | --------------------- | ------------------------------------------------------------ |
| Name        | String (optional)     | Name for the zone.<br />**NOTE**: Ensure that there are no spaces. |
| Description | String (optional)     | The description for the zone.                                |
| ZoneType    | String<br/>(required) | The type of the zone to be created. Options include:<br/>• ZoneOfZones<br/>• ZoneOfEndpoints<br/>• Default<br/>The type of the zone for a default zone is Default.<br/> |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow: "GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control: no-cache
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Location: /redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/adce4bd8-0f39-421d-9b78-5fb6981ca68b
Odata-Version: 4.0
X-Frame-Options: sameorigin
Date: Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
Transfer-Encoding: chunked

```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#Zone.Zone",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/adce4bd8-0f39-421d-9b78-5fb6981ca68b",
"@odata.type":"#Zone.v1_4_0.Zone",
"Description":"Default Zone",
"Id":"adce4bd8-0f39-421d-9b78-5fb6981ca68b",
"Name":"HPE-Tenant",
"Status":{
"Health":"OK",
"State":"Enabled"
},
"ZoneType":"Default"
}
```

## Creating a zone of zones

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/Zones`                       |
| **Description**    | This operation creates an empty container zone for all the other zones in a specific fabric. You can assign address pools, endpoints, other zones, or switches to this zone. |
| **Returns**        | JSON schema representing the created zone.                   |
| **Response code**  | On success, `201 Created`                                    |
| **Authentication** | Yes                                                          |


>**curl command**


```
curl -i POST \
-H "X-Auth-Token:{X-Auth-Token}" \
-d \
'{
"Name":"HPE-App1",
"Description":"Zone of endpoints",
"ZoneType":"ZoneOfZones",
"Links":{
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Zones/e5badcc7-707c-443d-b06f-b59686e1352d"
}
],
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/AddressPools/311a78f1-ce37-4ace-9020-04688e55c398"
}
]
}
}'
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'
```

>**Sample request body**

```
{
"Name":"HPE-App1",
"Description":"Zone of endpoints",
"ZoneType":"ZoneOfZones",
"Links":{
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Zones/e5badcc7-707c-443d-b06f-b59686e1352d"
}
],
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/AddressPools/311a78f1-ce37-4ace-9020-04688e55c398"
}
]
}
}
```

**Request parameters**

| Parameter        | Type                  | Description                                                  |
| ---------------- | --------------------- | ------------------------------------------------------------ |
| Name             | String (optional)     | Name for the zone.<br />**NOTE**: Ensure that there are no spaces. |
| Description      | String (optional)     | Description for the zone.                                    |
| ZoneType         | String<br/>(required) | The type of the zone to be created. Options include:<br/>• ZoneOfZones<br/>• ZoneOfEndpoints<br/>• Default<br/>The type of the zone for a default zone is ZoneOfZones.<br/> |
| Links{           | (required)            |                                                              |
| ContainedByZones | Array<br/>(required)  | Represents an array of default zones for the zone being created. |
| AddressPools     | Array<br/>(required)  | AddressPool links supported for the Zone of Zones (AddressPool links created for ZoneOfZones). |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/6415d9aa-47a3-439d-93bb-5b23dccf5d60
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#Zone.Zone",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/
Zones/6415d9aa-47a3-439d-93bb-5b23dccf5d60",
"@odata.type":"#Zone.v1_4_0.Zone",
"Description":"Zone of endpoints",
"Id":"6415d9aa-47a3-439d-93bb-5b23dccf5d60",
"Links":{
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/AddressPools/
3d251ab9-2566-410a-9416-8164a0080d9a"
}
],
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/Zones/adce4bd8-0f39-421d-9b78-5fb6981ca68b"
}
],
"ContainedByZones@odata.count":1
},
"Name":"HPE-App1",
"Status":{
"Health":"OK",
"State":"Enabled"
},
"ZoneType":"ZoneOfZones"
}
```

## Updating the connected ports

| **Method**         | `PATCH`                                                      |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricid}/Switches/{switchid}/Ports/{portid}` |
| **Description**    | This operation updates a connected port.                     |
| **Returns**        | JSON schema representing the updated connected port.         |
| **Response code**  | On success, `200 OK`                                         |
| **Authentication** | Yes                                                          |

>**curl command**


```
curl -i PATCH \
-H "X-Auth-Token:{X-Auth-Token}" \
-d \
'{
"Links":{
"ConnectedPorts":[
{
"@odata.id":"/redfish/v1/Systems/4f65617c-7337-4bc0-a277-e9b3b4865af7:1/EthernetInterfaces/1"
}
]
}
}'
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches/{switchid}/Ports/{portid}'
```

>**Sample request body**

```
{
"Links":{
"ConnectedPorts":[
{
"@odata.id":"/redfish/v1/Systems/4f65617c-7337-4bc0-a277-e9b3b4865af7:1/EthernetInterfaces/1"
}
]
}
}
```

**Request parameters**

| Parameter      | Type                 | Description                                      |
| -------------- | -------------------- | ------------------------------------------------ |
| Links{         | (required)           |                                                  |
| ContainedPorts | Array<br/>(required) | Represents an array of links to connected ports. |
>**Sample response header** 

```
HTTP/1.1 200 OK
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Switches/
9c5735d8-5598-40a4-896f-41cbc364f2fd:101/Ports/
ccae270d-4524-44de-95ba-62a92d9476d6:eth1-2
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Fri, 02 Apr 2021 07:39:26 GMT-2d 22h
Transfer-Encoding:chunked
```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#Port.Port",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/
Switches/9c5735d8-5598-40a4-896f-41cbc364f2fd:101/Ports/
ccae270d-4524-44de-95ba-62a92d9476d6:eth1-2",
"@odata.type":"#Port.v1_3_0.Port",
"Id":"ccae270d-4524-44de-95ba-62a92d9476d6:eth1-2",
"InterfaceEnabled":false,
"LinkNetworkTechnology":"Ethernet",
"Links":{
"ConnectedPorts":[
{
"@odata.id":"/redfish/v1/Systems/
951ed562-0323-4351-9c0f-6240a25ec478:1/EthernetInterfaces/1"
}
]
},
"MaxFrameSize":9000,
"Name":"Port-eth1/2",
"PortId":"eth1/2",
"PortProtocol":"Ethernet",
"PortType":"BidirectionalPort"
}
```

## Creating a redundant endpoint

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/Endpoints`                   |
| **Description**    | This operation creates a redundant endpoint in a specific fabric. |
| **Returns**        | • Link to the created endpoint in the `Location` header.<br/>• JSON schema representing the created endpoint. |
| **Response code**  | On success, `201 Created`                                    |
| **Authentication** | Yes                                                          |

>**curl command**


```
curl -i POST \
-H "X-Auth-Token:{X-Auth-Token}" \
-d \
'{
"Name": "Redundant-Endpoint-1",
"Description": "Redundant Endpoint to provide redundancy between two
Leaf switch ports",
"Redundancy": [
{
"Mode": "Sharing",
"RedundancySet": [
{
"@odata.id": "/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Switches/71b33448-9423-4c23-b517-eb6b3ce3b751:101/Ports/
90069713-cf34-4948-a5ca-abc22a13c56b:eth1-2"
},
{
"@odata.id": "/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Switches/b1b3d0e2-3860-4caf-ade4-2db7e9f6c075:102/Ports/
1655f138-2c46-49b6-aedf-61645d73ad3f:eth1-2"
}
]
}
]
}'
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints'
```

>**Sample request body**

```
{
"Name": "Redundant-Endpoint-1",
"Description": "Redundant Endpoint to provide redundancy between two
Leaf switch ports",
"Redundancy": [
{
"Mode": "Sharing",
"RedundancySet": [
{
"@odata.id": "/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Switches/71b33448-9423-4c23-b517-eb6b3ce3b751:101/Ports/
90069713-cf34-4948-a5ca-abc22a13c56b:eth1-2"
},
{
"@odata.id": "/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Switches/b1b3d0e2-3860-4caf-ade4-2db7e9f6c075:102/Ports/
1655f138-2c46-49b6-aedf-61645d73ad3f:eth1-2"
}
]
}
]
}
```

**Request parameters**

| Parameter      | Type                  | Description                                                  |
| -------------- | --------------------- | ------------------------------------------------------------ |
| Name           | String<br/>(optional) | Name for the endpoint.                                       |
| Description    | String<br/>(optional) | Description for the endpoint.                                |
| Redundancy[    | Array                 |                                                              |
| Mode           | String                | Redundancy mode.                                             |
| RedundancySet] | Array                 | Set of redundancy ports connected to the switches.<br/>These links must be switch leaf ports URIs. |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Endpoints/1cf55323-c1be-43d6-bc51-7ea0d06190d8
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
Transfer-Encoding:chunked
```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#Endpoint.Endpoint",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/
Endpoints/1cf55323-c1be-43d6-bc51-7ea0d06190d8",
"@odata.type":"#Endpoint.v1_5_0.Endpoint",
"Description":"Redundant Endpoint to provide redundancy between two Leaf
switch ports",
"Id":"1cf55323-c1be-43d6-bc51-7ea0d06190d8",
"Name":"Redundant-Endpoint-1",
"Redundancy":[
{
"Mode":"Sharing",
"RedundancySet":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/Switches/
8bfa29b9-7fec-412d-8b29-042df4ba46f5:101/Ports/
903f2727-2bf8-49b1-8ebd-97729a8f1460:eth1-2"
},
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/Switches/e941a68e-4ffc-4d65-
b3a5-3afe84f73fd7:102/Ports/43730998-10fe-491e-94a9-f48eeaa1e202:eth1-2"
}
]
}
]
}
```

## Creating a zone of endpoints

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/zones`                       |
| **Description**    | This operation creates a zone of endpoints in a specific fabric.<br />**NOTE**: Ensure that the endpoints are created first before assigning them to the zone of endpoints. |
| **Returns**        | JSON schema representing the created zone.                   |
| **Response code**  | On success, `201 Created`                                    |
| **Authentication** | Yes                                                          |

>**curl command**


```
curl -i POST \
-H "X-Auth-Token:{X-Auth-Token}" \
-d \
'{
"Name":"HPE-ZOE-1",
"Description":"Zone of endpoints",
"ZoneType":"ZoneOfEndpoints",
"Links":{
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Zones/451a7e26-00a4-4139-87b0-49e419bfa1ee"
}
],
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/AddressPools/1b695701-6ce3-457e-a530-2bc55cac5fc7"
}
],
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Endpoints/1309dbc6-85d7-4bf8-8df7-92b9f56b0092"
}
]
}
}'
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'
```

>**Sample request body**

```
{
"Name":"HPE-ZOE-1",
"Description":"Zone of endpoints",
"ZoneType":"ZoneOfEndpoints",
"Links":{
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Zones/451a7e26-00a4-4139-87b0-49e419bfa1ee"
}
],
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/AddressPools/1b695701-6ce3-457e-a530-2bc55cac5fc7"
}
],
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Endpoints/1309dbc6-85d7-4bf8-8df7-92b9f56b0092"
}
]
}
}
```

**Request parameters**

| Parameter           | Type                       | Description                                                  |
| ------------------- | -------------------------- | ------------------------------------------------------------ |
| Name                | String<br/>(optional)      | Name for the zone.<br />**NOTE**: Ensure that there are no spaces. |
| Description         | String<br/>(optional)      | Description for the zone.                                    |
| ZoneType            | String<br/>(required)<br/> | The type of the zone to be created. Options include:<br/>• ZoneOfZones<br/>• ZoneOfEndpoints<br/>• Default<br/>The type of the zone for a zone of endpoints is<br/>ZoneOfEndpoints. |
| Links{              | Object<br/>(required)      | Contains references to other resources that are related to the zone. |
| ContainedByZones [{ | Array<br/>(required)       | Represents an array of ZoneOfZones for the zone being created. |
| @odata.id }]        | String                     | Link to a Zone of zones.                                     |
| AddressPools [{     | Array<br/>(required)       | Represents an array of address pools linked with a ZoneOfZones. |
| @odata.id }]        | String                     | Link to an address pool.                                     |
| Endpoints [{        | Array<br/>(required)       | Represents an array of endpoints to be included in the zone. |
| @odata.id }]        | String                     | Link to an endpoint.                                         |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow: "GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control: no-cache
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Location: /redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/8e18e640-a91b-4d9b-9810-b63af3d9ce9b
Odata-Version: 4.0
X-Frame-Options: sameorigin
Date: Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
Transfer-Encoding: chunked
```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#Zone.Zone",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/
Zones/8e18e640-a91b-4d9b-9810-b63af3d9ce9b",
"@odata.type":"#Zone.v1_4_0.Zone",
"Description":"Zone of endpoints",
"Id":"8e18e640-a91b-4d9b-9810-b63af3d9ce9b",
"Links":{
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/AddressPools/bb2cd119-01e5-499d-8465-
c219ad891842"
}
],
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/Zones/6415d9aa-47a3-439d-93bb-5b23dccf5d60"
}
],
"ContainedByZones@odata.count":1,
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/Endpoints/1cf55323-c1be-43d6-bc51-7ea0d06190d8"
}
]
},
"Name":"HPE-ZOE-1",
"Status":{
"Health":"OK",
"State":"Enabled"
},
"ZoneType":"ZoneOfEndpoints"
}
```

## Updating a zone of endpoints

| **Method**         | `PATCH`                                                 |
| ------------------ | ------------------------------------------------------- |
| **URI**            | `/redfish/v1/Fabrics/{fabricid}/Zones/{zoneid}`         |
| **Description**    | This operation updates a zone of endpoints.             |
| **Returns**        | JSON schema representing the updated zone of endpoints. |
| **Response code**  | On success, `200 OK`                                    |
| **Authentication** | Yes                                                     |

>**curl command**


```
curl -i PATCH \
-H "X-Auth-Token:{X-Auth-Token}" \
-d \
'{
"Links":{
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Endpoints/1309dbc6-85d7-4bf8-8df7-92b9f56b0092"
}
]
}
}'
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones/{zoneid}
```

>**Sample request body**

```
{
"Links":{
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Endpoints/1309dbc6-85d7-4bf8-8df7-92b9f56b0092"
}
]
}
}
```

**Request parameters**

| Parameter | Type                 | Description                                                  |
| --------- | -------------------- | ------------------------------------------------------------ |
| Links{    | (required)           |                                                              |
| Endpoints | Array<br/>(required) | Represents an array of links to endpoints.<br/>**NOTE**: Adding new endpoint links replaces the existing ones in the zone of endpoints being updated.To retain the existing links, add them in this array along with the new ones. |


>**Sample response header** 

```
HTTP/1.1 200 OK
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/8e18e640-a91b-4d9b-9810-b63af3d9ce9b
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Fri, 02 Apr 2021 07:39:26 GMT-2d 22h
Transfer-Encoding:chunked
```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#Zone.Zone",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/
Zones/8e18e640-a91b-4d9b-9810-b63af3d9ce9b",
"@odata.type":"#Zone.v1_4_0.Zone",
"Description":"Zone of endpoints",
"Id":"8e18e640-a91b-4d9b-9810-b63af3d9ce9b",
"Links":{
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/AddressPools/bb2cd119-01e5-499d-8465-
c219ad891842"
}
],
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351:1/Zones/6415d9aa-47a3-439d-93bb-5b23dccf5d60"
}
],
"ContainedByZones@odata.count":1,
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b:
1/Endpoints/1309dbc6-85d7-4bf8-8df7-92b9f56b0092"
}
]
},
"Name":"HPE-ZOE-1",
"Status":{
"Health":"OK",
"State":"Enabled"
},
"ZoneType":"ZoneOfEndpoints"
}
```

## Deleting an ACI fabric entity

## 

| **Method**         | `DELETE`                                                     |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}`<br/>`/redfish/v1/Fabrics/{fabricID}/Zones/{zoneid}`<br/>`/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointid}` |
| **Description**    | This operation deletes a fabric entity such as an addresspool, zone, or an endpoint in a<br/>specific fabric.<br/>When deleting fabric entities, ensure to delete them in the following order:<br/>1. Zone of endpoints<br/>2. Endpoints<br/>3. Zone of zones<br/>4. Default zones<br/>5. Address pools |
| **Response code**  | On success, `204 No Content`                                 |
| **Authentication** | Yes                                                          |

>**curl command**


```
curl -i -X DELETE \
-H "X-Auth-Token:{X-Auth-Token}" \
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}'
curl -i -X DELETE \
-H "X-Auth-Token:{X-Auth-Token}" \
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones/{zoneid}'
curl -i -X DELETE \
-H "X-Auth-Token:{X-Auth-Token}" \
'https://{odim_host}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointid}'
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

   | Placeholder   | Description                                  |
   | ------------- | -------------------------------------------- |
| <plugin_name> | Name of the plugin. Example: "aciplugin"<br> |
   
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
| groupID               | Group ID to be used for creating the odimra group. The default value is 2021. You can optionally change it to a different value.<br>**NOTE**: Ensure that the group id is not already in use on any of the nodes. |
| haDeploymentEnabled   | When set to true, it deploys third-party services as a three-instance cluster. By default, it is set to true. Before setting it to false, ensure that there are at least three nodes in the Kubernetes cluster. |
| eventListenerNodePort | The port used for listening to plugin events. Refer to the sample plugin yaml configuration file to view the sample port information. |
| RootServiceUUID       | RootServiceUUID to be used by the plugin service. To generate an UUID, run the following command:<br> ```$ uuidgen```<br> Copy the output and paste it as the value for rootServiceUUID. |
| username              | Username of the plugin.                                      |
| password              | The encrypted password of the plugin.                        |
| lbHost                | If there is only one cluster node, the lbHost is the IP address of the cluster node. If there is more than one cluster node \(haDeploymentEnabled is true\), lbHost is the virtual IP address configured in Nginx and Keepalived. |
| lbPort                | If it is a one-cluster configuration, the lbPort must be same as eventListenerNodePort. <br>If there is more than one cluster node \(haDeploymentEnabled is true\), lbPort is the Nginx API node port configured in the Nginx plugin configuration file. |
| logPath               | The path where the plugin logs are stored. The default path is `/var/log/<plugin_name>_logs`<br/>**Example**: `/var/log/aciplugin\_logs`<br/> |
