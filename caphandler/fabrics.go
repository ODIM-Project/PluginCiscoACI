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
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

//GetFabricResource : Fetches details of the given resource from the device
func GetFabricResource(ctx iris.Context) {
	ctx.StatusCode(http.StatusNotImplemented)
}

// GetFabricData fetches the fabric information
func GetFabricData(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	fabricData, ok := capdata.FabricDataStore.Data[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", uri})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	var fabricResponse = model.Fabric{
		ODataContext: "/ODIM/v1/$metadata#Fabric.Fabric",
		ODataID:      uri,
		ODataType:    "#Fabric.v1_2_0.Fabric",
		Name:         "ACI Fabric",
		ID:           fabricID,
		AddressPools: &model.Link{
			Oid: uri + "/AddressPools",
		},
		Endpoints: &model.Link{
			Oid: uri + "/Endpoints",
		},
		Switches: &model.Link{
			Oid: uri + "/Switches",
		},
		Zones: &model.Link{
			Oid: uri + "/Zones",
		},
		FabricType: "Ethernet",
		MaxZones:   800,
		Status: &model.Status{
			State:  "Enabled",
			Health: getFabricHealthData(fabricData.PodID),
		},
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(fabricResponse)

}

func getFabricHealthData(podID string) string {
	fabricHealthResposne, err := caputilities.GetFabricHealth(podID)
	if err != nil {
		log.Info("Unable to get fabric health" + err.Error())
		return ""
	}
	log.Info(fabricHealthResposne)
	data := fabricHealthResposne.IMData[0].FabricHealthData.Attributes
	currentHealthValue := data["cur"].(string)
	healthValue, err := strconv.Atoi(currentHealthValue)
	if err != nil {
		log.Error("Unable to convert current Health value:" + currentHealthValue + " go the error" + err.Error())
		return ""
	}
	if healthValue > 90 {
		return "OK"
	} else if healthValue <= 90 && healthValue < 30 {
		return "Warning"
	}
	return "Critical"
}
