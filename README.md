[![build_deploy_test Actions Status](https://github.com/ODIM-Project/ODIM/workflows/build_deploy_test/badge.svg)](https://github.com/ODIM-Project/ODIM/actions)
[![build_unittest Actions Status](https://github.com/ODIM-Project/ODIM/workflows/build_unittest/badge.svg)](https://github.com/ODIM-Project/ODIM/actions)

# Table of contents

- [Overview of Cisco ACI](#overview-of-cisco-aci)
  - [Filing defects in ACI plugin](#filing-defects-in-aci-plugin)
  - [Deploying the Cisco ACI plugin](#deploying-the-cisco-aci-plugin)
  - [Adding a plugin into the Resource Aggregator for ODIM framework](#adding-a-plugin-into-the-resource-aggregator-for-odim-framework)
  - [Plugin configuration parameters](#plugin-configuration-parameters)
  - [Resource Aggregator for ODIM default ports](#resource-aggregator-for-odim-default-ports)
- [Cisco ACI fabric APIs](#Cisco-ACI-fabric-APIs)
  - [Creating an addresspool for a zone of zones](#creating-an-addresspool-for-a-zone-of-zones)
  - [Creating an addresspool for a zone of endpoints (with tagged VLAN)](#creating-an-addresspool-for-a-zone-of-endpoints-with-tagged-vlan)
  - [Creating an addresspool for a zone of endpoints (with untagged VLAN)](#creating-an-addresspool-for-a-zone-of-endpoints-with-untagged-vlan)
  - [Creating a default zone](#creating-a-default-zone)
  - [Creating a zone of zones](#creating-a-zone-of-zones)
  - [Updating the connected ports](#updating-the-connected-ports)
  - [Creating a redundant endpoint](#creating-a-redundant-endpoint)
  - [Creating a zone of endpoints](#creating-a-zone-of-endpoints)
  - [Updating a zone of endpoints](#updating-a-zone-of-endpoints)
  - [Deleting an ACI fabric entity](#deleting-an-ACI-fabric-entity)
- [Mapping of Redfish logical entities to Cisco ACI entities](#Mapping-of-Redfish-logical-entities-to-Cisco-ACI-entities)



# Overview of Cisco ACI

Cisco ACI (Application Centric Infrastructure) is an open ecosystem model that uses a holistic, systems-based approach to integrate hardware and software, and the physical and virtual elements, to enable unique business value for modern data centers.

Resource Aggregator for ODIM supports Cisco ACI plugin that can abstract, translate, and expose southbound resource information to the resource aggregator through RESTful APIs.

## Filing defects in ACI plugin

**Important**: In case of any unforeseen issues you experience while deploying or using ACI plugin, log on to the following website and file your defect by clicking **Create**.

**Prerequisite**: You must have valid LFN Jira credentials to create defects.

- Website: https://jira.lfnetworking.org/secure/Dashboard.jspa
- Discussion Forums: https://odim.slack.com/archives/C01DG9MH479
- Documentation:
  - Cisco ACI Plugin Deployment Document - https://github.com/ODIM-Project/PluginCiscoACI/blob/main/README.md
  - Resource Aggregator for ODIM Deployment Document- https://github.com/ODIM-Project/ODIM#readme
  - Additional documents - https://github.com/ODIM-Project/ODIM/blob/main/docs

## Deploying the Cisco ACI plugin

**Prerequisites**

Kubernetes cluster is set up and the resource aggregator is successfully deployed.

1. Create a directory called `plugins` on the deployment node.

   ```
   mkdir -p ~/plugins
   ```

3. In the `plugins` directory, create a directory called `aciplugin`.

   ```
   mkdir ~/plugins/aciplugin
   ```

3. Run the following commands on the deployment node:

   1. ```
      git clone https://github.com/ODIM-Project/PluginCiscoACI.git
      ```

   2. ```
      cd PluginCiscoACI/
      ```

   3. ```
      export ODIMRA_USER_ID=2021
      ```
      
   4. ```
      export ODIMRA_GROUP_ID=2021
      ```

   5. ```
      ./build_images.sh
      ```

4. On the deployment node, copy the Cisco ACI plugin configuration file and the hook script to `~/plugins/aciplugin`.

   ```
   cp ~/PluginCiscoACI/install/Kubernetes/helmcharts/aciplugin-config.yaml ~/plugins/aciplugin
   ```

   ```
   cp ~/PluginCiscoACI/install/Kubernetes/helmcharts/aciplugin.sh ~/plugins/aciplugin
   ```

5. Open the Cisco ACI plugin configuration YAML file.

   ```
   vi ~/plugins/aciplugin/aciplugin-config.yaml
   ```

   **Sample aciplugin-config.yaml file:**

   ```
   aciplugin:
    eventListenerNodePort: 30086
    aciPluginRootServiceUUID: a127eedc-c29b-416c-8c82-413153a3c351
    username: admin
    password: sTfTyTZFvNj5zU5Tt0TfyDYU-ye3_ZqTMnMIj-
   LAeXaa8vCnBqq8Ga7zV6ZdfqQCdSAzmaO5AJxccD99UHLVlQ==
    lbHost: aciplugin
    lbPort: 30086
    logPath: /var/log/aciplugin_logs
    apicHost: 17.5.7.8
    apicUserName: admin
    apicPassword: Apic123
    odimURL: https://api:45000
    odimUserName: admin
    odimPassword: H/r7PSBpgBafwA2UPPm6CrkGTBT9H0VJX0EQKz61ktpCJbdUXUiJdoX1LoU2JMxPEQPPv2tU4z1BO0HtiELe8muJ7VilCmW51zMWv0D7O+qjV4IxhvZ5EZT4tHqfjJwBSBBZZ5cV11ceic5p8L26soCT8KMNTRhksYVQJXUJnyT6qiNuTrLAIouJ4kj4xIdelpP4Zgzy7fdyd+x+yQP2DPWgCF5fYErmk7H7gxVibay1YUaE6qVAbYypqwRmUHIjnv3VC3qTRyRfwGMWEm+xD5ySNUOocXqUORuFcSPWDpZYXWRKSYnwKA+XZuCdm6KUiqU84Hyq4O5hWLwz51XZ/SutnOIoZoooqKxhMwmqLvAsx8/ndG9m2j+M/Vx+Cm22OWweGMKvXP5xKqR5X2bMybvLbKb+mJLW8WxjM+EI+Y4XpgunRlsaExYRW/4GCg7vWcvQ8Sc5a74n20+sNKqjqs/SgLdmJTzfh/6MN0TSfn8DtALJiN/17KAyTjH/2YO/arguin/eYiMfO9X6avgjy7x2ceOzUJFaWkWEOYMV8Ksm4msvlfhHcZ+2NgIsJRgfZgbO49+K+0jwQ7p7fXv5GOcFJ6HMVPNTJ8kCayU0Yh50bsqv3e7KTIERT1XyI6zXa7LYk5sswOvl7gsndE3vkPddrHg+m194tFo92chsnv0=
   ```

6. Update the following parameters in the plugin configuration file:

   - **eventListenerNodePort**: The port used for listening to the ACI plugin events. Default port is 30086.

   - **aciPluginRootServiceUUID**: The RootServiceUUID to be used by the ACI plugin service. To generate an UUID, run the following command:

     `uuidgen`

     Copy the output and paste it as the value for rootServiceUUID.

   - **lbHost**: Default value is `aciplugin` for one node and three node cluster configuration.

   - **lbPort**: Default port is 30086 for one node and three node cluster configuration.

     <blockquote>Note: The lbport is used as proxy port for eventlistenernodeport, which is used for subscribing to events.</blockquote>

   - **apicHost**: The IP address of the machine where Cisco APIC UI is launched.

   - **apicUserName**: The Cisco APIC username.

   - **apicPassword**: The Cisco APIC password.

   - **odimURL**: The URL of the ODIMRA API service. URL is https://api:45000.

   - **odimUserName**: The username of the default administrator account of Resource Aggregator for ODIM.

   - **odimPassword**: The encrypted password of the default administrator account of Resource Aggregator for ODIM. To generate the encrypted password, run the following command:
     
     ```
     echo -n '< ODIMRA password>' |openssl pkeyutl -encrypt -inkey <odimCertsPath>/odimra_rsa.private -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512|openssl base64 -A
     ```

   Other parameters can have default values. Optionally, you can update them with values based on your requirements. For more information on each parameter, see [Plugin configuration parameters](#plugin-configuration-parameters).

7. Generate the Helm package for the Cisco ACI plugin on the deployment node:

   1. Navigate to `PluginCiscoACI/install/Kubernetes/helmcharts`.

      ```
      cd ~/PluginCiscoACI/install/Kubernetes/helmcharts
      ```

   2. Run the following command to create `aciplugin` Helm package at `~/plugins/aciplugin`:

      ```
      helm package aciplugin -d ~/plugins/aciplugin
      ```

      The Helm package for the Cisco ACI plugin is created in the tgz format.

8. Save the Cisco ACI plugin Docker image on the deployment node at `~/plugins/aciplugin`.

   ```
   docker save aciplugin:3.1 -o ~/plugins/aciplugin/aciplugin.tar
   ```

9. Navigate to the `PluginCiscoACI` directory.

   ```
   cd ~/PluginCiscoACI
   ```

10. Copy the proxy configuration file `install/templates/aciplugin_proxy_server.conf.j2` to `~/plugins/aciplugin`.

    ```
    cp install/templates/aciplugin_proxy_server.conf.j2 ~/plugins/aciplugin
    ```

    **Important**: Do NOT change the value of any parameter in this file. 

11. Navigate to the `/ODIM/odim-controller/scripts` directory on the deployment node.

    ```
    cd ~/ODIM/odim-controller/scripts
    ```

12. Open the `kube_deploy_nodes.yaml` file.

        vi kube_deploy_nodes.yaml

13. Specify values for the following parameters in the `kube_deploy_nodes.yaml` file: 

    | Parameter                    | Value                                                        |
    | ---------------------------- | ------------------------------------------------------------ |
    | connectionMethodConf         | The connection method associated with Cisco ACI plugin: ConnectionMethodVariant: <br />`Fabric:BasicAuth:ACI_v1.0.0`<br> |
    | odimraKafkaClientCertFQDNSan | The FQDN to be included in the Kafka client certificate of Resource Aggregator for ODIM for deploying the ACI plugin:<br />`aciplugin`, `aciplugin-events`<br>Add these values to the existing comma-separated list.<br> |
    | odimraServerCertFQDNSan      | The FQDN to be included in the server certificate of Resource Aggregator for ODIM for deploying the ACI plugin:<br /> `aciplugin`, `aciplugin-events`<br> Add these values to the existing comma-separated list.<br> |

        odimPluginPath: /home/bruce/plugins
        odimra:
          groupID: 2021
          userID: 2021
          namespace: odim
          fqdn:
          rootServiceUUID:
          haDeploymentEnabled: True
          connectionMethodConf:
          - ConnectionMethodType: Redfish
            ConnectionMethodVariant: Fabric:BasicAuth:ACI_v1.0.0
          odimraKafkaClientCertFQDNSan: aciplugin,aciplugin-events
          odimraServerCertFQDNSan: aciplugin,aciplugin-events

14. Move `odimra_kafka_client.key`, `odimra_kafka_client.crt`, `odimra_server.key`, and `odimra_server.crt` stored in odimCertsPath to a different folder.

    <blockquote> NOTE: odimCertsPath is the absolute path of the directory where certificates required by the services of Resource Aggregator for ODIM are present. This parameter is configured in the `kube_deploy_nodes.yaml` file.</blockquote>

15. Update odimra-secrets:

       ```
    python3 odim-controller.py --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml --upgrade odimra-secret
       ```

16. Run the following command: 

        python3 odim-controller.py --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml --upgrade odimra-config

17. In `~/ODIM/odim-controller/scripts`, run the following command to install the Cisco ACI plugin: 

    ```
    python3 odim-controller.py --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml --add plugin --plugin aciplugin
    ```

18. Run the following command on the cluster nodes to verify the Cisco ACI plugin pod is up and running: 

    ```
    kubectl get pods -n odim
    ```

    Example output showing the Cisco ACI plugin pod details:

    | NAME                      | READY | STATUS  | RESTARTS | AGE   |
    | ------------------------- | ----- | ------- | -------- | ----- |
    | aciplugin-5fc4b6788-2xx97 | 1/1   | Running | 0        | 4d22h |

19. [Add the Cisco ACI plugin into the Resource Aggregator for ODIM framework](#adding-a-plugin-into-the-resource-aggregator-for-odim-framework). 

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
         "HostName":"aciplugin:45020",
         "UserName":"admin",
         "Password":"Plug!n12$4",
         "Links":{
                 "ConnectionMethod": {
                   "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/{ConnectionMethodId}"
               }
         }
      }
   ```
   
   **Request payload parameters**
   
   | Parameter        | Type                    | Description                                                  |
   | ---------------- | ----------------------- | ------------------------------------------------------------ |
   | HostName         | String \(required\)<br> | It is the plugin service name and the port specified in the Kubernetes environment. For default plugin ports, see *[Resource Aggregator for ODIM default ports](#resource-aggregator-for-odim-default-ports)*.<br>**NOTE**: If you are using a different port for a plugin, ensure that the port is greater than `45000`. |
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
     'https://{odim_host}:30080/redfish/v1/AggregationService/AggregationSources' -k
   ```
   
   <blockquote>
       NOTE: To generate a base64 encoded string of `{odim_username:odim_password}`, run the following command:</blockquote>
   
   ```
   echo -n '{odim_username}:{odim_password}' | base64 -w0
   ```
   
   Default username is `admin` and default password is `Od!m12$4`.
   Replace `{base64_encoded_string_of_[odim_username:odim_password]}` with the generated base64 encoded string in the curl command. You will receive:
   
      -   An HTTP `202 Accepted` status code.
      -   A link of the executed task. Performing a `GET` operation on this link displays the task monitor associated with this operation in the response header.
   
   To know the status of this task, perform HTTP `GET` on the `taskmon` URI until the task is complete. If the plugin is added successfully, you will receive an HTTP `200 OK` status code.
   
   After the plugin is successfully added, it will also be available as a manager resource at:
   
      `/redfish/v1/Managers`
   
   For more information, see *Adding a plugin* in the *[Resource Aggregator for Open Distributed Infrastructure Management™ API Reference and User Guide](https://github.com/ODIM-Project/ODIM/tree/development/docs)*. 
   
2. To verify that the added plugin is active and running, do the following: 

   1. To get the list of all available managers, perform HTTP `GET` on: 

      `/redfish/v1/Managers` 

      You will receive JSON response having a collection of links to the manager resources. You will see the following links in the collection:

      -   A link to the resource aggregator manager.

      -   Links to all the added plugin managers.

   2. To identify the plugin Id of the added plugin, perform HTTP `GET` on each manager link in the response. 
      The JSON response body for a plugin manager has `Name` as the plugin name.
      **Example**:
      The JSON response body for the Cisco ACI plugin manager has `Name` as `CiscoACI`.

      **Sample response**

      ```
      {
          "@odata.context": "/redfish/v1/$metadata#Manager.Manager",
          "@odata.id": "/redfish/v1/Managers/fb40f2dc-0c6d-4464-bc98-fea775adbbb9",
          "@odata.type": "#Manager.v1_10_0.Manager",
          "FirmwareVersion": "v1.0.0",
          "Id": "fb40f2dc-0c6d-4464-bc98-fea775adbbb9",
          "Links": {
              "ManagerForSwitches": [
                  {
                      "@odata.id": "/redfish/v1/Fabrics/fb40f2dc-0c6d-4464-bc98-fea775adbbb9.1/Switches/af10c386-68d5-45aa-b3c3-431e3e4c3647.101"
                  },
                  {
                      "@odata.id": "/redfish/v1/Fabrics/fb40f2dc-0c6d-4464-bc98-fea775adbbb9.1/Switches/668f20cf-b6e7-4ded-a180-bf8e33dc18fc.102"
                  },
                  {
                      "@odata.id": "/redfish/v1/Fabrics/fb40f2dc-0c6d-4464-bc98-fea775adbbb9.1/Switches/7d5a25b3-3ac4-49f4-a929-243b5b97bba0.201"
                  }
              ],
              "ManagerForSwitches@odata.count": 3
          },
          "ManagerType": "Service",
          "Name": "ACI",
          "Status": {
              "Health": "OK",
              "State": "Enabled"
          },
          "UUID": "fb40f2dc-0c6d-4464-bc98-fea775adbbb9"
      }
      ```

3. Check in the JSON response of the plugin manager, if: 

      - `State` is `Enabled` 

      - `Health` is `Ok` 

      For more information, refer to "Managers" in *[Resource Aggregator for Open Distributed Infrastructure Management™ API Reference and User Guide](https://github.com/ODIM-Project/ODIM/tree/development/docs)*.



## Plugin configuration parameters

The following table lists all the configuration parameters required to deploy a plugin service:

| Parameter           | Description                                                  |
| ------------------- | ------------------------------------------------------------ |
| odimra              | List of configurations required for deploying the services of Resource Aggregator for ODIM and third-party services.<br> **NOTE**: Ensure the values of the parameters listed under odimra are the same as the ones specified in the `kube_deploy_nodes.yaml` file. |
| namespace           | Namespace to be used for creating the service pods of Resource Aggregator for ODIM. Default value is "odim". You can optionally change it to a different value. |
| groupID             | Group ID to be used for creating the odimra group. Default value is 2021. You can optionally change it to a different value.<br>**NOTE**: Ensure that the group id is not already in use on any of the nodes. |
| haDeploymentEnabled | When set to true, it deploys third-party services as a three-instance cluster. By default, it is set to true. Before setting it to false, ensure there are at least three nodes in the Kubernetes cluster. |
| username            | Username of the plugin                                       |
| password            | The encrypted password of the plugin                         |
| logPath             | The path where the plugin logs are stored. Default path is `/var/log/<plugin_name>_logs`<br/>**Example**: `/var/log/aciplugin_logs`<br/> |
| odimPassword        | The encrypted password of the default administrator account of Resource Aggregator for ODIM. To generate the encrypted password, run the following command:<br />`echo -n '<HPE ODIMRA password>' |openssl pkeyutl -encrypt -inkey <odimCertsPath>/odimra_rsa.private -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512|openssl base64 -A` |

## Resource Aggregator for ODIM default ports

The following table lists all the default ports used by the resource aggregator, plugins, and third-party services. The following ports (except container ports) must be free:

| Port name                                                    | Ports                                                        |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| Container ports (access restricted only to the Kubernetes cluster network) | 45000 — API service port<br />45101- 45201 — Resource Aggregator for ODIM service ports<br />9082, 9092 — Kafka ports<br />6379 — Redis port<br />26379 — Redis Sentinel port<br />2181 — Zookeeper port<br>2379, 2380 — etcd ports |
| API node port (for external access)                          | 30080                                                        |
| Kafka node port (for external access)                        | 30092 for a one-node cluster configuration<br />30092, 30093, and 30094 for a three-node cluster configuration |
| ACI port<br />EventListenerNodePort<br />lbport              | 45020 — Port to be used while adding Cisco ACI plugin<br />30083 — Port used for event subscriptions in one-node cluster configuration <br />lbport — For three-node cluster configuration, specify lbport as per your requirement. This port must be assigned with a free port (preferably above 45000) available on all cluster nodes. This port is used as nginx proxy port for the plugin<br />For one-node cluster configuration, it is the same as EventListenerNodePort |

# Cisco ACI fabric APIs

Resource Aggregator for ODIM exposes Redfish APIs to view and manage simple fabrics. A fabric is a network topology consisting of entities such as interconnecting switches, zones, endpoints, and address pools. The Redfish `Fabrics` APIs allow you to create and remove these entities in a fabric.

When creating fabric entities, ensure to create them in the following order:

1.  Address pools
2.  Default zones
3.  Zone of zones
4.  Endpoints
5.  Zone of endpoints (with tagged and untagged VLAN)

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
| /redfish/v1/Fabrics/\{fabricId\}/Switches/\{switchId\}/Ports/\{portid\}<br> | GET                  | `Login`                        |

## Creating an addresspool for a zone of zones

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/AddressPools`                |
| **Description**    | This operation creates an address pool for a zone of zones in a specific fabric. |
| **Returns**        | - Link to the created address pool in the `Location` header<br />- JSON schema representing the created address pool |
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
| Name                        | String (optional)         | Name for the address pool                                    |
| Description                 | String (optional)         | Description for the address pool                             |
| Ethernet{                   |                           |                                                              |
| IPv4\{                      | \(required\)<br>          |                                                              |
| VlanIdentifierAddressRange{ | (required)                | A single VLAN (virtual LAN) used for creating the IP interface for the user Virtual Routing and Forwarding (VRF) |
| Lower                       | Integer \(required\)<br>  | VLAN lower address                                           |
| Upper\}<br />}}             | Integer \(required\)<br/> | VLAN upper address<br />                                     |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/AddressPools/e08cdf22-3c69-4548-a73b-0532111876de
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
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

## Creating an addresspool for a zone of endpoints with Tagged VLAN

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/AddressPools`                |
| **Description**    | This operation creates an address pool (with tagged VLAN) that can be used by a zone of endpoints. |
| **Returns**        | - Link to the created address pool in the `Location` header<br />- JSON schema representing the created address pool |
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
"GatewayIPAddress":"17.5.7.8/24",
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
 "GatewayIPAddress":"17.5.7.8/24",
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
| Name                        | String (optional)        | Name for the address pool                                    |
| Ethernet{                   |                          |                                                              |
| IPv4\{                      | \(required\)<br>         |                                                              |
| GatewayIPAddressList\{      | Array \(required\)<br>   | IP pool to assign IPv4 address to the IP interface for VLAN per switch |
| VlanIdentifierAddressRange{ | (required)               | A single VLAN (virtual LAN) used for creating the IP interface for the user Virtual Routing and Forwarding (VRF) |
| Lower                       | Integer \(required\)<br> | VLAN lower address                                           |
| Upper\}<br />}}             |                          | VLAN upper address.<br />Lower and Upper must have the same values for the addresspool created for ZoneOfEndpoints. |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/AddressPools/bb2cd119-01e5-499d-8465-c219ad891842
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#AddressPool.AddressPool",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/AddressPools/bb2cd119-01e5-499d-8465-c219ad891842",
"@odata.type":"#AddressPool.v1_1_0.AddressPool",
"Ethernet":{
"IPv4":{
"GatewayIPAddress":"17.5.7.8/24",
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

## Creating an addresspool for a zone of endpoints with untagged VLAN

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricID}/AddressPools`                |
| **Description**    | This operation creates an address pool (with untagged VLAN) that can be used by a zone of endpoints. |
| **Returns**        | - Link to the created address pool in the `Location` header<br />- JSON schema representing the created address pool |
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
"GatewayIPAddress":"17.5.7.8/24",
"NativeVLAN":101
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
 "GatewayIPAddress":"17.5.7.8/24",
 "NativeVLAN":101
 		  }
 		 	  }
 }
```

**Request parameters**

| Parameter              | Type                   | Description                                                  |
| ---------------------- | ---------------------- | ------------------------------------------------------------ |
| Name                   | String (optional)      | Name for the address pool                                    |
| Ethernet{              |                        |                                                              |
| IPv4\{                 | \(required\)<br>       |                                                              |
| GatewayIPAddressList\{ | Array \(required\)<br> | IP pool to assign IPv4 address to the IP interface for VLAN per switch |
| NativeVLAN             | (required)             | A single VLAN (virtual LAN) used for creating the IP interface for the user Virtual Routing and Forwarding (VRF) |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/AddressPools/bb2cd119-01e5-499d-8465-c219ad891842
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#AddressPool.AddressPool",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351.1/AddressPools/bb2cd119-01e5-499d-8465-c219ad891842",
"@odata.type":"#AddressPool.v1_1_0.AddressPool",
"Ethernet":{
"IPv4":{
"GatewayIPAddress":"17.5.7.8/24",
"NativeVLAN":101
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
| **Returns**        | JSON schema representing the created zone                   |
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
| Description | String (optional)     | The description for the zone.<br />**NOTE**: Ensure that there are no spaces. |
| ZoneType    | String<br/>(required) | The type of the zone to be created. Options include:<br/>• ZoneOfZones<br/>• ZoneOfEndpoints<br/>• Default<br/>The type of the zone for a default zone is Default.<br/> |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow: "GET", "PUT", "POST", "PATCH", "DELETE"
Location: /redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/adce4bd8-0f39-421d-9b78-5fb6981ca68b
Date: Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
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
| **Returns**        | JSON schema representing the created zone                    |
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
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/Zones/e5badcc7-707c-443d-b06f-b59686e1352d"
}
],
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
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
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/Zones/e5badcc7-707c-443d-b06f-b59686e1352d"
}
],
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
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
| Description      | String (optional)     | Description for the zone.<br />**NOTE**: Ensure that there are no spaces. |
| ZoneType         | String<br/>(required) | The type of the zone to be created. Options include:<br/>• ZoneOfZones<br/>• ZoneOfEndpoints<br/>• Default<br/>The type of the zone for a default zone is ZoneOfZones.<br/> |
| Links{           | (required)            |                                                              |
| ContainedByZones | Array<br/>(required)  | Represents an array of default zones for the zone being created |
| AddressPools     | Array<br/>(required)  | AddressPool links supported for the Zone of Zones (AddressPool links created for ZoneOfZones) |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/6415d9aa-47a3-439d-93bb-5b23dccf5d60
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
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
416c-8c82-413153a3c351.1/AddressPools/
3d251ab9-2566-410a-9416-8164a0080d9a"
}
],
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351.1/Zones/adce4bd8-0f39-421d-9b78-5fb6981ca68b"
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
| **Returns**        | JSON schema representing the updated connected port          |
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

| Parameter      | Type                 | Description                                     |
| -------------- | -------------------- | ----------------------------------------------- |
| Links{         | (required)           |                                                 |
| ContainedPorts | Array<br/>(required) | Represents an array of links to connected ports |
>**Sample response header** 

```
HTTP/1.1 200 OK
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Switches/
9c5735d8-5598-40a4-896f-41cbc364f2fd:101/Ports/
ccae270d-4524-44de-95ba-62a92d9476d6:eth1-2
Date:Fri, 02 Apr 2021 07:39:26 GMT-2d 22h
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
951ed562-0323-4351-9c0f-6240a25ec478.1/EthernetInterfaces/1"
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
| **Returns**        | • Link to the created endpoint in the `Location` header<br/>• JSON schema representing the created endpoint |
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
"@odata.id": "/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/Switches/71b33448-9423-4c23-b517-eb6b3ce3b751.101/Ports/
90069713-cf34-4948-a5ca-abc22a13c56b.eth1-2"
},
{
"@odata.id": "/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/Switches/b1b3d0e2-3860-4caf-ade4-2db7e9f6c075.102/Ports/
1655f138-2c46-49b6-aedf-61645d73ad3f.eth1-2"
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
"@odata.id": "/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/Switches/71b33448-9423-4c23-b517-eb6b3ce3b751.101/Ports/
90069713-cf34-4948-a5ca-abc22a13c56b.eth1-2"
},
{
"@odata.id": "/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/Switches/b1b3d0e2-3860-4caf-ade4-2db7e9f6c075.102/Ports/
1655f138-2c46-49b6-aedf-61645d73ad3f.eth1-2"
}
]
}
]
}
```

**Request parameters**

| Parameter      | Type                  | Description                                                  |
| -------------- | --------------------- | ------------------------------------------------------------ |
| Name           | String<br/>(optional) | Name for the endpoint                                        |
| Description    | String<br/>(optional) | Description for the endpoint                                 |
| Redundancy[    | Array                 |                                                              |
| Mode           | String                | Redundancy mode                                              |
| RedundancySet] | Array                 | Set of redundancy ports connected to the switches.<br/>These links must be switch leaf ports URIs. |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Endpoints/1cf55323-c1be-43d6-bc51-7ea0d06190d8
Date:Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#Endpoint.Endpoint",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351.1/
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
416c-8c82-413153a3c351.1/Switches/
8bfa29b9-7fec-412d-8b29-042df4ba46f5.101/Ports/
903f2727-2bf8-49b1-8ebd-97729a8f1460.eth1-2"
},
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351.1/Switches/e941a68e-4ffc-4d65-
b3a5-3afe84f73fd7.102/Ports/43730998-10fe-491e-94a9-f48eeaa1e202.eth1-2"
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
| **Returns**        | JSON schema representing the created zone                    |
| **Response code**  | On success, `201 Created`                                    |
| **Authentication** | Yes                                                          |

<blockquote> NOTE: MultiVLAN is supported for the creation of multiple zone of endpoints.</blockquote>


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
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/Zones/451a7e26-00a4-4139-87b0-49e419bfa1ee"
}
],
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/AddressPools/1b695701-6ce3-457e-a530-2bc55cac5fc7"
}
],
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
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
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/Zones/451a7e26-00a4-4139-87b0-49e419bfa1ee"
}
],
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
1/AddressPools/1b695701-6ce3-457e-a530-2bc55cac5fc7"
}
],
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
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
| Description         | String<br/>(optional)      | Description for the zone.<br />**NOTE**: Ensure that there are no spaces. |
| ZoneType            | String<br/>(required)<br/> | The type of the zone to be created. Options include:<br/>• ZoneOfZones<br/>• ZoneOfEndpoints<br/>• Default<br/>The type of the zone for a zone of endpoints is<br/>ZoneOfEndpoints. |
| Links{              | Object<br/>(required)      | Contains references to other resources that are related to the zone |
| ContainedByZones [{ | Array<br/>(required)       | Represents an array of ZoneOfZones for the zone being created |
| @odata.id }]        | String                     | Link to a Zone of zones                                      |
| AddressPools [{     | Array<br/>(required)       | Represents an array of address pools linked with a ZoneOfZones |
| @odata.id }]        | String                     | Link to an address pool                                      |
| Endpoints [{        | Array<br/>(required)       | Represents an array of endpoints to be included in the zone  |
| @odata.id }]        | String                     | Link to an endpoint                                          |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow: "GET", "PUT", "POST", "PATCH", "DELETE"
Location: /redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/8e18e640-a91b-4d9b-9810-b63af3d9ce9b
Date: Wed, 31 Mar 2021 12:55:55 GMT-20h 45m
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
416c-8c82-413153a3c351.1/Endpoints/1cf55323-c1be-43d6-bc51-7ea0d06190d8"
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

| **Method**         | `PATCH`                                                |
| ------------------ | ------------------------------------------------------ |
| **URI**            | `/redfish/v1/Fabrics/{fabricid}/Zones/{zoneid}`        |
| **Description**    | This operation updates a zone of endpoints.            |
| **Returns**        | JSON schema representing the updated zone of endpoints |
| **Response code**  | On success, `200 OK`                                   |
| **Authentication** | Yes                                                    |

>**curl command**


```
curl -i PATCH \
-H "X-Auth-Token:{X-Auth-Token}" \
-d \
'{
"Links":{
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
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
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
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
Location:/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351:1/Zones/8e18e640-a91b-4d9b-9810-b63af3d9ce9b
Date:Fri, 02 Apr 2021 07:39:26 GMT-2d 22h
```

>**Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#Zone.Zone",
"@odata.id":"/redfish/v1/Fabrics/a127eedc-c29b-416c-8c82-413153a3c351.1/
Zones/8e18e640-a91b-4d9b-9810-b63af3d9ce9b",
"@odata.type":"#Zone.v1_4_0.Zone",
"Description":"Zone of endpoints",
"Id":"8e18e640-a91b-4d9b-9810-b63af3d9ce9b",
"Links":{
"AddressPools":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351.1/AddressPools/bb2cd119-01e5-499d-8465-
c219ad891842"
}
],
"ContainedByZones":[
{
"@odata.id":"/redfish/v1/Fabrics/a127eedcc29b-
416c-8c82-413153a3c351.1/Zones/6415d9aa-47a3-439d-93bb-5b23dccf5d60"
}
],
"ContainedByZones@odata.count":1,
"Endpoints":[
{
"@odata.id":"/redfish/v1/Fabrics/16b17167-de3e-483d-9f6daad629b8829b.
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

# Mapping of Redfish logical entities to Cisco ACI entities

| Redfish Logical Entity                                       | Equivalent Cisco ACI entity                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| Default zone                                                 | Tenant                                                       |
| *ZoneOfZone*                                                 | Application profile and VRF                                  |
| *VLANIdentifierAddressRange of ZoneOfZoneAddressPool*        | VLAN range of a domain                                       |
| Redundant ports of an endpoint                               | VPC Policy Group                                             |
| *ZoneOfEndpoints*                                            | BridgeDomain and Application EPGs (Endpoint Group)           |
| *GatewayIPAddress of ZoneOfEndpoints'* addresspool           | Subnet in BridgeDomain                                       |
| *VLANIdentifierAddressRange of ZoneOfEndpoints'* addresspool | VLAN in StaticPort of Application EPGs                       |
| Health and status of Redfish fabrics, switches, and ports    | Health and status of ACI fabrics, switches, and ports. See the following table |

**Mapping of Redfish fabric health to ACI fabric health range**

| Redfish representation | Equivalent ACI range |
| ---------------------- | -------------------- |
| OK                     | 91 to 100            |
| Warning                | 31 to 90             |
| Critical               | Lesser than 31       |

testing
