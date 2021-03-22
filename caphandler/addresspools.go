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
	"fmt"
	"net"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/db"

	iris "github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// GetAddressPoolCollection fetches the addresspool which are linked to that fabric
func GetAddressPoolCollection(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	// get all switches which are store under that fabric

	if _, err := capmodel.GetFabric(fabricID); errors.Is(err, db.ErrorKeyNotFound) {
		errMsg := fmt.Sprintf("Address data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"AddressPool", uri})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	var members = []*model.Link{}

	for addressPoolOID, addressPooldData := range capdata.AddressPoolDataStore {
		if addressPooldData.FabricID == fabricID {
			members = append(members, &model.Link{
				Oid: addressPoolOID,
			})
		}
	}

	addressPoolCollectionResponse := model.Collection{
		ODataContext: "/ODIM/v1/$metadata#AddressPoolCollection.AddressPoolCollection",
		ODataID:      uri,
		ODataType:    "#AddressPoolCollection.AddressPoolCollection",
		Description:  "AddressPool view",
		Name:         "AddressPools",
		Members:      members,
		MembersCount: len(members),
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(addressPoolCollectionResponse)
}

// GetAddressPoolInfo fetches the addresspool info for given addresspool id
func GetAddressPoolInfo(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	if _, err := capmodel.GetFabric(fabricID); errors.Is(err, db.ErrorKeyNotFound) {
		errMsg := fmt.Sprintf("AddressPool data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	// Get the addresspool data from the memory
	addressPoolResponse, ok := capdata.AddressPoolDataStore[uri]
	if !ok {
		errMsg := fmt.Sprintf("AddressPool data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"AddressPool", uri})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(addressPoolResponse.AddressPool)
}

// CreateAddressPool stores the given addresspool against given fabric
func CreateAddressPool(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	if _, err := capmodel.GetFabric(fabricID); errors.Is(err, db.ErrorKeyNotFound) {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	var addresspoolData model.AddressPool
	err := ctx.ReadJSON(&addresspoolData)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the  request: " + err.Error()
		log.Error(errorMessage)
		resp := updateErrorResponse(response.MalformedJSON, errorMessage, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	// Todo :Add required validation for the request params
	missingAttribute, err := validateAddressPoolRequest(addresspoolData)
	if err != nil {
		log.Error(err.Error())
		resp := updateErrorResponse(response.PropertyMissing, err.Error(), []interface{}{missingAttribute})
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	// validate cidr given in request
	if _, _, err := net.ParseCIDR(addresspoolData.Ethernet.IPv4.GatewayIPAddress); err != nil {
		errorMessage := "Invalid value for GatewayIPAddress:" + err.Error()
		log.Errorf(errorMessage)
		resp := updateErrorResponse(response.PropertyValueFormatError, errorMessage, []interface{}{addresspoolData.Ethernet.IPv4.GatewayIPAddress, "GatewayIPAddress"})
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return

	}
	if addresspoolData.Ethernet.IPv4.NativeVLAN < 2 ||
		(addresspoolData.Ethernet.IPv4.NativeVLAN > 1001 && addresspoolData.Ethernet.IPv4.NativeVLAN < 1006) ||
		addresspoolData.Ethernet.IPv4.NativeVLAN > 4094 {
		errorMessage := "Invalid value for NativeVLAN: it should in range of 2 to 1001 or 1006 to 4094"
		log.Errorf(errorMessage)
		resp := updateErrorResponse(response.PropertyValueNotInList, errorMessage, []interface{}{fmt.Sprintf("%d", addresspoolData.Ethernet.IPv4.NativeVLAN), "NativeVLAN"})
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(resp)
		return
	}
	for _, data := range capdata.AddressPoolDataStore {
		if data.AddressPool.Ethernet.IPv4.GatewayIPAddress == addresspoolData.Ethernet.IPv4.GatewayIPAddress {
			errorMessage := "Requested GatewayIPAddress is already present in the addresspool " + data.AddressPool.ODataID
			log.Error(errorMessage)
			resp := updateErrorResponse(response.ResourceAlreadyExists, errorMessage, []interface{}{"AddressPool", "GatewayIPAddress", addresspoolData.Ethernet.IPv4.GatewayIPAddress})
			ctx.StatusCode(http.StatusConflict)
			ctx.JSON(resp)
			return
		}
	}
	addressPoolID := uuid.NewV4().String()
	addresspoolData.ODataContext = "/ODIM/v1/$metadata#AddressPool.AddressPool"
	addresspoolData.ODataType = "#AddressPool.v1_1_0.AddressPool"
	addresspoolData.ODataID = fmt.Sprintf("%s/%s", uri, addressPoolID)
	addresspoolData.ID = addressPoolID

	capdata.AddressPoolDataStore[addresspoolData.ODataID] = &capdata.AddressPoolsData{
		FabricID:    fabricID,
		AddressPool: &addresspoolData,
	}
	common.SetResponseHeader(ctx, map[string]string{
		"Location": addresspoolData.ODataID,
	})
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(addresspoolData)
}

func validateAddressPoolRequest(request model.AddressPool) (string, error) {
	if request.Ethernet == nil {
		return "Ethernet", fmt.Errorf("Ethernet data in request is missing")
	}
	if request.Ethernet.IPv4 == nil {
		return "IPv4", fmt.Errorf("Ethernet IPV4 data  in request is missing")
	}
	if request.Ethernet.IPv4.GatewayIPAddress == "" {
		return "GatewayIPAddress", fmt.Errorf("IPV4 GatewayIPAddress data  in request is missing")
	}
	return "", nil
}

// DeleteAddressPoolInfo stores the given addresspool against given fabric
func DeleteAddressPoolInfo(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fabricID := ctx.Params().Get("id")
	if _, err := capmodel.GetFabric(fabricID); errors.Is(err, db.ErrorKeyNotFound) {
		errMsg := fmt.Sprintf("Fabric data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"Fabric", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	addresspoolData, ok := capdata.AddressPoolDataStore[uri]
	if !ok {
		errMsg := fmt.Sprintf("AddressPool data for uri %s not found", uri)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"AddressPool", fabricID})
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(resp)
		return
	}
	if addresspoolData.AddressPool.Links != nil && len(addresspoolData.AddressPool.Links.Zones) > 0 {
		errMsg := fmt.Sprintf("AddressPool cannot be deleted as there are depZ Zone  still tied to it")
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceCannotBeDeleted, errMsg, []interface{}{uri, "AddressPool"})
		ctx.StatusCode(http.StatusNotAcceptable)
		ctx.JSON(resp)
		return
	}
	// Todo:Add the validation  to verify the links
	delete(capdata.AddressPoolDataStore, uri)
	ctx.StatusCode(http.StatusNoContent)
}

func getAddressPoolData(addresspoolOID string) (*model.AddressPool, int, interface{}) {
	addressPoolData, ok := capdata.AddressPoolDataStore[addresspoolOID]
	if !ok {
		errMsg := fmt.Sprintf("AddressPool data for uri %s not found", addresspoolOID)
		log.Error(errMsg)
		resp := updateErrorResponse(response.ResourceNotFound, errMsg, []interface{}{"AddressPool", addresspoolOID})

		return nil, http.StatusNotFound, resp
	}
	return addressPoolData.AddressPool, http.StatusOK, nil
}

func updateAddressPoolData(zoneOID, addresspoolOID, operation string) {
	addresspoolData := capdata.AddressPoolDataStore[addresspoolOID].AddressPool
	if addresspoolData.Links == nil {
		addresspoolData.Links = &model.AddressPoolLinks{}
	}
	if operation == "Add" {
		addresspoolData.Links.Zones = []model.Link{
			model.Link{
				Oid: zoneOID,
			},
		}
		addresspoolData.Links.ZonesCount = len(addresspoolData.Links.Zones)
	} else {
		addresspoolData.Links.Zones = []model.Link{}
		if len(addresspoolData.Links.Endpoints) == 0 {
			addresspoolData.Links = nil
		}
	}
	capdata.AddressPoolDataStore[addresspoolOID].AddressPool = addresspoolData
}
