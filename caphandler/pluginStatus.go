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
	"github.com/ODIM-Project/PluginCiscoACI/capresponse"
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"
	pluginConfig "github.com/ODIM-Project/PluginCiscoACI/config"
	iris "github.com/kataras/iris/v12"
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

}

// GetPluginStartup ...
func GetPluginStartup(ctx iris.Context) {
	// TODO: implementation pending
	ctx.StatusCode(http.StatusNotImplemented)
}
