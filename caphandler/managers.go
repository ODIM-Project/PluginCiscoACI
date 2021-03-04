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

//Package caphandler ...
package caphandler

import (
	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	pluginConfig "github.com/ODIM-Project/PluginCiscoACI/config"
	iris "github.com/kataras/iris/v12"
	"net/http"
)

//GetManagersCollection Fetches details of the manager collection
func GetManagersCollection(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	var members = []*model.Link{
		&model.Link{
			Oid: "/ODIM/v1/Managers/" + pluginConfig.Data.RootServiceUUID,
		},
	}

	managers := model.Collection{
		ODataContext: "/ODIM/v1/$metadata#ManagerCollection.ManagerCollection",
		ODataID:      uri,
		ODataType:    "#ManagerCollection.ManagerCollection",
		Description:  "Managers view",
		Name:         "Managers",
		Members:      members,
		MembersCount: len(members),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(managers)
	return

}

//GetManagersInfo Fetches details of the given manager info
func GetManagersInfo(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	// Get all switch data uri
	managedSwitches := []model.Link{}
	for fabricID, fabricData := range capdata.FabricDataStore.Data {
		for i := 0; i < len(fabricData.SwitchData); i++ {
			managedSwitches = append(managedSwitches, model.Link{
				Oid: "/ODIM/v1/Fabrics/" + fabricID + "/Switches/" + fabricData.SwitchData[i],
			})
		}
	}

	managers := model.Manager{
		ODataContext:    "/ODIM/v1/$metadata#Manager.Manager",
		ODataID:         uri,
		ODataType:       "#Manager.v1_10_0.Manager",
		Name:            pluginConfig.Data.PluginConf.ID,
		ManagerType:     "Service",
		ID:              pluginConfig.Data.RootServiceUUID,
		UUID:            pluginConfig.Data.RootServiceUUID,
		FirmwareVersion: pluginConfig.Data.FirmwareVersion,
		Status: &model.Status{
			State:  "Enabled",
			Health: "OK",
		},
		Links: &model.ManagerLinks{
			ManagerForSwitches:      &managedSwitches,
			ManagerForSwitchesCount: len(managedSwitches),
		},
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(managers)
	return

}

func getInfoFromDevice(uri string, deviceDetails capmodel.Device, ctx iris.Context) {
	// TODO: implementation pending
}
