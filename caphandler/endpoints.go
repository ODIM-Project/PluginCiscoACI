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

// Package caphandler ...
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

// GetEndpointCollection : Fetches details of the given resource from the device
func GetEndpointCollection(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")

	endpointData, err := capmodel.GetAllEndpoints(fabricID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch endpoint data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Endpoint", uri})
		return
	}

	var members = []*model.Link{}
	for endpointID := range endpointData {
		members = append(members, &model.Link{
			Oid: endpointID,
		})
	}

	endpointCollection := model.Collection{
		ODataContext: "/ODIM/v1/$metadata#EndpointCollection.EndpointCollection",
		ODataID:      uri,
		ODataType:    "#EndpointCollection.EndpointCollection",
		Description:  "EndpointCollection view",
		Name:         "Endpoints",
		Members:      members,
		MembersCount: len(members),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(endpointCollection)
}

// CreateEndpoint : created endpoints under given fabric
func CreateEndpoint(ctx iris.Context) {
	// Add logic to check if given ports exits
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	fabricData, err := capmodel.GetFabric(fabricID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch fabric data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Fabric", fabricID})
		return
	}

	var endpoint model.Endpoint
	if err = ctx.ReadJSON(&endpoint); err != nil {
		errorMessage := "error while trying to get JSON body from the  request: " + err.Error()
		log.Error(errorMessage)
		resp := updateErrorResponse(response.MalformedJSON, errorMessage, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	if len(endpoint.Redundancy) < 1 {
		errMsg := fmt.Sprintf("Endpoint cannot be created, Redudancy in the request is missing: " + err.Error())
		resp := updateErrorResponse(response.PropertyMissing, errMsg, []interface{}{"Redundancy"})
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	if len(endpoint.Redundancy[0].RedundancySet) == 0 {
		errMsg := fmt.Sprintf("Endpoint cannot be created, RedudancySet in the request is missing: " + err.Error())
		resp := updateErrorResponse(response.PropertyMissing, errMsg, []interface{}{"RedudancySet"})
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	// get all existing endpoints under fabric check for the name
	data, err := capmodel.GetAllEndpoints(fabricID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch endpoint data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Fabric", fabricID})
		return
	}
	for _, endpointData := range data {
		if endpoint.Name == endpointData.Endpoint.Name {
			errMsg := "Endpoint name is already assigned to other endpoint:" + endpointData.Endpoint.Name
			resp := updateErrorResponse(response.ResourceAlreadyExists, errMsg, []interface{}{"Endpoint", endpointData.Endpoint.Name, endpoint.Name})
			ctx.StatusCode(http.StatusConflict)
			ctx.JSON(resp)
			return
		}
	}
	var switchURI = ""
	var portPattern = ""
	portList := make(map[string]bool)
	// check if given ports are present in plugin database
	for i := 0; i < len(endpoint.Redundancy[0].RedundancySet); i++ {
		portURI := endpoint.Redundancy[0].RedundancySet[i].Oid
		if _, ok := portList[endpoint.Redundancy[0].RedundancySet[i].Oid]; ok {
			errMsg := "Duplicate port passed in the request"
			resp := updateErrorResponse(response.PropertyValueConflict, errMsg, []interface{}{endpoint.Redundancy[0].RedundancySet[i].Oid, endpoint.Redundancy[0].RedundancySet[i].Oid})
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(resp)
			return

		}
		portList[endpoint.Redundancy[0].RedundancySet[i].Oid] = true

		portData := getPortData(ctx, portURI)
		if portData == nil {
			return
		}
		statusCode, resp := checkEndpointPortMapping(endpoint.Redundancy[0].RedundancySet[i].Oid)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		portURIData := strings.Split(portURI, "/")
		switchID := portURIData[6]
		switchIDData := strings.Split(switchID, ":")
		switchURI = switchURI + "-" + switchIDData[1]
		portIDData := strings.Split(portURIData[8], ":")
		tmpPortPattern := strings.Replace(portIDData[1], "eth", "", -1)
		tmpPortPattern = strings.Replace(tmpPortPattern, "-", "-ports-", -1)
		portPattern = tmpPortPattern
	}
	policyGroupDN := ""

	// create policyGroup for the given ports
	resp, statusCode, aciPolicyGroupData := createPolicyGroup(switchURI, portPattern)
	if statusCode != http.StatusCreated {
		ctx.StatusCode(statusCode)
		ctx.JSON(resp)
		return
	}

	log.Info("Dn of Policy group:" + policyGroupDN)
	aciPolicyGroupData.PolicyGroupDN = fmt.Sprintf("topology/pod-%s/protpaths%s/pathep-[%s]", fabricData.PodID, switchURI, aciPolicyGroupData.PcVPCPolicyGroupName)
	if err = saveEndpointData(uri, fabricID, aciPolicyGroupData, &endpoint); err != nil {
		errMsg := fmt.Sprintf("failed to store endpoint data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Fabric", fabricID})
		return
	}
	common.SetResponseHeader(ctx, map[string]string{
		"Location": endpoint.ODataID,
	})
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(endpoint)
}

// GetEndpointInfo : gets endpoints under given fabric
func GetEndpointInfo(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")

	endpointData, err := capmodel.GetEndpoints(fabricID, uri)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch endpoint data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Endpoint", fabricID})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(endpointData.Endpoint)
}

// DeleteEndpointInfo : deletes  endpoints under given fabric
func DeleteEndpointInfo(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")

	endpointData, err := capmodel.GetEndpoints(fabricID, uri)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch endpoint data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Endpoint", fabricID})
		return
	}

	if endpointData.Endpoint.Links != nil && len(endpointData.Endpoint.Links.AddressPools) > 0 {
		errMsg := fmt.Sprintf("Endpoint cannot be deleted as there are dependent upon AddressPool")
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceCannotBeDeleted, errMsg, []interface{}{uri, "Endpoint"})
		ctx.StatusCode(http.StatusNotAcceptable)
		ctx.JSON(resp)
		return
	}
	// Todo:Add the validation  to verify the links
	resp, statusCode := deletePolicyGroup(endpointData.ACIPolicyGroupData)
	if statusCode != http.StatusOK {
		ctx.JSON(resp)
		ctx.StatusCode(statusCode)
		return
	}
	// check if endpoint is associated with any ZoneOfEndpoints
	// get all zones
	zoneCollectionData, err := capmodel.GetAllZones(fabricID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch zone data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", fabricID})
		return

	}
	for zoneURI, zoneData := range zoneCollectionData {
		if zoneData.ZoneType == "ZoneOfEndpoints" {
			updateRequired := false
			endpointsList := zoneData.Links.Endpoints
			for i := 0; i < len(endpointsList); i++ {
				if endpointsList[i].Oid == uri {
					endpointsList = append(endpointsList[:i], endpointsList[i+1:]...)
					updateRequired = true
					break
				}
			}
			if updateRequired {
				zoneData.Links.Endpoints = endpointsList
				if err = capmodel.UpdateZone(fabricID, zoneURI, &zoneData); err != nil {
					errMsg := fmt.Sprintf("failed to update zone data for %s: %s", zoneURI, err.Error())
					createDbErrResp(ctx, err, errMsg, []interface{}{"Zone", fabricID})
					return
				}
			}

		}
	}

	if err := capmodel.DeleteEndpoint(fabricID, uri); err != nil {
		errMsg := fmt.Sprintf("failed to delete endpoint data for uri %s: %s", uri, err.Error())
		createDbErrResp(ctx, err, errMsg, []interface{}{"Endpoint", fabricID})
		return
	}
	ctx.StatusCode(http.StatusNoContent)
}

func saveEndpointData(uri, fabricID string, aciPolicyGroupData *capdata.ACIPolicyGroupData, endpoint *model.Endpoint) error {
	endpointID := uuid.NewV4().String()
	endpoint.ID = endpointID
	endpoint.ODataContext = "/ODIM/v1/$metadata#Endpoint.Endpoint"
	endpoint.ODataType = "#Endpoint.v1_5_0.Endpoint"
	endpoint.ODataID = fmt.Sprintf("%s/%s", uri, endpointID)
	data := &capdata.EndpointData{
		Endpoint:           endpoint,
		ACIPolicyGroupData: aciPolicyGroupData,
	}
	if err := capmodel.SaveEndpoint(fabricID, endpoint.ODataID, data); err != nil {
		return err
	}
	return nil
}

func checkEndpointPortMapping(portOID string) (int, interface{}) {
	fabricData, err := capmodel.GetAllFabric("")
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch fabric data: %s", err.Error())
		return createDbErrResp(nil, err, errMsg, nil)
	}

	for fabricID := range fabricData {
		// get all existing endpoints check if port is assinged to other endpoint
		data, err := capmodel.GetAllEndpoints(fabricID)
		if err != nil {
			errMsg := fmt.Sprintf("failed to fetch endpoint data belonging to fabric %s: %s", fabricID, err.Error())
			return createDbErrResp(nil, err, errMsg, []interface{}{"Fabric", fabricID})
		}

		for _, endpointData := range data {
			for i := 0; i < len(endpointData.Endpoint.Redundancy[0].RedundancySet); i++ {
				if endpointData.Endpoint.Redundancy[0].RedundancySet[i].Oid == portOID {
					errMsg := "Port already assigned to other endpoint:" + portOID
					resp := updateErrorResponse(response.ResourceAlreadyExists, errMsg, []interface{}{"Endpoint", endpointData.Endpoint.Redundancy[0].RedundancySet[i].Oid, portOID})
					return http.StatusConflict, resp
				}
			}
		}
	}
	return http.StatusOK, nil
}

func getEndpointData(fabricID, endpoinOID string) (capdata.EndpointData, int, interface{}) {
	respData, err := capmodel.GetEndpoints(fabricID, endpoinOID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to fetch endpoint data for %s:%s: %s", fabricID, endpoinOID, err.Error())
		statuscode, resp := createDbErrResp(nil, err, errMsg, []interface{}{"Endpoint", endpoinOID})
		return respData, statuscode, resp
	}
	return respData, http.StatusOK, nil
}

func createPolicyGroup(switchPattern, portPattern string) (interface{}, int, *capdata.ACIPolicyGroupData) {
	// check if switch profile is present
	aciClient := caputilities.GetConnection()
	var err error

	switchProfileSelectorName := "Switch" + switchPattern + "_Profile_ifselector"
	accessPortSelectorName := "Switch" + switchPattern + "_" + portPattern

	var switchInterfaceProfileResp *aciModels.LeafInterfaceProfile
	portPatternData := strings.Split(portPattern, "-ports-")
	switchInterfaceProfileResp, err = aciClient.ReadLeafInterfaceProfile(switchProfileSelectorName)
	if err != nil {
		if !strings.Contains(err.Error(), "Object may not exists") {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
		// switch profile is not found creating the switch profile
		leafInterfaceAttributes := aciModels.LeafInterfaceProfileAttributes{
			Name: switchProfileSelectorName,
		}
		switchInterfaceProfileResp, err = aciClient.CreateLeafInterfaceProfile(switchProfileSelectorName, "", leafInterfaceAttributes)
		if err != nil {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
	}
	// create access port seletor
	accessPortSelectorAttributes := aciModels.AccessPortSelectorAttributes{
		Name:                    accessPortSelectorName,
		AccessPortSelector_type: "range",
	}
	accessPortSelectorResp, err := aciClient.CreateAccessPortSelector(accessPortSelectorAttributes.AccessPortSelector_type, accessPortSelectorName, switchProfileSelectorName, "", accessPortSelectorAttributes)
	if err != nil {
		errMsg := "Error while creating Endpoint: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil
	}
	portBlockName := "block-" + portPatternData[1]
	portBlockAttributes := aciModels.AccessPortBlockAttributes{
		Name:     portBlockName,
		FromPort: portPatternData[1],
		ToPort:   portPatternData[1],
	}
	_, err = aciClient.CreateAccessPortBlock(portBlockName, accessPortSelectorAttributes.AccessPortSelector_type, accessPortSelectorName, switchProfileSelectorName, "", portBlockAttributes)
	if err != nil {
		errMsg := "Error while creating Endpoint: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil
	}
	// check if vpc port policy is created with name ODIM-PORT-VPCPolicy
	portVPCPolicyName := "ODIM-PORT-VPCPolicy"

	_, err = aciClient.ReadLACPPolicy(portVPCPolicyName)
	if err != nil {
		if !strings.Contains(err.Error(), "Object may not exists") {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
		// switch profile is not found creating the switch profile
		lacpPolicyAttributes := aciModels.LACPPolicyAttributes{
			Name: portVPCPolicyName,
			Mode: "active",
		}
		_, err = aciClient.CreateLACPPolicy(portVPCPolicyName, "", lacpPolicyAttributes)
		if err != nil {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
	}
	// createPCVPC interface policy group
	pcVPCPolicyGroupName := "Switch" + switchPattern + "_" + portPattern + "_PolGrp"
	var pcVPCPolicyGroupAtrributes = aciModels.PCVPCInterfacePolicyGroupAttributes{
		Name: pcVPCPolicyGroupName,
		LagT: "node",
	}
	pcVPCPolicyGroupResp, err := aciClient.CreatePCVPCInterfacePolicyGroup(pcVPCPolicyGroupName, "", pcVPCPolicyGroupAtrributes)
	if err != nil {
		errMsg := "Error while creating Endpoint: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil
	}
	log.Info("Attaching policy group to port selector")
	err = aciClient.CreateRelationinfraRsAccBaseGrpFromAccessPortSelector(accessPortSelectorResp.BaseAttributes.DistinguishedName, pcVPCPolicyGroupResp.BaseAttributes.DistinguishedName)
	if err != nil {
		errMsg := "Error while creating Endpoint: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil

	}
	err = aciClient.CreateRelationinfraRsLacpPolFromPCVPCInterfacePolicyGroup(pcVPCPolicyGroupResp.BaseAttributes.DistinguishedName, portVPCPolicyName)
	if err != nil {
		errMsg := "Error while creating Endpoint: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest, nil

	}
	// if leaf profile is created else create the same
	var switchProfileName = "Switch" + switchPattern + "_Profile"
	switchPatternData := strings.Split(switchPattern, "-")
	var switchProfileResp *aciModels.LeafProfile
	switchProfileResp, err = aciClient.ReadLeafProfile(switchProfileName)
	if err != nil {
		if !strings.Contains(err.Error(), "Object may not exists") {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
		// switch profile is not found creating the switch profile
		leafprofileAttributes := aciModels.LeafProfileAttributes{
			Name: switchProfileName,
		}
		switchProfileResp, err = aciClient.CreateLeafProfile(switchProfileName, "", leafprofileAttributes)
		if err != nil {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
	}
	// check if switch assoication exist for given switch profile
	switchAssociationName := switchProfileName + "selector_"
	for i := 0; i < len(switchPatternData); i++ {
		switchAssociationName = switchAssociationName + switchPatternData[i]
	}
	_, err = aciClient.ReadSwitchAssociation("range", switchAssociationName, switchProfileName)
	if err != nil {
		if !strings.Contains(err.Error(), "Object may not exists") {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
		// switch profile is not found creating the switch profile
		switchAssociationAttributes := aciModels.SwitchAssociationAttributes{
			Name:                    switchAssociationName,
			Switch_association_type: "range",
		}
		_, err = aciClient.CreateSwitchAssociation("range", switchAssociationName, switchProfileName, "", switchAssociationAttributes)
		if err != nil {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
	}

	for i := 0; i < len(switchPatternData); i++ {
		//createNodeBlock for all switches
		resp, statusCode := createNodeBlock(switchProfileName, switchAssociationName, switchPatternData[i], i)
		if statusCode != http.StatusCreated {
			return resp, statusCode, nil
		}
	}

	// check if switch profile is associated with the switch interface profile
	_, err = aciClient.ReadRelationinfraRsAccPortPFromLeafProfile(switchProfileResp.BaseAttributes.DistinguishedName)
	if err != nil {
		if !strings.Contains(err.Error(), "Object may not exists") {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil
		}
		// associate switch profile with the switch interface profile
		err = aciClient.CreateRelationinfraRsAccPortPFromLeafProfile(switchProfileResp.BaseAttributes.DistinguishedName, switchInterfaceProfileResp.BaseAttributes.DistinguishedName)
		if err != nil {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest, nil

		}

	}

	aciPolicyGroupData := capdata.ACIPolicyGroupData{
		SwitchProfileName:         switchProfileName,
		SwitchAssociationName:     switchAssociationName,
		SwitchProfileSelectorName: switchProfileSelectorName,
		AccessPortSelectorName:    accessPortSelectorName,
		PcVPCPolicyGroupName:      pcVPCPolicyGroupName,
		PCVPCPolicyGroupDN:        pcVPCPolicyGroupResp.BaseAttributes.DistinguishedName,
	}
	return nil, http.StatusCreated, &aciPolicyGroupData
}

func createNodeBlock(switchProfileName, switchAssociationName, switchID string, index int) (interface{}, int) {
	// check if node block exist for given switch
	nodeblockName := fmt.Sprintf("single-%d", index)
	aciClient := caputilities.GetConnection()

	_, err := aciClient.ReadNodeBlock(nodeblockName, "range", switchAssociationName, switchProfileName)
	if err != nil {
		if !strings.Contains(err.Error(), "Object may not exists") {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest
		}
		// switch profile is not found creating the switch profile
		nodeBlockAttributes := aciModels.NodeBlockAttributes{
			Name:  nodeblockName,
			From_: switchID,
			To_:   switchID,
		}
		_, err = aciClient.CreateNodeBlock(nodeblockName, "range", switchAssociationName, switchProfileName, "", nodeBlockAttributes)
		if err != nil {
			errMsg := "Error while creating Endpoint: " + err.Error()
			log.Error(errMsg)
			resp := updateErrorResponse(response.GeneralError, errMsg, nil)
			return resp, http.StatusBadRequest
		}
	}
	return nil, http.StatusCreated
}

func deletePolicyGroup(aciPolicyGroupData *capdata.ACIPolicyGroupData) (interface{}, int) {
	aciClient := caputilities.GetConnection()

	err := aciClient.DeleteAccessPortSelector("range", aciPolicyGroupData.AccessPortSelectorName, aciPolicyGroupData.SwitchProfileSelectorName)
	if err != nil {
		errMsg := "Error while deleting Endpoint: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	err = aciClient.DeletePCVPCInterfacePolicyGroup(aciPolicyGroupData.PcVPCPolicyGroupName)
	if err != nil {
		errMsg := "Error while deleting  Endpoint: " + err.Error()
		log.Error(errMsg)
		resp := updateErrorResponse(response.GeneralError, errMsg, nil)
		return resp, http.StatusBadRequest
	}
	return nil, http.StatusOK
}
