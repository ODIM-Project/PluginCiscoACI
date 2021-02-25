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
	"net/http"
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
	portID := ctx.Params().Get("portID")
	portData := capdata.PortDataStore[portID].(map[string]interface{})
	portResponse := model.Port{
		ODataContext: "/ODIM/v1/$metadata#Port.Port",
		ODataID:      uri,
		ODataType:    "#Port.v1_3_0.Port",
		ID:           portID,
		PortID:       portData["id"].(string),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(portResponse)
}