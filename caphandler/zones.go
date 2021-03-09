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
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"
	aciModels "github.com/ciscoecosystem/aci-go-client/models"
	iris "github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// CreateApplicationProfile creates Application profiles using APIC
func CreateApplicationProfile(name string, tenant string, description string, fvApattr aciModels.ApplicationProfileAttributes) (*aciModels.ApplicationProfile, error) {
	aciServiceManager := caputilities.GetConnection()

	return aciServiceManager.CreateApplicationProfile(name, tenant, description, fvApattr)
}

// CreateVRF creates VRF's using APIC
func CreateVRF() {
	return
}

// GetZones returns the collection of zones present under a fabric
func GetZones(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	zonesData, ok := capdata.FabricToZoneDataStore[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Zone data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Zone", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}

	var members = []*model.Link{}
	for i := 0; i < len(zonesData); i++ {
		members = append(members, &model.Link{
			Oid: zonesData[i],
		})
	}
	zoneCollection := model.Collection{
		ODataContext: "/ODIM/v1/$metadata#ZoneCollection.ZoneCollection",
		ODataID:      uri,
		ODataType:    "#ZoneCollection.ZoneCollection",
		Description:  "ZoneCollection view",
		Name:         "Zones",
		Members:      members,
		MembersCount: len(members),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(zoneCollection)
}

// GetZone returns a specific zone present under a fabric
func GetZone(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	_, ok := capdata.FabricDataStore.Data[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}

	respData, ok := capdata.ZoneDataStore[uri]
	if !ok {
		errMsg := fmt.Sprintf("Zone data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Zone", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(respData)
}

// CreateZone default function called for creation of any type of zone
func CreateZone(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	_, ok := capdata.FabricDataStore.Data[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}

	var zone model.Zone
	err := ctx.ReadJSON(&zone)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the  request: " + err.Error()
		log.Error(errorMessage)
		resp := updateErrorResponse(response.MalformedJSON, errorMessage, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	if zone.ZoneType != "Default" {
		ctx.StatusCode(http.StatusNotImplemented)
		return
	}
	if zone.ZoneType == "ZoneOfZones" {
		ctx.StatusCode(http.StatusNotImplemented)
		//resp, statusCode := CreateZoneOfZones(zone)
		//ctx.StatusCode(statusCode)
		//ctx.JSON(resp)
		return
	}
	if zone.ZoneType == "Default" {
		resp, statusCode := CreateDefaultZone(zone)
		if statusCode != http.StatusCreated {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		defaultZoneID := uuid.NewV4().String()
		zone.ID = defaultZoneID
		zone.ODataContext = "/ODIM/v1/$metadata#Zone.Zone"
		zone.ODataType = "#Zone.v1_4_0.Zone"
		zone.ODataID = fmt.Sprintf("%s/%s", uri, defaultZoneID)
		common.SetResponseHeader(ctx, map[string]string{
			"Location": zone.ODataID,
		})
		data, ok := capdata.FabricToZoneDataStore[fabricID]
		if ok {
			data = append(data, zone.ODataID)
			capdata.FabricToZoneDataStore[fabricID] = data
		} else {
			capdata.FabricToZoneDataStore[fabricID] = []string{zone.ODataID}
		}
		log.Println("HERE")
		log.Println(zone.ODataID)
		capdata.ZoneDataStore[zone.ODataID] = &zone
		ctx.StatusCode(statusCode)
		ctx.JSON(zone)
	}
}

// CreateDefaultZone creates a zone of type 'Default'
func CreateDefaultZone(zone model.Zone) (interface{}, int) {
	var tenantAttributesStruct aciModels.TenantAttributes
	tenantAttributesStruct.Name = zone.Name
	aciClient := caputilities.GetClient()
	//var resp map[string]interface{}
	resp, err := aciClient.CreateTenant(zone.Name, zone.Description, tenantAttributesStruct)
	if err != nil {
		errMsg := "Error while creating default Zone: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	//if resp["totalCount"].(string) != "0"{
	//var errResponse capModels.ErrorResponse
	log.Println(resp)
	log.Println(err)
	//return resp, http.StatusBadRequest
	//}
	return resp, http.StatusCreated
}

// DeleteZone deletes the zone from the resource
func DeleteZone(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	_, ok := capdata.FabricDataStore.Data[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}

	respData, ok := capdata.ZoneDataStore[uri]
	if !ok {
		errMsg := fmt.Sprintf("Zone data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Zone", uri})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	if respData.Links != nil {
		if respData.Links.ContainsZonesCount != 0 {
			errMsg := fmt.Sprintf("Zone cannot be deleted as there are dependent resources still tied to it")
			log.Error(errMsg)
			resp := updateErrorResponse(response.ResourceCannotBeDeleted, errMsg, []interface{}{"Zone", uri})
			ctx.StatusCode(http.StatusNotAcceptable)
			ctx.JSON(resp)
			return
		}
	}
	aciClient := caputilities.GetClient()
	err := aciClient.DeleteTenant(respData.Name)
	if err != nil {
		errMsg := "Error while deleting Zone: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	delete(capdata.ZoneDataStore, uri)
	ctx.JSON(http.StatusNoContent)

}
