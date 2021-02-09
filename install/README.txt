This document will provide steps on how to create docker image & deploy through Helm chart.

- Clone the CISCO ACI repository
- Run below commands to set envoirnment variables:
   $ export ODIMRA_GROUP_ID=<group_id>
   $ export ODIMRA_USER_ID=<user_id>

- Run below command to create docker image
   $ cd PluginCiscoACI
   $ ./build_images.sh

- Run below command to deploy ACI plugin using helm chart
   $ cd PluginCiscoACI/install/Kubernetes/helmcharts
   $ ./deploying_chart.sh <namespace>

Note: Before deploying ACI plugin we need to populate values(now manually by editing the required files) to be picked by helm chart before deployment.
Folowing files are required to be edited with the user configured values:
   - PluginCiscoACI/install/Kubernetes/helmcharts/aci-platformconfig/values.yaml
   - PluginCiscoACI/install/Kubernetes/helmcharts/aciplugin-config/values.yaml
   - PluginCiscoACI/install/Kubernetes/helmcharts/aciplugin-pv-pvc/values.yaml
   - PluginCiscoACI/install/Kubernetes/helmcharts/aciplugin/values.yaml

- Run below command to undeploy ACI plugin 
   $ cd PluginCiscoACI/install/Kubernetes/helmcharts
   $ ./undeploy_chart.sh <namespace>
