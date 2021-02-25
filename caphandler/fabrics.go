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
	iris "github.com/kataras/iris/v12"
	"net/http"
)

//GetFabricResource : Fetches details of the given resource from the device
func GetFabricResource(ctx iris.Context) {
	ctx.StatusCode(http.StatusNotImplemented)
}

// GetFabricData fetches the fabric information
func GetFabricData(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")

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
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(fabricResponse)

}
