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
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"

	aciModels "github.com/ciscoecosystem/aci-go-client/models"
	iris "github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// CreateApplicationProfile creates Application profiles using APIC
func CreateApplicationProfile(name string, tenant string, description string, fvApattr aciModels.ApplicationProfileAttributes) (*aciModels.ApplicationProfile, error) {
	aciServiceManager := caputilities.GetConnection()

	return aciServiceManager.CreateApplicationProfile(name, tenant, description, fvApattr)
}

// CreateVRF creates VRF's using APIC
func CreateVRF(name string, tenant string, description string, fvCtxattr aciModels.VRFAttributes) (*aciModels.VRF, error) {
	aciServiceManager := caputilities.GetConnection()
	return aciServiceManager.CreateVRF(name, tenant, description, fvCtxattr)
}

// GetZones returns the collection of zones present under a fabric
func GetZones(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	if _, err := capmodel.GetFabric(fabricID); err != nil {
		errMsg := fmt.Sprintf("failed to fetch fabric data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"AddressPool", uri})
		return
	}

	var members = []*model.Link{}
	zoneData, err := capmodel.GetAllZones(fabricID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", fabricID})
		return

	}
	for zoneID := range zoneData {
		members = append(members, &model.Link{
			Oid: zoneID,
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
	if _, err := capmodel.GetFabric(fabricID); err != nil {
		errMsg := fmt.Sprintf("failed to fetch fabric data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Fabric", fabricID})
		return
	}

	zoneData, err := capmodel.GetZone(fabricID, uri)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", fabricID})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(zoneData)
}

// CreateZone default function called for creation of any type of zone
func CreateZone(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	if _, err := capmodel.GetFabric(fabricID); err != nil {
		errMsg := fmt.Sprintf("failed to fetch fabric data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Fabric", fabricID})
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
	switch zone.ZoneType {
	case "Default":
		resp, statusCode := CreateDefaultZone(zone)
		if statusCode != http.StatusCreated {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		conflictFlag := false
		var defaultZoneID string
		fabricData, err := capmodel.GetAllFabric("")
		if err != nil {
			errMsg := fmt.Sprintf("failed to fetch fabric data: %s", err.Error())
			createDbErrResp(ctx, err, errMsg, nil)
			return
		}

		for fabricID := range fabricData {
			data, err := capmodel.GetAllZones(fabricID)
			if err != nil {
				errMsg := fmt.Sprintf("failed to fetch all zones data of fabric %s: %s", fabricID, err.Error())
				createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", fabricID})
				return
			}
			for _, zoneData := range data {
				if zoneData.Name == zone.Name {
					conflictFlag = true
				}
			}
		}
		if !conflictFlag {
			defaultZoneID = uuid.NewV4().String()
			if zone, err = saveZoneData(defaultZoneID, uri, fabricID, zone); err != nil {
				errMsg := fmt.Sprintf("failed to store default zone data for uri %s: %s", uri, err.Error())
				createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", fabricID})
				return
			}
		}
		common.SetResponseHeader(ctx, map[string]string{
			"Location": zone.ODataID,
		})
		ctx.StatusCode(statusCode)
		ctx.JSON(zone)
		return
	case "ZoneOfZones":
		defaultZoneLink, resp, statusCode, domainData := CreateZoneOfZones(uri, fabricID, zone)
		if statusCode != http.StatusCreated {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		conflictFlag := false
		var defaultZoneID string
		fabricData, err := capmodel.GetAllFabric("")
		if err != nil {
			errMsg := fmt.Sprintf("failed to fetch fabric data: %s", err.Error())
			createDbErrResp(ctx, err, errMsg, nil)
			return
		}

		for fabricID := range fabricData {
			data, err := capmodel.GetAllZones(fabricID)
			if err != nil {
				errMsg := fmt.Sprintf("failed to fetch all zones data of fabric %s: %s", fabricID, err.Error())
				createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", fabricID})
				return
			}
			for _, zoneData := range data {
				if zoneData.Name == zone.Name {
					conflictFlag = true
				}
			}
		}
		if !conflictFlag {
			defaultZoneID = uuid.NewV4().String()
			if zone, err = saveZoneData(defaultZoneID, uri, fabricID, zone); err != nil {
				errMsg := fmt.Sprintf("failed to store zone of zone data for uri %s: %s", uri, err.Error())
				createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", defaultZoneID})
				return
			}
			if err = saveZoneToDomainDNData(zone.ODataID, domainData); err != nil {
				errMsg := fmt.Sprintf("failed to update zone domain data for uri %s: %s", uri, err.Error())
				createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", zone.ODataID})
				return
			}
		}
		if err = updateZoneData(fabricID, defaultZoneLink, zone); err != nil {
			errMsg := fmt.Sprintf("failed to update zone data for uri %s: %s", uri, err.Error())
			createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", defaultZoneLink})
			return
		}
		if err = updateAddressPoolData(fabricID, zone.ODataID, zone.Links.AddressPools[0].Oid, "Add"); err != nil {
			errMsg := fmt.Sprintf("failed to update AddressPool data for uri %s: %s", uri, err.Error())
			createDbErrResp(ctx, err, errMsg, []interface{}{"AddressPool", zone.Links.AddressPools[0].Oid})
			return
		}
		common.SetResponseHeader(ctx, map[string]string{
			"Location": zone.ODataID,
		})
		ctx.StatusCode(statusCode)
		ctx.JSON(zone)
		return
	case "ZoneOfEndpoints":
		zoneofZoneOID, resp, statusCode := createZoneOfEndpoints(uri, fabricID, zone)
		if statusCode != http.StatusCreated {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		zoneID := uuid.NewV4().String()
		if zone, err = saveZoneData(zoneID, uri, fabricID, zone); err != nil {
			errMsg := fmt.Sprintf("failed to store zone of endpoints data for uri %s: %s", uri, err.Error())
			createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", zoneID})
			return
		}
		if err = updateZoneData(fabricID, zoneofZoneOID, zone); err != nil {
			errMsg := fmt.Sprintf("failed to update zone data for uri %s: %s", uri, err.Error())
			createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", zoneofZoneOID})
			return
		}
		if err = updateAddressPoolData(fabricID, zone.ODataID, zone.Links.AddressPools[0].Oid, "Add"); err != nil {
			errMsg := fmt.Sprintf("failed to update AddressPool data for uri %s: %s", uri, err.Error())
			createDbErrResp(ctx, err, errMsg, []interface{}{"AddressPool", zone.Links.AddressPools[0].Oid})
			return
		}
		common.SetResponseHeader(ctx, map[string]string{
			"Location": zone.ODataID,
		})
		ctx.StatusCode(statusCode)
		ctx.JSON(zone)
		return
	default:
		ctx.StatusCode(http.StatusNotImplemented)
		return
	}
}

func saveZoneData(defaultZoneID string, uri string, fabricID string, zone model.Zone) (model.Zone, error) {
	zone.ID = defaultZoneID
	zone.ODataContext = "/ODIM/v1/$metadata#Zone.Zone"
	zone.ODataType = "#Zone.v1_4_0.Zone"
	zone.ODataID = fmt.Sprintf("%s/%s", uri, defaultZoneID)
	zone.Status = &model.Status{}
	zone.Status.State = "Enabled"
	zone.Status.Health = "OK"
	if zone.Links != nil {
		if zone.Links.ContainedByZones != nil {
			zone.Links.ContainedByZonesCount = len(zone.Links.ContainedByZones)
		}
	}
	if err := capmodel.SaveZone(fabricID, zone.ODataID, &zone); err != nil {
		return zone, err
	}
	return zone, nil
}

// CreateDefaultZone creates a zone of type 'Default'
func CreateDefaultZone(zone model.Zone) (interface{}, int) {
	var tenantAttributesStruct aciModels.TenantAttributes
	tenantAttributesStruct.Name = zone.Name
	aciClient := caputilities.GetConnection()
	//var tenantList []*aciModels.Tenant
	tenantList, err := aciClient.ListTenant()
	if err != nil {
		errMsg := "Error while creating default Zone: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	for _, tenant := range tenantList {
		if tenant.TenantAttributes.Name == zone.Name {
			errMsg := "Default zone already exists with name: " + zone.Name
			resp := updateErrorResponse(response.ResourceAlreadyExists, errMsg, []interface{}{"DefaultZone", tenant.TenantAttributes.Name, zone.Name})
			return resp, http.StatusConflict
		}

	}

	resp, err := aciClient.CreateTenant(zone.Name, zone.Description, tenantAttributesStruct)
	if err != nil {
		errMsg := "Error while creating default Zone: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return resp, http.StatusCreated
}

// DeleteZone deletes the zone from the resource
func DeleteZone(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	if _, err := capmodel.GetFabric(fabricID); err != nil {
		errMsg := fmt.Sprintf("failed to fetch fabric data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Fabric", fabricID})
		return
	}

	//TODO: Get list of zones which are pre-populated from onstart and compare the members for item not present in odim but present in ACI

	zoneData, err := capmodel.GetZone(fabricID, uri)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", uri})
		return
	}
	if zoneData.Links != nil {
		if zoneData.Links.ContainsZonesCount != 0 {
			errMsg := fmt.Sprintf("Zone cannot be deleted as there are dependent resources still tied to it")
			log.Error(errMsg)
			resp := updateErrorResponse(response.ResourceCannotBeDeleted, errMsg, []interface{}{"Zone", uri})
			ctx.StatusCode(http.StatusNotAcceptable)
			ctx.JSON(resp)
			return
		}
	}
	if zoneData.ZoneType == "ZoneOfZones" {
		err := deleteZoneOfZone(fabricID, uri, &zoneData)
		if err != nil {
			if err.Error() == "Error deleting Application Profile" {
				resp := updateErrorResponse(response.GeneralError, err.Error(), nil)
				ctx.StatusCode(http.StatusBadRequest)
				ctx.JSON(resp)
				return
			}
		}
		if err != nil {
			errMsg := fmt.Sprintf("Zone data for uri %s not found", uri)
			log.Error(errMsg)
			resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Zone", uri})
			ctx.StatusCode(http.StatusNotFound)
			ctx.JSON(resp)
			return
		}
		if err = capmodel.DeleteZone(fabricID, uri); err != nil {
			errMsg := fmt.Sprintf("failed to delete zone data for %s: %s", uri, err.Error())
			createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", uri})
			return
		}
		ctx.StatusCode(http.StatusNoContent)
	}
	if zoneData.ZoneType == "Default" {
		aciClient := caputilities.GetConnection()
		err := aciClient.DeleteTenant(zoneData.Name)
		if err != nil {
			errMsg := "Error while deleting Zone: " + err.Error()
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(resp)
			return
		}
		if err = capmodel.DeleteZone(fabricID, uri); err != nil {
			errMsg := fmt.Sprintf("failed to delete zone data for %s: %s", uri, err.Error())
			createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", uri})
			return
		}
		ctx.StatusCode(http.StatusNoContent)
	}
	if zoneData.ZoneType == "ZoneOfEndpoints" {
		resp, statusCode := deleteZoneOfEndpoints(fabricID, &zoneData)
		ctx.StatusCode(statusCode)
		ctx.JSON(resp)
	}
}

func deleteZoneOfZone(fabricID, uri string, respData *model.Zone) error {
	var parentZoneLink model.Link
	var parentZone model.Zone
	if respData.Links != nil {
		var err error
		if respData.Links.ContainedByZonesCount != 0 {
			// Assuming contained by link is only one
			parentZoneLink = respData.Links.ContainedByZones[0]
			parentZone, err = capmodel.GetZone(fabricID, parentZoneLink.Oid)
			if err != nil {
				return fmt.Errorf("failed to fetch zone data for %s: %s", parentZoneLink.Oid, err.Error())
			}
		}
		aciServiceManager := caputilities.GetConnection()
		err = aciServiceManager.DeleteApplicationProfile(respData.Name, parentZone.Name)
		if err != nil {
			errMsg := fmt.Errorf("Error deleting Application Profile:%v", err)
			return errMsg
		}
		vrfErr := aciServiceManager.DeleteVRF(respData.Name+"-VRF", parentZone.Name)
		if vrfErr != nil {
			errMsg := fmt.Errorf("Error deleting VRF:%v", err)
			return errMsg
		}
		// delete contract
		contractErr := aciServiceManager.DeleteContract(respData.Name+"-VRF-Con", parentZone.Name)
		if contractErr != nil {
			errMsg := fmt.Errorf("Error deleting Contract:%v", contractErr)
			log.Error(errMsg.Error())
			return errMsg
		}
		err = aciServiceManager.DeleteAttachableAccessEntityProfile(respData.Name + "-DOM-EntityProfile")
		if err != nil {
			errMsg := fmt.Errorf("Error deleting  domain profile:%v", contractErr)
			log.Error(errMsg.Error())
			return errMsg
		}
		err = aciServiceManager.DeletePhysicalDomain(respData.Name + "-DOM")
		if err != nil {
			errMsg := fmt.Errorf("Error deleting Physical domain:%v", contractErr)
			log.Error(errMsg.Error())
			return errMsg
		}
		err = aciServiceManager.DeleteVLANPool("static", respData.Name+"-DOM-VLAN")
		if err != nil {
			errMsg := fmt.Errorf("Error deleting Physical domain:%v", contractErr)
			log.Error(errMsg.Error())
			return errMsg
		}
		if err = updateAddressPoolData(fabricID, respData.ODataID, respData.Links.AddressPools[0].Oid, "Remove"); err != nil {
			errMsg := fmt.Errorf("Error updating addressPool data:%v", err)
			return errMsg
		}
		if err = capmodel.DeleteZone(fabricID, uri); err != nil {
			return fmt.Errorf("failed to delete zone %s: %s", uri, err.Error())
		}
		if err = capmodel.DeleteZoneDomain(uri); err != nil {
			return fmt.Errorf("failed to delete zone domain %s: %s", uri, err.Error())
		}
		for i := 0; i < len(parentZone.Links.ContainsZones); i++ {
			if parentZone.Links.ContainsZones[i].Oid == respData.ODataID {
				parentZone.Links.ContainsZones[i] = parentZone.Links.ContainsZones[len(parentZone.Links.ContainsZones)-1] // Copy last element to index i.
				parentZone.Links.ContainsZones[len(parentZone.Links.ContainsZones)-1] = model.Link{}                      // Erase last element (write zero value).
				parentZone.Links.ContainsZones = parentZone.Links.ContainsZones[:len(parentZone.Links.ContainsZones)-1]
			}
		}
		parentZone.Links.ContainsZonesCount = len(parentZone.Links.ContainsZones)

		if err = capmodel.UpdateZone(fabricID, parentZoneLink.Oid, &parentZone); err != nil {
			return fmt.Errorf("failed to update zone data for %s: %s", parentZoneLink.Oid, err.Error())
		}

		return nil
	}
	return nil
}

// CreateZoneOfZones takes the request to create zone of zones and translates to create application profiles and VRFs
func CreateZoneOfZones(uri string, fabricID string, zone model.Zone) (string, interface{}, int, *capdata.ACIDomainData) {
	var apModel aciModels.ApplicationProfileAttributes
	var vrfModel aciModels.VRFAttributes
	apModel.Name = zone.Name
	vrfModel.Name = zone.Name + "-VRF"
	if zone.Links != nil {
		if len(zone.Links.ContainedByZones) == 0 {
			errMsg := fmt.Sprintf("Zone cannot be created as there are dependent resources missing")
			log.Error(errMsg)
			resp := updateErrorResponse(response.PropertyMissing, errMsg, []interface{}{"ContainedByZones"})
			return "", resp, http.StatusBadRequest, nil
		}
	}
	log.Println("Request Body")
	log.Println(zone)
	// Assuming there is only link under ContainedByZones
	defaultZoneLinks := zone.Links.ContainedByZones
	defaultZoneLink := defaultZoneLinks[0].Oid
	respData, err := capmodel.GetZone(fabricID, defaultZoneLink)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", defaultZoneLink, err.Error())
		statusCode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"Zone", defaultZoneLink})
		return "", resp, statusCode, nil
	}
	// validate all given addresspools if it's present
	if len(zone.Links.AddressPools) == 0 {
		errorMessage := "AddressPools attribute is missing in the request"
		return "", updateErrorResponse(response.PropertyMissing, errorMessage, []interface{}{"AddressPool"}), http.StatusBadRequest, nil
	}
	if len(zone.Links.AddressPools) > 1 {
		errorMessage := "More than one AddressPool not allowed for the creation of ZoneOfZones"
		return "", updateErrorResponse(response.PropertyValueFormatError, errorMessage, []interface{}{"AddressPools", "AddressPools"}), http.StatusBadRequest, nil
	}

	addresspoolData, statusCode, resp := getAddressPoolData(fabricID, zone.Links.AddressPools[0].Oid)
	if statusCode != http.StatusOK {
		return "", resp, statusCode, nil
	}
	if addresspoolData.Ethernet.IPv4.VLANIdentifierAddressRange == nil {
		errorMessage := "Provided AddressPool doesn't contain the VLANIdentifierAddressRange"
		return "", updateErrorResponse(response.PropertyMissing, errorMessage, []interface{}{"VLANIdentifierAddressRange"}), http.StatusBadRequest, nil
	}
	aciClient := caputilities.GetConnection()
	appProfileList, err := aciClient.ListApplicationProfile(respData.Name)
	if err != nil && !strings.Contains(err.Error(), "Object may not exists") {
		errMsg := fmt.Sprintf("Zone cannot be created, error while retriving existing Application profiles: " + err.Error())
		resp := updateErrorResponse(response.PropertyMissing, errMsg, []interface{}{"ContainedByZones"})
		return "", resp, http.StatusBadRequest, nil
	}
	for _, appProfile := range appProfileList {
		if appProfile.ApplicationProfileAttributes.Name == zone.Name {
			errMsg := "Application profile already exists with name: " + zone.Name
			resp := updateErrorResponse(response.ResourceAlreadyExists, errMsg, []interface{}{"ApplicationProfile", appProfile.ApplicationProfileAttributes.Name, zone.Name})
			return "", resp, http.StatusConflict, nil
		}
	}
	vrfList, err := aciClient.ListVRF(respData.Name)
	if err != nil && !strings.Contains(err.Error(), "Object may not exists") {
		errMsg := fmt.Sprintf("Zone cannot be created, error while retriving existing VRFs: " + err.Error())
		log.Error(errMsg)
		resp := updateErrorResponse(response.PropertyMissing, errMsg, []interface{}{"ContainedByZones"})
		return "", resp, http.StatusBadRequest, nil
	}
	for _, vrf := range vrfList {
		if vrf.VRFAttributes.Name == vrfModel.Name {
			errMsg := "VRF already exists with name: " + vrfModel.Name
			resp := updateErrorResponse(response.ResourceAlreadyExists, errMsg, []interface{}{"VRF", vrf.VRFAttributes.Name, vrfModel.Name})
			return "", resp, http.StatusConflict, nil
		}
	}

	apResp, err := CreateApplicationProfile(zone.Name, respData.Name, respData.Description, apModel)
	if err != nil {
		errMsg := "Error while creating application profile: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return "", resp, http.StatusBadRequest, nil
	}
	_, vrfErr := CreateVRF(vrfModel.Name, respData.Name, respData.Description, vrfModel)
	if vrfErr != nil {
		errMsg := "Error while creating application profile: " + vrfErr.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return "", resp, http.StatusBadRequest, nil
	}
	// create contract with name vrf and suffix-Con
	resp, statusCode = createContract(vrfModel.Name, respData.Name, zone.Name)
	if statusCode != http.StatusCreated {
		return "", resp, statusCode, nil
	}
	// create the domain for the given addresspool
	var domainData *capdata.ACIDomainData
	resp, statusCode, domainData = createACIDomain(addresspoolData, zone.Name)
	if statusCode != http.StatusCreated {
		return "", resp, statusCode, nil
	}
	return defaultZoneLink, apResp, http.StatusCreated, domainData

}

func updateZoneData(fabricID, defaultZoneLink string, zone model.Zone) error {
	defaultZoneData, err := capmodel.GetZone(fabricID, defaultZoneLink)
	if err != nil {
		return fmt.Errorf("failed to fetch zone data of %s: %s", defaultZoneLink, err.Error())
	}
	if defaultZoneData.Links == nil {
		defaultZoneData.Links = &model.ZoneLinks{}
	}
	if defaultZoneData.Links.ContainsZones == nil {
		var containsList []model.Link
		log.Println("List of contains")
		log.Println(containsList)
		var link model.Link
		link.Oid = zone.ODataID
		containsList = append(containsList, link)
		defaultZoneData.Links.ContainsZones = containsList
		defaultZoneData.Links.ContainsZonesCount = len(containsList)
	} else {
		var link model.Link
		link.Oid = zone.ODataID
		defaultZoneData.Links.ContainsZones = append(defaultZoneData.Links.ContainsZones, link)
		defaultZoneData.Links.ContainsZonesCount = len(defaultZoneData.Links.ContainsZones)
	}

	if err = capmodel.UpdateZone(fabricID, defaultZoneLink, &defaultZoneData); err != nil {
		return fmt.Errorf("failed to update zone data of %s: %s", defaultZoneLink, err.Error())
	}
	return nil
}

func createZoneOfEndpoints(uri, fabricID string, zone model.Zone) (string, interface{}, int) {
	// Create the BridgeDomain
	// get the Tenant name from the ZoneofZone data
	//validate the request
	if zone.Links == nil {
		errorMessage := "Links attribute is missing in the request"
		return "", updateErrorResponse(response.PropertyMissing, errorMessage, []interface{}{"Links"}), http.StatusBadRequest
	}
	if zone.Links.ContainedByZones == nil {
		errorMessage := "ContainedByZones attribute is missing in the request"
		return "", updateErrorResponse(response.PropertyMissing, errorMessage, []interface{}{"ContainedByZones"}), http.StatusBadRequest

	}
	zoneofZoneURL := zone.Links.ContainedByZones[0].Oid
	// get the zone of zone data
	zoneofZoneData, err := capmodel.GetZone(fabricID, zoneofZoneURL)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", uri, err.Error())
		statusCode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"ZoneofZone", zoneofZoneURL})
		return "", resp, statusCode
	}
	// validate all given addresspools if it's present
	if len(zone.Links.AddressPools) == 0 {
		errorMessage := "AddressPools attribute is missing in the request"
		return "", updateErrorResponse(response.PropertyMissing, errorMessage, []interface{}{"AddressPool"}), http.StatusBadRequest
	}
	if len(zone.Links.AddressPools) > 1 {
		errorMessage := "More than one AddressPool not allowed for the creation of ZoneOfEndpoints"
		return "", updateErrorResponse(response.PropertyValueFormatError, errorMessage, []interface{}{"AddressPools", "AddressPools"}), http.StatusBadRequest
	}

	addresspoolData, statusCode, resp := getAddressPoolData(fabricID, zone.Links.AddressPools[0].Oid)
	if statusCode != http.StatusOK {
		return "", resp, statusCode
	}
	// validate all given addresspools if it's present
	if addresspoolData.Links != nil && len(addresspoolData.Links.Zones) > 0 {
		errorMessage := fmt.Sprintf("Given AddressPool %s is assingned to other ZoneofEndpoints", zone.Links.AddressPools[0].Oid)
		return "", updateErrorResponse(response.ResourceInUse, errorMessage, []interface{}{"AddressPools", "AddressPools"}), http.StatusBadRequest
	}

	// validate the given addresspool
	if addresspoolData.Ethernet.IPv4.GatewayIPAddress == "" {
		errorMessage := fmt.Sprintf("Given AddressPool %s doesn't contain the GatewayIPAddress ", zone.Links.AddressPools[0].Oid)
		return "", updateErrorResponse(response.PropertyMissing, errorMessage, []interface{}{"GatewayIPAddress"}), http.StatusBadRequest
	}
	if addresspoolData.Ethernet.IPv4.VLANIdentifierAddressRange.Lower != addresspoolData.Ethernet.IPv4.VLANIdentifierAddressRange.Upper {
		errorMessage := fmt.Sprintf("Given AddressPool %s VLANIdentifierAddressRange Lower and Upper values are not matching ", zone.Links.AddressPools[0].Oid)
		return "", updateErrorResponse(response.PropertyUnknown, errorMessage, []interface{}{"VLANIdentifierAddressRange"}), http.StatusBadRequest
	}

	// Get the default zone data
	defaultZoneURL := zoneofZoneData.Links.ContainedByZones[0].Oid
	defaultZoneData, err := capmodel.GetZone(fabricID, defaultZoneURL)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", uri, err.Error())
		statusCode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"Zone", defaultZoneURL})
		return "", resp, statusCode
	}
	// Get the endpoints from the db
	endPointData := make(map[string]*capdata.EndpointData)
	for i := 0; i < len(zone.Links.Endpoints); i++ {
		data, statusCode, resp := getEndpointData(fabricID, zone.Links.Endpoints[i].Oid)
		if statusCode != http.StatusOK {
			return "", resp, statusCode
		}
		endPointData[zone.Links.Endpoints[i].Oid] = &data
	}
	// get domain from given addresspool native vlan from config
	domainData, err := getZoneTODomainDNData(zoneofZoneURL)
	if err != nil {
		errMsg := fmt.Sprintf("Domain not found for  %s", zoneofZoneURL)
		log.Error(errMsg)
		return "", updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{zoneofZoneURL, "Domain"}), http.StatusNotFound
	}
	bdResp, bdDN, statusCode := createBridgeDomain(defaultZoneData.Name, zone)
	if statusCode != http.StatusCreated {
		return "", bdResp, statusCode
	}

	// create the subnet for BD for all given address pool
	resp, statusCode = createSubnets(defaultZoneData.Name, zone.Name, addresspoolData)
	if statusCode != http.StatusCreated {
		return "", resp, statusCode
	}
	// link bridgedomain to vrf
	resp, statusCode = linkBDtoVRF(bdDN, zoneofZoneData.Name+"-VRF")
	if statusCode != http.StatusCreated {
		return "", resp, statusCode
	}
	resp, statusCode = applicationEPGOperation(defaultZoneData.Name, zoneofZoneData.Name, zone.Name, domainData, endPointData, addresspoolData.Ethernet.IPv4.VLANIdentifierAddressRange.Lower)
	return zoneofZoneURL, resp, statusCode
}

func createBridgeDomain(tenantName string, zone model.Zone) (interface{}, string, int) {
	var bridgeDomainAttributes aciModels.BridgeDomainAttributes
	bridgeDomainAttributes.Name = zone.Name
	aciClient := caputilities.GetConnection()
	//var tenantList []*aciModels.Tenant
	bridgeDomainList, err := aciClient.ListBridgeDomain(tenantName)
	if err != nil && !strings.Contains(err.Error(), "Object may not exists") {
		errMsg := "Error while creating Zone endpoints: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, "", http.StatusBadRequest
	}
	for _, bd := range bridgeDomainList {
		if bd.Name == zone.Name {
			errMsg := "ZoneOfEndpoints already exists with name: " + zone.Name + " for the default zone " + tenantName
			resp := updateErrorResponse(response.ResourceAlreadyExists, errMsg, []interface{}{"ZoneOfEndpoints", bd.BridgeDomainAttributes.Name, zone.Name})
			return resp, "", http.StatusConflict
		}

	}

	resp, err := aciClient.CreateBridgeDomain(zone.Name, tenantName, zone.Description, bridgeDomainAttributes)
	if err != nil {
		errMsg := "Error while creating  Zone of Endpoints: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, "", http.StatusBadRequest
	}
	return resp, resp.BaseAttributes.DistinguishedName, http.StatusCreated
}

func createSubnets(tenantName, bdName string, addresspoolData *model.AddressPool) (interface{}, int) {
	var subnetAttributes aciModels.SubnetAttributes
	subnetAttributes.Ip = addresspoolData.Ethernet.IPv4.GatewayIPAddress
	aciClient := caputilities.GetConnection()
	_, err := aciClient.CreateSubnet(subnetAttributes.Ip, bdName, tenantName, "subnet for ip"+subnetAttributes.Ip, subnetAttributes)
	if err != nil {
		errMsg := "Error while creating  Zone of Endpoints: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusCreated
}

func linkBDtoVRF(bdDN, vrfName string) (interface{}, int) {
	aciClient := caputilities.GetConnection()
	err := aciClient.CreateRelationfvRsCtxFromBridgeDomain(bdDN, vrfName)
	if err != nil {
		errMsg := "Error while creating  Zone of Endpoints: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusCreated
}

func applicationEPGOperation(tenantName, applicationProfileName, bdName string, domainData capdata.ACIDomainData, endPointData map[string]*capdata.EndpointData, nativeVLAN int) (interface{}, int) {
	//create EPG with name of bd adding -EPG suffix
	epgName := bdName + "-EPG"
	resp, appEPGDN, statusCode := createapplicationEPG(tenantName, applicationProfileName, epgName)
	if statusCode != http.StatusCreated {
		return resp, statusCode
	}
	// Link EPG to BD
	resp, statusCode = linkAPPEPGtoBD(appEPGDN, bdName)
	if statusCode != http.StatusCreated {
		return resp, statusCode
	}
	// Link EPG to Domain
	resp, statusCode = linkEpgtoDomain(appEPGDN, domainData.DomainDN)
	if statusCode != http.StatusCreated {
		return resp, statusCode
	}
	// Create static port
	for _, data := range endPointData {
		resp, statusCode = createStaticPort(epgName, tenantName, applicationProfileName, data.ACIPolicyGroupData, nativeVLAN, &domainData)
		if statusCode != http.StatusCreated {
			return resp, statusCode
		}
	}
	return nil, http.StatusCreated
}

func createapplicationEPG(tenantName, applicationProfileName, epgName string) (interface{}, string, int) {
	var epgAttributes = aciModels.ApplicationEPGAttributes{
		Name: epgName,
	}
	aciClient := caputilities.GetConnection()
	resp, err := aciClient.CreateApplicationEPG(epgName, applicationProfileName, tenantName, "Application EPG for "+epgName, epgAttributes)
	if err != nil {
		errMsg := "Error while creating  Zone of Endpoints: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, "", http.StatusBadRequest
	}
	return resp, resp.BaseAttributes.DistinguishedName, http.StatusCreated
}

func linkAPPEPGtoBD(appEPGDN, bdName string) (interface{}, int) {
	aciClient := caputilities.GetConnection()
	err := aciClient.CreateRelationfvRsBdFromApplicationEPG(appEPGDN, bdName)
	if err != nil {
		errMsg := "Error while creating  Zone of Endpoints: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusCreated
}

func linkEpgtoDomain(appEPGDN, domain string) (interface{}, int) {

	aciClient := caputilities.GetConnection()
	err := aciClient.CreateRelationfvRsDomAttFromApplicationEPG(appEPGDN, domain)
	if err != nil {
		errMsg := "Error while creating  Zone of Endpoints: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusCreated
}

func deleteZoneOfEndpoints(fabricID string, zoneData *model.Zone) (interface{}, int) {
	zoneofZoneURL := zoneData.Links.ContainedByZones[0].Oid
	// get the zone of zone data
	zoneofZoneData, err := capmodel.GetZone(fabricID, zoneofZoneURL)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", zoneofZoneURL, err.Error())
		statusCode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"Zone", zoneofZoneURL})
		return resp, statusCode

	}
	// Get the default zone data
	defaultZoneURL := zoneofZoneData.Links.ContainedByZones[0].Oid
	defaultZoneData, err := capmodel.GetZone(fabricID, defaultZoneURL)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", defaultZoneURL, err.Error())
		statusCode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"Zone", defaultZoneURL})
		return resp, statusCode
	}
	aciClient := caputilities.GetConnection()
	for i := 0; i < len(zoneData.Links.Endpoints); i++ {
		endpointData, statusCode, resp := getEndpointData(fabricID, zoneData.Links.Endpoints[i].Oid)
		if statusCode != http.StatusOK {
			return resp, statusCode
		}
		resp, statusCode = deleteRelationDomainEntityGroupInterfacePolicyGroup(endpointData.ACIPolicyGroupData.PCVPCPolicyGroupDN)
		if statusCode != http.StatusOK {
			return resp, statusCode
		}
	}
	if err = aciClient.DeleteApplicationEPG(zoneData.Name+"-EPG", zoneofZoneData.Name, defaultZoneData.Name); err != nil {
		errMsg := "Error while deleting Zone: " + err.Error()
		return updateErrorResponse(response.GeneralError, errMsg, nil), http.StatusBadRequest
	}
	err = aciClient.DeleteBridgeDomain(zoneData.Name, defaultZoneData.Name)
	if err != nil {
		errMsg := "Error while deleting Zone: " + err.Error()
		return updateErrorResponse(response.GeneralError, errMsg, nil), http.StatusBadRequest
	}
	//updating the contains zonesdata
	if zoneofZoneData.Links != nil {
		for i := 0; i < len(zoneofZoneData.Links.ContainsZones); i++ {
			if zoneofZoneData.Links.ContainsZones[i].Oid == zoneData.ODataID {
				zoneofZoneData.Links.ContainsZones[i] = zoneofZoneData.Links.ContainsZones[len(zoneofZoneData.Links.ContainsZones)-1] // Copy last element to index i.
				zoneofZoneData.Links.ContainsZones[len(zoneofZoneData.Links.ContainsZones)-1] = model.Link{}                          // Erase last element (write zero value).
				zoneofZoneData.Links.ContainsZones = zoneofZoneData.Links.ContainsZones[:len(zoneofZoneData.Links.ContainsZones)-1]
			}
		}
		zoneofZoneData.Links.ContainsZonesCount = len(zoneofZoneData.Links.ContainsZones)
		if err = capmodel.UpdateZone(fabricID, zoneofZoneURL, &zoneofZoneData); err != nil {
			errMsg := fmt.Sprintf("failed to update zone data for uri %s: %s", zoneofZoneURL, err.Error())
			statusCode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"Zone", zoneofZoneURL})
			return resp, statusCode
		}
	}
	if err = updateAddressPoolData(fabricID, zoneData.ODataID, zoneData.Links.AddressPools[0].Oid, "Remove"); err != nil {
		errMsg := fmt.Sprintf("failed to update AddressPool data for %s: %s", fabricID, err.Error())
		statusCode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"AddressPool", zoneData.Links.AddressPools[0].Oid})
		return resp, statusCode
	}
	if err = capmodel.DeleteZone(fabricID, zoneData.ODataID); err != nil {
		errMsg := fmt.Sprintf("failed to delete zone data for uri %s: %s", zoneData.ODataID, err.Error())
		statusCode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"Zone", zoneData.ODataID})
		return resp, statusCode

	}
	return nil, http.StatusNoContent
}

func createContract(vrfName, tenantName, description string) (interface{}, int) {
	contractName := vrfName + "-Con"
	contractAttributes := aciModels.ContractAttributes{
		Name:  contractName,
		Scope: "context",
	}
	aciClient := caputilities.GetConnection()
	contractResp, err := aciClient.CreateContract(contractName, tenantName, description, contractAttributes)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	// create the contract subject
	contractSubjectName := contractName + "-Subject"
	subejctatrribute := aciModels.ContractSubjectAttributes{
		Name: contractSubjectName,
	}
	subjectResp, err := aciClient.CreateContractSubject(contractSubjectName, contractName, tenantName, "Contract subject for the Contract "+contractResp.BaseAttributes.DistinguishedName, subejctatrribute)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		log.Error(errMsg)

		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	// create filter for the contract subject
	err = aciClient.CreateRelationvzRsSubjFiltAttFromContractSubject(subjectResp.BaseAttributes.DistinguishedName, "default")
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	// create vrfContract
	vzAnyAttributes := aciModels.AnyAttributes{
		MatchT: "All",
	}
	vzAnyresp, err := aciClient.CreateAny(vrfName, tenantName, "VRF any for the VRF "+vrfName, vzAnyAttributes)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	// relate VRF contract consumer
	err = aciClient.CreateRelationvzRsAnyToConsFromAny(vzAnyresp.BaseAttributes.DistinguishedName, contractName)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	err = aciClient.CreateRelationvzRsAnyToProvFromAny(vzAnyresp.BaseAttributes.DistinguishedName, contractName)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusCreated
}

func createStaticPort(epgName, tenantName, applicationProfileName string, aciPolicyGroupData *capdata.ACIPolicyGroupData, nativeVLAN int, domainData *capdata.ACIDomainData) (interface{}, int) {
	staticPathAttributes := aciModels.StaticPathAttributes{
		TDn:         aciPolicyGroupData.PolicyGroupDN,
		Encap:       fmt.Sprintf("vlan-%d", nativeVLAN),
		InstrImedcy: "immediate",
	}
	aciClient := caputilities.GetConnection()
	_, err := aciClient.CreateStaticPath(aciPolicyGroupData.PolicyGroupDN, epgName, applicationProfileName, tenantName, "", staticPathAttributes)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	// Attach the domain entity profile to given policy group
	err = aciClient.CreateRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(aciPolicyGroupData.PCVPCPolicyGroupDN, domainData.DomaineEntityProfileDn)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusCreated
}

func createACIDomain(addressPoolData *model.AddressPool, zoneName string) (interface{}, int, *capdata.ACIDomainData) {
	aciClient := caputilities.GetConnection()
	domainName := zoneName + "-DOM"
	physicalDomainAttributes := aciModels.PhysicalDomainAttributes{
		Name: domainName,
	}
	physDomResp, err := aciClient.CreatePhysicalDomain(domainName, "", physicalDomainAttributes)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil
	}
	// createVLANpool
	vlanPoolAttributes := aciModels.VLANPoolAttributes{
		Name:      domainName + "-VLAN",
		AllocMode: "static",
	}
	vlanPoolResp, err := aciClient.CreateVLANPool(vlanPoolAttributes.AllocMode, vlanPoolAttributes.Name, "", vlanPoolAttributes)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil
	}
	rangesAttribute := aciModels.RangesAttributes{
		From:      fmt.Sprintf("vlan-%d", addressPoolData.Ethernet.IPv4.VLANIdentifierAddressRange.Lower),
		To:        fmt.Sprintf("vlan-%d", addressPoolData.Ethernet.IPv4.VLANIdentifierAddressRange.Upper),
		AllocMode: vlanPoolAttributes.AllocMode,
	}
	_, err = aciClient.CreateRanges(rangesAttribute.To, rangesAttribute.From, rangesAttribute.AllocMode, vlanPoolAttributes.Name, "", rangesAttribute)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil
	}
	err = aciClient.CreateRelationinfraRsVlanNsFromPhysicalDomain(physDomResp.BaseAttributes.DistinguishedName, vlanPoolResp.BaseAttributes.DistinguishedName)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil
	}
	//CreateDomainEntityProfile for the given Domain
	entityProfileAttribute := aciModels.AttachableAccessEntityProfileAttributes{
		Name: domainName + "-EntityProfile",
	}
	entityProfileResp, err := aciClient.CreateAttachableAccessEntityProfile(entityProfileAttribute.Name, "", entityProfileAttribute)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil
	}
	err = aciClient.CreateRelationinfraRsDomPFromAttachableAccessEntityProfile(entityProfileResp.BaseAttributes.DistinguishedName, physDomResp.BaseAttributes.DistinguishedName)
	return nil, http.StatusCreated, &capdata.ACIDomainData{
		DomainDN:               physDomResp.BaseAttributes.DistinguishedName,
		DomaineEntityProfileDn: entityProfileResp.BaseAttributes.DistinguishedName,
	}
}

func saveZoneToDomainDNData(zoneID string, domainData *capdata.ACIDomainData) error {
	return capmodel.SaveZoneDomain(zoneID, domainData)
}

func getZoneTODomainDNData(zoneID string) (capdata.ACIDomainData, error) {
	return capmodel.GetZoneDomain(zoneID)
}

// UpdateZoneData provides patch operation on Zone
func UpdateZoneData(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	if _, err := capmodel.GetFabric(fabricID); err != nil {
		errMsg := fmt.Sprintf("failed to fetch fabric data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Fabric", fabricID})
		return
	}

	//TODO: Get list of zones which are pre-populated from onstart and compare the members for item not present in odim but present in ACI

	zoneData, err := capmodel.GetZone(fabricID, uri)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", uri})
		return
	}
	if zoneData.ZoneType != "ZoneOfEndpoints" {
		ctx.StatusCode(http.StatusMethodNotAllowed)
		resp := updateErrorResponse(response.ActionNotSupported, "", []interface{}{ctx.Request().Method})
		ctx.JSON(resp)
		return
	}
	var zoneRequest model.Zone
	if err = ctx.ReadJSON(&zoneRequest); err != nil {
		errorMessage := "error while trying to get JSON body from the  request: " + err.Error()
		log.Error(errorMessage)
		resp := updateErrorResponse(response.MalformedJSON, errorMessage, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}

	if zoneRequest.Links == nil {
		errMsg := fmt.Sprintf("Zone cannot be patched as there are Links is in the missing")
		log.Error(errMsg)
		resp := updateErrorResponse(response.PropertyMissing, errMsg, []interface{}{"Links"})
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	// get the AddressPoolData for the zone
	addresspoolData, statusCode, resp := getAddressPoolData(fabricID, zoneData.Links.AddressPools[0].Oid)
	if statusCode != http.StatusOK {
		ctx.StatusCode(statusCode)
		ctx.JSON(resp)
		return
	}
	// get the domaindata for the ZoneOfZone
	domainData, err := getZoneTODomainDNData(zoneData.Links.ContainedByZones[0].Oid)
	if err != nil {
		errMsg := fmt.Sprintf("Domain not found for  %s", zoneData.Links.ContainedByZones[0].Oid)
		log.Error(errMsg)
		resp = updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{zoneData.Links.ContainedByZones[0].Oid, "Domain"})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	// check all given endpoints
	// save all existing endpoints in the map
	endpointURIData := make(map[string]bool)
	endPointData := make(map[string]*capdata.EndpointData)
	for i := 0; i < len(zoneData.Links.Endpoints); i++ {
		endpointURIData[zoneData.Links.Endpoints[i].Oid] = false
		data, statusCode, resp := getEndpointData(fabricID, zoneData.Links.Endpoints[i].Oid)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		endPointData[zoneData.Links.Endpoints[i].Oid] = &data
	}
	endpointRequestData := make(map[string]*capdata.EndpointData)
	for i := 0; i < len(zoneRequest.Links.Endpoints); i++ {
		data, statusCode, resp := getEndpointData(fabricID, zoneRequest.Links.Endpoints[i].Oid)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		endpointRequestData[zoneRequest.Links.Endpoints[i].Oid] = &data
	}
	zoneofZoneURL := zoneData.Links.ContainedByZones[0].Oid
	// get the zone of zone data
	zoneofZoneData, err := capmodel.GetZone(fabricID, zoneofZoneURL)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", zoneofZoneURL, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", zoneofZoneURL})
		return
	}
	// Get the default zone data
	defaultZoneURL := zoneofZoneData.Links.ContainedByZones[0].Oid
	defaultZoneData, err := capmodel.GetZone(fabricID, defaultZoneURL)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", defaultZoneURL, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", defaultZoneURL})
		return
	}

	for endpointOID, data := range endpointRequestData {
		_, ok := endPointData[endpointOID]
		if !ok {
			resp, statusCode = createStaticPort(zoneData.Name+"-EPG", defaultZoneData.Name, zoneofZoneData.Name, data.ACIPolicyGroupData, addresspoolData.Ethernet.IPv4.VLANIdentifierAddressRange.Lower, &domainData)
			if statusCode != http.StatusCreated {
				ctx.StatusCode(statusCode)
				ctx.JSON(resp)
				return
			}
		}
		delete(endPointData, endpointOID)
	}

	for endpointOID, data := range endPointData {
		resp, statusCode = deleteStaticPort(data.ACIPolicyGroupData.PolicyGroupDN, zoneData.Name+"-EPG", defaultZoneData.Name, zoneofZoneData.Name)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		resp, statusCode = deleteRelationDomainEntityGroupInterfacePolicyGroup(data.ACIPolicyGroupData.PCVPCPolicyGroupDN)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		delete(endPointData, endpointOID)
	}
	zoneData.Links.Endpoints = zoneRequest.Links.Endpoints
	if err = updatezoneData(fabricID, uri, &zoneData); err != nil {
		errMsg := fmt.Sprintf("failed to update zone data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", uri})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(zoneData)
}

func deleteStaticPort(policyGroupDN, epgName, tenantName, applicationProfileName string) (interface{}, int) {
	aciClient := caputilities.GetConnection()
	err := aciClient.DeleteStaticPath(policyGroupDN, epgName, applicationProfileName, tenantName)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusOK
}

func deleteRelationDomainEntityGroupInterfacePolicyGroup(policyGroupDN string) (interface{}, int) {
	aciClient := caputilities.GetConnection()
	err := aciClient.DeleteRelationinfraRsAttEntPFromPCVPCInterfacePolicyGroup(policyGroupDN)
	if err != nil {
		errMsg := "Error while creating  Zone of Zones: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusOK
}

func updatezoneData(fabricID, zoneOID string, zoneData *model.Zone) error {
	return capmodel.UpdateZone(fabricID, zoneOID, zoneData)
}
