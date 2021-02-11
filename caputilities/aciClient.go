//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

//Package caputilities ...
package caputilities

import (
	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	log "github.com/sirupsen/logrus"
)

// GetFabricNodeData collects the all switch and fabric  deatails from the aci
func GetFabricNodeData() ([]*models.FabricNodeMember, error) {
	aciClient := client.NewClient(config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	log.Info("token: ")
	log.Info(aciClient.AuthToken)
	serviceManager := client.NewServiceManager("", aciClient)
	return serviceManager.ListFabricNodeMember()

}
