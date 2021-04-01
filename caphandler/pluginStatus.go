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
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/PluginCiscoACI/capmessagebus"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/capresponse"
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"
	"github.com/ODIM-Project/PluginCiscoACI/config"
	pluginConfig "github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/ODIM-Project/PluginCiscoACI/constants"
	iris "github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// GetPluginStatus defines the GetPluginStatus iris handler.
// and returns status
func GetPluginStatus(ctx iris.Context) {
	//Get token from Request
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
	var messageQueueInfo []capresponse.EmbQueue
	var resp = capresponse.PluginStatusResponse{
		Comment: "Plugin Status Response",
		Name:    "Common Redfish Plugin Status",
		Version: pluginConfig.Data.FirmwareVersion,
	}
	resp.Status = caputilities.Status
	resp.Status.TimeStamp = time.Now().Format(time.RFC3339)
	resp.EventMessageBus = capresponse.EventMessageBus{
		EmbType: pluginConfig.Data.MessageBusConf.EmbType,
	}
	//messageQueueInfo := make([]capresponse.EmbQueue, 0)
	for i := 0; i < len(pluginConfig.Data.MessageBusConf.EmbQueue); i++ {
		messageQueueInfo = append(messageQueueInfo, capresponse.EmbQueue{
			QueueName: pluginConfig.Data.MessageBusConf.EmbQueue[i],
			QueueDesc: "Queue for redfish events",
		})
	}
	resp.EventMessageBus.EmbQueue = messageQueueInfo

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
	if capmodel.PluginIntialStatus == false {
		go publishResourceAddedEvent()
		capmodel.PluginIntialStatus = true
	}

}

// GetPluginStartup ...
func GetPluginStartup(ctx iris.Context) {
	// TODO: implementation pending
	ctx.StatusCode(http.StatusNotImplemented)
}

func publishResourceAddedEvent() {
	time.Sleep(5 * time.Second)
	// Send resource added event odim
	allFabric, err := capmodel.GetAllFabric("")
	if err != nil {
		log.Fatal("while fetching all stored fabric data got: " + err.Error())
	}
	for fabricID := range allFabric {
		var event = common.Event{
			EventID:   uuid.NewV4().String(),
			MessageID: constants.ResourceCreatedMessageID,
			EventType: "ResourceAdded",
			OriginOfCondition: &common.Link{
				Oid: "/ODIM/v1/Fabrics/" + fabricID,
			},
		}
		var events = []common.Event{event}
		var messageData = common.MessageData{
			Name:      "Fabric added event",
			Context:   "/redfish/v1/$metadata#Event.Event",
			OdataType: constants.EventODataType,
			Events:    events,
		}
		data, _ := json.Marshal(messageData)
		eventData := common.Events{
			IP:      config.Data.LoadBalancerConf.Host,
			Request: data,
		}
		capmessagebus.Publish(eventData)
	}
}
