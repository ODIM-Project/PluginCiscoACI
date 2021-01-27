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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"
	evtConfig "github.com/ODIM-Project/PluginCiscoACI/config"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//CreateEventSubscription : Subscribes for events
func CreateEventSubscription(ctx iris.Context) {

	device, deviceDetails, err := getDeviceDetails(ctx)
	if err != nil {
		return
	}
	//First delete existing matching subscription(our subscription) from device
	deleteMatchingSubscriptions(device)

	var reqPostBody capmodel.EvtSubPost
	var reqData string

	//replacing the reruest  with south bound translation URL
	for key, value := range evtConfig.Data.URLTranslation.SouthBoundURL {
		reqData = strings.Replace(string(deviceDetails.PostBody), key, value, -1)
	}

	err = json.Unmarshal([]byte(reqData), &reqPostBody)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	// remove the mesaageids, resourcestypes and originresources from the request and post it to device
	// since some of device doesnt support these
	req := capmodel.EvtSubPost{
		Destination: "https://" + evtConfig.Data.LoadBalancerConf.Host + ":" + evtConfig.Data.LoadBalancerConf.Port + evtConfig.Data.EventConf.DestURI,
		EventTypes:  reqPostBody.EventTypes,
		Context:     reqPostBody.Context,
		HTTPHeaders: reqPostBody.HTTPHeaders,
		Protocol:    reqPostBody.Protocol,
	}
	device.PostBody, err = json.Marshal(req)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	redfishClient, err := caputilities.GetRedfishClient()
	if err != nil {
		log.Error(err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	var resp *http.Response
	//Subscribe to Events
	resp, err = redfishClient.SubscribeForEvents(device)
	if err != nil {
		log.Error(err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	defer resp.Body.Close()
	if err := validateResponse(ctx, device, resp, http.MethodPost); err != nil {
		return
	}
}

// Delete match subscription from device
func deleteMatchingSubscriptions(device *caputilities.RedfishDevice) {
	// get all subscriptions
	device.Location = "https://" + device.Host + "/redfish/v1/EventService/Subscriptions"
	redfishClient, err := caputilities.GetRedfishClient()
	if err != nil {
		log.Error(err.Error())
		return
	}

	//Get Subscription details to check if it is really ours
	resp, err := redfishClient.GetSubscriptionDetail(device)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage := "while getting subscription details for URI " + device.Location + " PluginCiscoACI got: " + strconv.Itoa(resp.StatusCode)
		log.Error(errorMessage)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var subscriptionCollectionBody interface{}
	err = json.Unmarshal(body, &subscriptionCollectionBody)
	if err != nil {
		log.Error(err.Error())
		return
	}
	members := subscriptionCollectionBody.(map[string]interface{})["Members"]
	for _, member := range members.([]interface{}) {
		device.Location = member.(map[string]interface{})["@odata.id"].(string)
		device.Location = "https://" + device.Host + device.Location
		if isOurSubscription(device) {
			resp, err = redfishClient.DeleteSubscriptionDetail(device)
			if err != nil {
				log.Error(err.Error())
				return
			}
			resp.Body.Close()
		}
	}
	return
}

func isOurSubscription(device *caputilities.RedfishDevice) bool {

	redfishClient, err := caputilities.GetRedfishClient()
	if err != nil {
		log.Error(err.Error())
		return false
	}
	//Get Subscription details to check if it is really ours
	resp, err := redfishClient.GetSubscriptionDetail(device)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage := "while getting subscription details for URI " + device.Location + ", PluginCiscoACI got " + strconv.Itoa(resp.StatusCode)
		log.Error(errorMessage)
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	var subscriptionBody interface{}
	err = json.Unmarshal(body, &subscriptionBody)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	subscriptionDestinationFromDevice := subscriptionBody.(map[string]interface{})["Destination"].(string)
	// if the subscription is ours then the destination should match with LBHOST:LBPORT.
	//If it is not matching then retrun with MethodNotAllowed
	if !strings.Contains(subscriptionDestinationFromDevice, evtConfig.Data.LoadBalancerConf.Host+":"+evtConfig.Data.LoadBalancerConf.Port) {
		return false
	}
	return true
}

//DeleteEventSubscription : Delete subscription
func DeleteEventSubscription(ctx iris.Context) {
	device, _, err := getDeviceDetails(ctx)
	if err != nil {
		return
	}
	redfishClient, err := caputilities.GetRedfishClient()
	if err != nil {
		log.Error(err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	//Delete Subscription details
	resp, err := redfishClient.DeleteSubscriptionDetail(device)
	if err != nil {
		log.Error(err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	defer resp.Body.Close()
	if err := validateResponse(ctx, device, resp, http.MethodDelete); err != nil {
		return
	}
}

// getDeviceDetails will accepts iris context and it will extract device details from context
// then decrypt the password and return device details
func getDeviceDetails(ctx iris.Context) (*caputilities.RedfishDevice, *capmodel.Device, error) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return nil, nil, fmt.Errorf("Invalid/Expired X-Auth-Token")
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
		return nil, nil, err
	}

	device := &caputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
		Location: deviceDetails.Location,
	}

	device.Password = string(deviceDetails.Password)

	return device, &deviceDetails, nil
}

// validateResponse will accepts iris context to write status code and resopnse
// method is to return status created incase of create subscription
// otherwise return statusok
func validateResponse(ctx iris.Context, device *caputilities.RedfishDevice, resp *http.Response, method string) error {
	var body []byte
	var err error
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return err
	}
	defer resp.Body.Close()
	if strings.EqualFold(method, http.MethodPost) {
		// if there was an error for message ids means device haven't support of MessageIds
		// So remove the MessageIds from the request and subscribe again.
		if resp.StatusCode != http.StatusCreated {
			removeMessageID(ctx, device)
			// Subscribe to Events
			redfishClient, err := caputilities.GetRedfishClient()
			if err != nil {
				log.Error(err.Error())
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return err
			}

			resp, err = redfishClient.SubscribeForEvents(device)
			if err != nil {
				log.Error(err.Error())
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return err
			}
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error(err.Error())
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return err
			}
		}

	}
	header := make(map[string]string)
	for k, v := range resp.Header {
		var value string
		for i := 0; i < len(v); i++ {
			value = value + " " + v[i]
		}
		header[k] = value
	}
	if resp.StatusCode == 401 {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Authtication with the device failed")
		return errors.New("Authtication with the device failed")
	}
	if resp.StatusCode >= 300 {
		log.Warn("Subscription operation failed: " + string(body))
	}
	common.SetResponseHeader(ctx, header)
	ctx.StatusCode(resp.StatusCode)
	log.Info("Redfish plugin response body: " + string(body))
	ctx.WriteString(string(body))
	return nil
}

func removeMessageID(ctx iris.Context, device *caputilities.RedfishDevice) {
	var ReqPostBody capmodel.EvtSubPost
	err := json.Unmarshal(device.PostBody, &ReqPostBody)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	ReqPostBody.MessageIds = []string{}
	device.PostBody, err = json.Marshal(ReqPostBody)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	return
}
