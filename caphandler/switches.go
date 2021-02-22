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
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// GetSwitchCollection fetches the switches which are linked to that fabric
func GetSwitchCollection(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Println("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}

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
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Println("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}
	switchID := ctx.Params().Get("rid")
	// Get the switch data from the memory
	switchData := capdata.SwitchDataStore.Data[switchID]
	switchUUIDData := strings.Split(switchID, ":")
	switchResponse := model.Switch{
		ODataContext: "/ODIM/v1/$metadata#Switch.Switch",
		ODataID:      uri,
		ODataType:    "#Switch.v1_4_0.Switch",
		ID:           switchID,
		Name:         switchData.Name,
		SwitchType:   "Ethernet",
		UUID:         switchUUIDData[0],
		SerialNumber: switchData.Serial,
		Ports: &model.Link{
			Oid: uri + "/Ports",
		},
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(switchResponse)
}
