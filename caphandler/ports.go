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

// GetPortCollection fetches the ports  which are linked to that switch
func GetPortCollection(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	switchID := ctx.Params().Get("switchID")

	// get all port which are store under that switch
	portData := capdata.SwitchToPortDataStore[switchID]

	var members = []*model.Link{}

	for i := 0; i < len(portData); i++ {
		members = append(members, &model.Link{
			Oid: uri + "/" + portData[i],
		})
	}

	portCollectionResponse := model.Collection{
		ODataContext: "/ODIM/v1/$metadata#PortCollection.PortCollection",
		ODataID:      uri,
		ODataType:    "#PortCollection.PortCollection",
		Description:  "PortCollection view",
		Name:         "Ports",
		Members:      members,
		MembersCount: len(members),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(portCollectionResponse)
}

// GetPortInfo fetches the port info for given port id
func GetPortInfo(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	switchID := ctx.Params().Get("rid")
	fabricID := ctx.Params().Get("id")
	fabricData := capdata.FabricDataStore.Data[fabricID]
	portID := ctx.Params().Get("portID")
	portData := capdata.PortDataStore[portID]
	portData.ODataID = uri
	getPortAddtionalAttributes(fabricData.PodID, switchID, portData)
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(portData)

}

func getPortAddtionalAttributes(fabricID, switchID string, p *model.Port) {
	switchIDData := strings.Split(switchID, ":")
	PortInfoResponse, err := caputilities.GetPortInfo(fabricID, switchIDData[1], p.PortID)
	if err != nil {
		log.Error("Unable to get addtional port info " + err.Error())
		return
	}
	portInfoData := PortInfoResponse.IMData[0].PhysicalInterface.Attributes
	operationState := portInfoData["operSt"].(string)
	if operationState == "up" {
		p.LinkState = "Enabled"
		p.LinkStatus = "LinkUp"
	} else {
		p.LinkState = "Disabled"
		p.LinkStatus = "LinkDown"

	}
	portsHelathResposne, err := caputilities.GetPortHealth(fabricID, switchIDData[1], p.PortID)
	if err != nil {
		log.Error("Unable to get helath of switch " + err.Error())
		return
	}

	helathdata := portsHelathResposne.IMData[0].HelathData.Attributes
	currentHealthValue := helathdata["cur"].(string)
	healthValue, err := strconv.Atoi(currentHealthValue)
	if err != nil {
		log.Error("Unable to convert current helath value:" + currentHealthValue + " go the error" + err.Error())
		return
	}
	var portStatus = model.Status{
		State: p.LinkState,
	}
	if healthValue > 90 {
		portStatus.Health = "OK"
	} else if healthValue <= 90 && healthValue < 30 {
		portStatus.Health = "Warning"
	} else {
		portStatus.Health = "Critical"
	}

	p.Status = &portStatus
	return
}
