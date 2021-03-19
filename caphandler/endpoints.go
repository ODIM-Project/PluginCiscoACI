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
	iris "github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

//GetEndpointCollection : Fetches details of the given resource from the device
func GetEndpointCollection(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	_, ok := capdata.FabricDataStore.Data[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Endpoint data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Endpoint", uri})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	var members = []*model.Link{}

	for endpointID, endpointData := range capdata.EndpointDataStore {
		if endpointData.FabricID == fabricID {
			members = append(members, &model.Link{
				Oid: endpointID,
			})
		}
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

//CreateEndpoint : created endpoints under given fabric
func CreateEndpoint(ctx iris.Context) {
	// Add logic to check if given ports exits
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	fabricData, ok := capdata.FabricDataStore.Data[fabricID]
	if !ok {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}

	var endpoint model.Endpoint
	err := ctx.ReadJSON(&endpoint)
	if err != nil {
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
	for _, endpointData := range capdata.EndpointDataStore {
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
			errMsg := "Given Port already present in the request"
			resp := updateErrorResponse(response.PropertyValueConflict, errMsg, []interface{}{endpoint.Redundancy[0].RedundancySet[i].Oid, "RedundancySet"})
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(resp)
			return

		}
		portList[endpoint.Redundancy[0].RedundancySet[i].Oid] = true

		_, statusCode, resp := getPortData(portURI)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.JSON(resp)
			return
		}
		statusCode, resp = checkEndpointPortMapping(endpoint.Redundancy[0].RedundancySet[i].Oid)
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

	portPolicyGroupList, err := caputilities.GetPortPolicyGroup(fabricData.PodID, switchURI)
	if err != nil || len(portPolicyGroupList) == 0 {
		errMsg := "Port policy group not found for given ports"
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"protpaths" + switchURI, "PolicyGroup"})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return

	}
	policyGroupDN := ""
	for i := 0; i < len(portPolicyGroupList); i++ {
		if strings.Contains(portPolicyGroupList[i].BaseAttributes.DistinguishedName, portPattern) {
			policyGroupDN = portPolicyGroupList[i].BaseAttributes.DistinguishedName
		}
	}
	if policyGroupDN == "" {
		errMsg := "Port policy group not found for given ports"
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{portPattern, "PolicyGroup"})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	log.Info("Dn of Policy group:" + policyGroupDN)
	saveEndpointData(uri, fabricID, policyGroupDN, &endpoint)
	common.SetResponseHeader(ctx, map[string]string{
		"Location": endpoint.ODataID,
	})
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(endpoint)
}

//GetEndpointInfo : gets endpoints under given fabric
func GetEndpointInfo(ctx iris.Context) {
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

	respData, ok := capdata.EndpointDataStore[uri]
	if !ok {
		errMsg := fmt.Sprintf("Endpoint data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Endpoint", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(respData.Endpoint)

}

//DeleteEndpointInfo : deletes  endpoints under given fabric
func DeleteEndpointInfo(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	if _, ok := capdata.FabricDataStore.Data[fabricID]; !ok {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	endpointData, ok := capdata.EndpointDataStore[uri]
	if !ok {
		errMsg := fmt.Sprintf("Endpoint data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Endpoint", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
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
	delete(capdata.EndpointDataStore, uri)
	ctx.StatusCode(http.StatusNoContent)
}

func saveEndpointData(uri, fabricID, policyGroupDN string, endpoint *model.Endpoint) {
	endpointID := uuid.NewV4().String()
	endpoint.ID = endpointID
	endpoint.ODataContext = "/ODIM/v1/$metadata#Endpoint.Endpoint"
	endpoint.ODataType = "#Endpoint.v1_5_0.Endpoint"
	endpoint.ODataID = fmt.Sprintf("%s/%s", uri, endpointID)
	capdata.EndpointDataStore[endpoint.ODataID] = &capdata.EndpointData{
		FabricID:      fabricID,
		Endpoint:      endpoint,
		PolicyGroupDN: policyGroupDN,
	}

}

func checkEndpointPortMapping(portOID string) (int, interface{}) {
	// get all existing endpoints check if port is assinged to other endpoint
	for _, endpointData := range capdata.EndpointDataStore {
		for i := 0; i < len(endpointData.Endpoint.Redundancy[0].RedundancySet); i++ {
			if endpointData.Endpoint.Redundancy[0].RedundancySet[i].Oid == portOID {
				errMsg := "Port already assigned to other endpoint:" + portOID
				resp := updateErrorResponse(response.ResourceAlreadyExists, errMsg, []interface{}{"Endpoint", endpointData.Endpoint.Redundancy[0].RedundancySet[i].Oid, portOID})
				return http.StatusConflict, resp
			}
		}
	}
	return http.StatusOK, nil
}
