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
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

// GetSwitchCollection fetches the switches which are linked to that fabric
func GetSwitchCollection(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	// get all switches which are store under that fabric

	fabricData := capdata.FabricDataStore.Data[fabricID]

	var members = []*model.Link{}

	for i := 0; i < len(fabricData.SwitchData); i++ {
		members = append(members, &model.Link{
			Oid: uri + "/" + fabricData.SwitchData[i],
		})
	}

	switchCollectionResponse := model.Collection{
		ODataContext: "/ODIM/v1/$metadata#SwitchCollection.SwitchCollection",
		ODataID:      uri,
		ODataType:    "#SwitchCollection.SwitchCollection",
		Description:  "Switches view",
		Name:         "Switches",
		Members:      members,
		MembersCount: len(members),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(switchCollectionResponse)
}

// GetSwitchInfo fetches the switch info for given swith id
func GetSwitchInfo(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	switchID := ctx.Params().Get("rid")
	// Get the switch data from the memory
	switchResponse := capdata.SwitchDataStore.Data[switchID]
	switchResponse.ODataID = uri
	switchResponse.Ports = &model.Link{
		Oid: uri + "/Ports",
	}
	fabricID := ctx.Params().Get("id")
	fabricData := capdata.FabricDataStore.Data[fabricID]
	switchResponse.Status = &model.Status{
		State:  "Enabled",
		Health: getSwitchHealthData(fabricData.PodID, switchID),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(switchResponse)
}

func getSwitchHealthData(podID, switchID string) string {
	switchIDData := strings.Split(switchID, ":")
	switchHealthResposne, err := caputilities.GetSwitchHealth(podID, switchIDData[1])
	if err != nil {
		log.Error("Unable to get Health of switch " + err.Error())
		return ""
	}
	data := switchHealthResposne.IMData[0].HealthData.Attributes
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
