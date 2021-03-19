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

func GetChassisCollection(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	var members []*model.Link
	for _, chassis := range capdata.ChassisData {
		members = append(members, &model.Link{
			Oid: chassis.Oid,
		})
	}
	chassisCollection := model.Collection{
		ODataContext: "/ODIM/v1/$metadata#ChassisCollection.ChassisCollection",
		ODataID:      uri,
		ODataType:    "#ChassisCollection.ChassisCollection",
		Description:  "Chassis view",
		Name:         "ChassisCollection",
		Members:      members,
		MembersCount: len(members),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(chassisCollection)

}

func GetChassis(ctx iris.Context) {
	chassisID := ctx.Params().Get("id")
	var respData *model.Chassis
	for _, chassis := range capdata.ChassisData {
		if chassis.ID == chassisID {
			respData = chassis
		}
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(respData)
}
