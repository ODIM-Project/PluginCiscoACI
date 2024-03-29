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
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/capresponse"
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"
	pluginConfig "github.com/ODIM-Project/PluginCiscoACI/config"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

// TokenValidation validates sent token with the list of plugin generated tokens
func TokenValidation(token string) bool {
	var flag bool
	flag = false
	for _, v := range tokenDetails {
		if token == v.Token {
			flag = true
			if time.Since(v.LastUsed).Minutes() > pluginConfig.Data.SessionTimeoutInMinutes {
				return flag
			}
		}
	}
	return flag
}

// Validate does Basic authentication with device and returns UUID of device in response
func Validate(ctx iris.Context) {
	//Get token from Request
	if ctx.GetHeader("X-Auth-Token") == "" && ctx.GetHeader("Authorization") == "" {
		ctx.StatusCode(http.StatusUnauthorized)
		log.Error("No valid authorization")
		ctx.WriteString("No valid authorization")
		return
	}
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}
	var deviceDetails capmodel.Device
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		errMsg := "while trying to collect data from request, PluginCiscoACI got: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(errMsg)
		return
	}
	device := &caputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
	}

	redfishClient, err := caputilities.GetRedfishClient()
	if err != nil {
		log.Error(err.Error())
		capresponse.SetErrorResponse(ctx, http.StatusInternalServerError, response.InternalError, err.Error(), nil)
		return
	}

	//Get ServiceRoot of the device
	err = redfishClient.GetRootService(device)
	if err != nil {
		log.Error(err.Error())
		capresponse.SetErrorResponse(ctx, http.StatusInternalServerError, response.InternalError, err.Error(), nil)
		return
	}
	//Doing Get on device using basic Authentication
	resp, err := redfishClient.BasicAuthWithDevice(device, device.RootNode.Systems.Oid)
	if err != nil {
		log.Error(err.Error())
		capresponse.SetErrorResponse(ctx, http.StatusUnauthorized, response.ResourceAtURIUnauthorized, err.Error(), []interface{}{device.RootNode.Systems.Oid})
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		capresponse.SetErrorResponse(ctx, http.StatusInternalServerError, response.InternalError, err.Error(), nil)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		ctx.StatusCode(resp.StatusCode)
		ctx.JSON(string(body))
		return
	}
	if resp.StatusCode >= 300 {
		log.Error("Could not retreive ComputerSystems for " + device.Host + ":" + string(body))
	}

	response := capresponse.Device{
		ServerIP:   device.Host,
		Username:   device.Username,
		DeviceUUID: device.RootNode.UUID,
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.JSON(response)
}
