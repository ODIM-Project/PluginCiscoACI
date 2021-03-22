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
	"errors"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/capresponse"
	"github.com/ODIM-Project/PluginCiscoACI/db"

	iris "github.com/kataras/iris/v12"
)

//GetChassisCollection fetches details of the chassis collection
func GetChassisCollection(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	var members []*model.Link
	chassisData, err := capmodel.GetAllSwitchChassis("")
	if err != nil {
		capresponse.SetErrorResponse(ctx, http.StatusInternalServerError, response.InternalError, err.Error(), nil)
		return
	}
	for _, chassis := range chassisData {
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

//GetChassis fetches details of the chassis
func GetChassis(ctx iris.Context) {
	chassisID := ctx.Params().Get("id")
	data, err := capmodel.GetSwitchChassis(chassisID)
	if errors.Is(err, db.ErrorKeyNotFound) {
		msgArgs := []interface{}{"Chassis", chassisID}
		capresponse.SetErrorResponse(ctx, http.StatusNotFound, response.ResourceNotFound, err.Error(), msgArgs)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(data)
	return
}
