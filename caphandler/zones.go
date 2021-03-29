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
	"strings"
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
	_, ok := capdata.FabricDataStore.Data[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Address data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"AddressPool", uri})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	var members = []*model.Link{}

	for zoneID, zoneData := range capdata.ZoneDataStore {
		if zoneData.FabricID == fabricID {
			members = append(members, &model.Link{
				Oid: zoneID,
			})
		}
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
	ctx.JSON(respData.Zone)
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
		for _, value := range capdata.ZoneDataStore {
			if value.Zone.Name == zone.Name {
				conflictFlag = true
			}
		}
		if !conflictFlag {
			defaultZoneID = uuid.NewV4().String()
			zone = saveZoneData(defaultZoneID, uri, fabricID, zone)
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
		for _, value := range capdata.ZoneDataStore {
			if value.Zone.Name == zone.Name {
				conflictFlag = true
			}
		}
		if !conflictFlag {
			defaultZoneID = uuid.NewV4().String()
			zone = saveZoneData(defaultZoneID, uri, fabricID, zone)
			saveZoneToDomainDNData(zone.ODataID, domainData)
		}
		updateZoneData(defaultZoneLink, zone)
		updateAddressPoolData(zone.ODataID, zone.Links.AddressPools[0].Oid, "Add")
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
		zone = saveZoneData(zoneID, uri, fabricID, zone)
		updateZoneData(zoneofZoneOID, zone)
		updateAddressPoolData(zone.ODataID, zone.Links.AddressPools[0].Oid, "Add")
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

func saveZoneData(defaultZoneID string, uri string, fabricID string, zone model.Zone) model.Zone {
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
	capdata.ZoneDataStore[zone.ODataID] = &capdata.ZoneData{
		FabricID: fabricID,
		Zone:     &zone,
	}
	return zone
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
	_, ok := capdata.FabricDataStore.Data[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}

	//TODO: Get list of zones which are pre-populated from onstart and compare the members for item not present in odim but present in ACI

	respData, ok := capdata.ZoneDataStore[uri]
	log.Println(capdata.ZoneDataStore)
	if !ok {
		errMsg := fmt.Sprintf("Zone data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Zone", uri})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	if respData.Zone.Links != nil {
		if respData.Zone.Links.ContainsZonesCount != 0 {
			errMsg := fmt.Sprintf("Zone cannot be deleted as there are dependent resources still tied to it")
			log.Error(errMsg)
			resp := updateErrorResponse(response.ResourceCannotBeDeleted, errMsg, []interface{}{"Zone", uri})
			ctx.StatusCode(http.StatusNotAcceptable)
			ctx.JSON(resp)
			return
		}
	}
	if respData.Zone.ZoneType == "ZoneOfZones" {
		err := deleteZoneOfZone(respData, uri)
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
		delete(capdata.ZoneDataStore, uri)
		ctx.StatusCode(http.StatusNoContent)
	}
	if respData.Zone.ZoneType == "Default" {
		aciClient := caputilities.GetConnection()
		err := aciClient.DeleteTenant(respData.Zone.Name)
		if err != nil {
			errMsg := "Error while deleting Zone: " + err.Error()
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(resp)
			return
		}

		delete(capdata.ZoneDataStore, uri)
		ctx.StatusCode(http.StatusNoContent)
	}
	if respData.Zone.ZoneType == "ZoneOfEndpoints" {
		resp, statusCode := deleteZoneOfEndpoints(respData.Zone)
		ctx.StatusCode(statusCode)
		ctx.JSON(resp)
	}
}

func deleteZoneOfZone(respData *capdata.ZoneData, uri string) error {
	var parentZoneLink model.Link
	var parentZone *model.Zone
	if respData.Zone.Links != nil {
		if respData.Zone.Links.ContainedByZonesCount != 0 {
			// Assuming contained by link is only one
			parentZoneLink = respData.Zone.Links.ContainedByZones[0]
			parentZoneData, ok := capdata.ZoneDataStore[parentZoneLink.Oid]
			if !ok {
				errMsg := fmt.Errorf("Zone data for uri %s not found " + uri)
				return errMsg
			}
			parentZone = parentZoneData.Zone
			links := parentZone.Links.ContainsZones
			var parentZoneIndex int
			for index, value := range links {
				if value.Oid == uri {
					parentZoneIndex = index
					break
				}
			}
			parentZone.Links.ContainsZones = append(links[:parentZoneIndex], links[parentZoneIndex+1:]...)
			parentZone.Links.ContainsZonesCount = len(parentZone.Links.ContainsZones)
			parentZoneData.Zone = parentZone
			capdata.ZoneDataStore[parentZoneLink.Oid] = parentZoneData
		}
		aciServiceManager := caputilities.GetConnection()
		err := aciServiceManager.DeleteApplicationProfile(respData.Zone.Name, parentZone.Name)
		if err != nil {
			errMsg := fmt.Errorf("Error deleting Application Profile")
			return errMsg
		}
		vrfErr := aciServiceManager.DeleteVRF(respData.Zone.Name+"-VRF", parentZone.Name)
		if vrfErr != nil {
			errMsg := fmt.Errorf("Error deleting VRF")
			return errMsg
		}
		// delete contract
		contractErr := aciServiceManager.DeleteContract(respData.Zone.Name+"-VRF-Con", parentZone.Name)
		if contractErr != nil {
			errMsg := fmt.Errorf("Error deleting Contract:%v", contractErr)
			log.Error(errMsg.Error())
			return errMsg
		}
		err = aciServiceManager.DeleteAttachableAccessEntityProfile(respData.Zone.Name + "-DOM-EntityProfile")
		if err != nil {
			errMsg := fmt.Errorf("Error deleting  domain profile:%v", contractErr)
			log.Error(errMsg.Error())
			return errMsg
		}
		err = aciServiceManager.DeletePhysicalDomain(respData.Zone.Name + "-DOM")
		if err != nil {
			errMsg := fmt.Errorf("Error deleting Physical domain:%v", contractErr)
			log.Error(errMsg.Error())
			return errMsg
		}
		err = aciServiceManager.DeleteVLANPool("static", respData.Zone.Name+"-DOM-VLAN")
		if err != nil {
			errMsg := fmt.Errorf("Error deleting Physical domain:%v", contractErr)
			log.Error(errMsg.Error())
			return errMsg
		}
		updateAddressPoolData(respData.Zone.ODataID, respData.Zone.Links.AddressPools[0].Oid, "Remove")
		delete(capdata.ZoneDataStore, uri)
		delete(capdata.ZoneTODomainDN, uri)
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
	respData, ok := capdata.ZoneDataStore[defaultZoneLink]
	if !ok {
		errMsg := fmt.Sprintf("Zone data for uri %s not found", defaultZoneLink)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Zone", defaultZoneLink})
		return "", resp, http.StatusNotFound, nil
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

	addresspoolData, statusCode, resp := getAddressPoolData(zone.Links.AddressPools[0].Oid)
	if statusCode != http.StatusOK {
		return "", resp, statusCode, nil
	}
	if addresspoolData.Ethernet.IPv4.VLANIdentifierAddressRange == nil {
		errorMessage := "Provided AddressPool doesn't contain the VLANIdentifierAddressRange"
		return "", updateErrorResponse(response.PropertyMissing, errorMessage, []interface{}{"VLANIdentifierAddressRange"}), http.StatusBadRequest, nil
	}
	aciClient := caputilities.GetConnection()
	appProfileList, err := aciClient.ListApplicationProfile(respData.Zone.Name)
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
	vrfList, err := aciClient.ListVRF(respData.Zone.Name)
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

	apResp, err := CreateApplicationProfile(zone.Name, respData.Zone.Name, respData.Zone.Description, apModel)
	if err != nil {
		errMsg := "Error while creating application profile: " + err.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return "", resp, http.StatusBadRequest, nil
	}
	_, vrfErr := CreateVRF(vrfModel.Name, respData.Zone.Name, respData.Zone.Description, vrfModel)
	if vrfErr != nil {
		errMsg := "Error while creating application profile: " + vrfErr.Error()
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return "", resp, http.StatusBadRequest, nil
	}
	// create contract with name vrf and suffix-Con
	resp, statusCode = createContract(vrfModel.Name, respData.Zone.Name, zone.Name)
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

func updateZoneData(defaultZoneLink string, zone model.Zone) {
	defaultZoneStore := capdata.ZoneDataStore[defaultZoneLink]
	defaultZoneData := defaultZoneStore.Zone
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

	capdata.ZoneDataStore[defaultZoneLink].Zone = defaultZoneData
	return
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
	zoneofZoneData, ok := capdata.ZoneDataStore[zoneofZoneURL]
	if !ok {
		errMsg := fmt.Sprintf("ZoneofZone data for uri %s not found", uri)
		log.Error(errMsg)
		return "", updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"ZoneofZone", zoneofZoneURL}), http.StatusNotFound
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

	addresspoolData, statusCode, resp := getAddressPoolData(zone.Links.AddressPools[0].Oid)
	if statusCode != http.StatusOK {
		return "", resp, statusCode
	}

	if addresspoolData.Links != nil && len(addresspoolData.Links.Zones) > 0 {
		errorMessage := fmt.Sprintf("Given AddressPool %s is assingned to other ZoneofEndpoints", zone.Links.AddressPools[0].Oid)
		return "", updateErrorResponse(response.ResourceInUse, errorMessage, []interface{}{"AddressPools", "AddressPools"}), http.StatusBadRequest
	}
	// Get the endpoints from the db
	// validate all given addresspools if it's present
	if len(zone.Links.Endpoints) == 0 {
		errorMessage := "Endpoints attribute is missing in the request"
		return "", updateErrorResponse(response.PropertyMissing, errorMessage, []interface{}{"Endpoints"}), http.StatusBadRequest
	}
	if len(zone.Links.Endpoints) > 1 {
		errorMessage := "More than one Endpoints not allowed for the creation of ZoneOfEndpoints"
		return "", updateErrorResponse(response.PropertyValueFormatError, errorMessage, []interface{}{"Endpoints", "Endpoints"}), http.StatusBadRequest
	}
	// Get the default zone data
	defaultZoneURL := zoneofZoneData.Zone.Links.ContainedByZones[0].Oid
	defaultZoneData := capdata.ZoneDataStore[defaultZoneURL]
	endPointData := make(map[string]*capdata.EndpointData)
	for i := 0; i < len(zone.Links.Endpoints); i++ {
		data, statusCode, resp := getEndpointData(zone.Links.Endpoints[i].Oid)
		if statusCode != http.StatusOK {
			return "", resp, statusCode
		}
		endPointData[zone.Links.Endpoints[i].Oid] = data
	}
	// get domain from given addresspool native vlan from config
	domainData, ok := getZoneTODomainDNData(zoneofZoneURL)
	if !ok {
		errMsg := fmt.Sprintf("Domain not found for  %s", zoneofZoneURL)
		log.Error(errMsg)
		return "", updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{zoneofZoneURL, "Domain"}), http.StatusNotFound
	}
	bdResp, bdDN, statusCode := createBridgeDomain(defaultZoneData.Zone.Name, zone)
	if statusCode != http.StatusCreated {
		return "", bdResp, statusCode
	}

	// create the subnet for BD for all given address pool
	resp, statusCode = createSubnets(defaultZoneData.Zone.Name, zone.Name, addresspoolData)
	if statusCode != http.StatusCreated {
		return "", resp, statusCode
	}
	// link bridgedomain to vrf
	resp, statusCode = linkBDtoVRF(bdDN, zoneofZoneData.Zone.Name+"-VRF")
	if statusCode != http.StatusCreated {
		return "", resp, statusCode
	}
	resp, statusCode = applicationEPGOperation(defaultZoneData.Zone.Name, zoneofZoneData.Zone.Name, zone.Name, domainData, endPointData, addresspoolData.Ethernet.IPv4.NativeVLAN)
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

func applicationEPGOperation(tenantName, applicationProfileName, bdName string, domainData *capdata.ACIDomainData, endPointData map[string]*capdata.EndpointData, nativeVLAN int) (interface{}, int) {
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
	// Create static port'
	for _, data := range endPointData {
		resp, statusCode = createStaticPort(epgName, tenantName, applicationProfileName, data.ACIPolicyGroupData, nativeVLAN, domainData)
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

func deleteZoneOfEndpoints(zoneData *model.Zone) (interface{}, int) {
	zoneofZoneURL := zoneData.Links.ContainedByZones[0].Oid
	// get the zone of zone data
	zoneofZoneData := capdata.ZoneDataStore[zoneofZoneURL].Zone
	// Get the default zone data
	defaultZoneURL := zoneofZoneData.Links.ContainedByZones[0].Oid
	defaultZoneData := capdata.ZoneDataStore[defaultZoneURL].Zone
	aciClient := caputilities.GetConnection()
	for i := 0; i < len(zoneData.Links.Endpoints); i++ {
		endpointData, statusCode, resp := getEndpointData(zoneData.Links.Endpoints[i].Oid)
		if statusCode != http.StatusOK {
			return resp, statusCode
		}
		resp, statusCode = deleteRelationDomainEntityGroupInterfacePolicyGroup(endpointData.ACIPolicyGroupData.PCVPCPolicyGroupDN)
		if statusCode != http.StatusOK {
			return resp, statusCode
		}
	}
	err := aciClient.DeleteApplicationEPG(zoneData.Name+"-EPG", zoneofZoneData.Name, defaultZoneData.Name)
	if err != nil {
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
		capdata.ZoneDataStore[zoneofZoneURL].Zone = zoneofZoneData
	}
	updateAddressPoolData(zoneData.ODataID, zoneData.Links.AddressPools[0].Oid, "Remove")
	delete(capdata.ZoneDataStore, zoneData.ODataID)
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

func saveZoneToDomainDNData(zoneID string, domainData *capdata.ACIDomainData) {
	capdata.ZoneTODomainDN[zoneID] = domainData
}

func getZoneTODomainDNData(zoneID string) (*capdata.ACIDomainData, bool) {
	data, ok := capdata.ZoneTODomainDN[zoneID]
	return data, ok
}
// UpdateZoneData provides patch operation on Zone 
func UpdateZoneData(ctx iris.Context) {
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

	//TODO: Get list of zones which are pre-populated from onstart and compare the members for item not present in odim but present in ACI

	zoneData, ok := capdata.ZoneDataStore[uri]
	if !ok {
		errMsg := fmt.Sprintf("Zone data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Zone", uri})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	if zoneData.Zone.ZoneType != "ZoneOfEndpoints" {
		ctx.StatusCode(http.StatusMethodNotAllowed)
		resp := updateErrorResponse(response.ActionNotSupported, "", []interface{}{ctx.Request().Method})
		ctx.JSON(resp)
		return
	}
	var zoneRequest model.Zone
	err := ctx.ReadJSON(&zoneRequest)
	if err != nil {
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
	addresspoolData, statusCode, resp := getAddressPoolData(zoneData.Zone.Links.AddressPools[0].Oid)
	if statusCode != http.StatusOK {
		ctx.StatusCode(statusCode)
		ctx.JSON(resp)
		return
	}
	// get the domaindata for the ZoneOfZone
	domainData, ok := getZoneTODomainDNData(zoneData.Zone.Links.ContainedByZones[0].Oid)
	if !ok {
		errMsg := fmt.Sprintf("Domain not found for  %s", zoneData.Zone.Links.ContainedByZones[0].Oid)
		log.Error(errMsg)
		resp = updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{zoneData.Zone.Links.ContainedByZones[0].Oid, "Domain"})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	// check all given endpoints
	// save all existing endpoints in the map
	endpointURIData := make(map[string]bool)
	endPointData := make(map[string]*capdata.EndpointData)
	for i := 0; i < len(zoneData.Zone.Links.Endpoints); i++ {
		endpointURIData[zoneData.Zone.Links.Endpoints[i].Oid] = false
		data, statusCode, resp := getEndpointData(zoneData.Zone.Links.Endpoints[i].Oid)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		endPointData[zoneData.Zone.Links.Endpoints[i].Oid] = data
	}
	endpointRequestData := make(map[string]*capdata.EndpointData)
	for i := 0; i < len(zoneRequest.Links.Endpoints); i++ {
		data, statusCode, resp := getEndpointData(zoneRequest.Links.Endpoints[i].Oid)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		endpointRequestData[zoneRequest.Links.Endpoints[i].Oid] = data
	}
	zoneofZoneURL := zoneData.Zone.Links.ContainedByZones[0].Oid
	// get the zone of zone data
	zoneofZoneData := capdata.ZoneDataStore[zoneofZoneURL]
	// get default zone data
	defaultZoneURL := zoneofZoneData.Zone.Links.ContainedByZones[0].Oid
	defaultZoneData := capdata.ZoneDataStore[defaultZoneURL]
	for endpointOID, data := range endpointRequestData {
		_, ok := endPointData[endpointOID]
		if !ok {
			resp, statusCode = createStaticPort(zoneData.Zone.Name+"-EPG", defaultZoneData.Zone.Name, zoneofZoneData.Zone.Name, data.ACIPolicyGroupData, addresspoolData.Ethernet.IPv4.NativeVLAN, domainData)
			if statusCode == http.StatusCreated {
				ctx.StatusCode(statusCode)
				ctx.JSON(resp)
				return
			}
		}
		delete(endPointData, endpointOID)
	}

	for endpointOID, data := range endPointData {
		resp, statusCode = deleteStaticPort(data.ACIPolicyGroupData.PCVPCPolicyGroupDN, zoneData.Zone.Name+"-EPG", defaultZoneData.Zone.Name, zoneofZoneData.Zone.Name)
		if statusCode == http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		resp, statusCode = deleteRelationDomainEntityGroupInterfacePolicyGroup(data.ACIPolicyGroupData.PCVPCPolicyGroupDN)
		if statusCode == http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		delete(endPointData, endpointOID)
	}
	zoneData.Zone.Links.Endpoints = zoneRequest.Links.Endpoints
	updatezoneData(uri, zoneData)
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(zoneData.Zone)
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

func updatezoneData(zoneOID string, zoneData *capdata.ZoneData) {
	capdata.ZoneDataStore[zoneOID] = zoneData
}
