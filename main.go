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
package main

import (
	"encoding/json"
	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	lutilconf "github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/caphandler"
	"github.com/ODIM-Project/PluginCiscoACI/capmessagebus"
	"github.com/ODIM-Project/PluginCiscoACI/capmiddleware"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/caputilities"
	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/ODIM-Project/PluginCiscoACI/constants"
	"github.com/ciscoecosystem/aci-go-client/models"
	iris "github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var subscriptionInfo []capmodel.Device
var log = logrus.New()

// TokenObject will contains the generated token and public key of odimra
type TokenObject struct {
	AuthToken string `json:"authToken"`
	PublicKey []byte `json:"publicKey"`
}

func main() {
	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatal("Plugin Service should not be run as the root user")
	}

	if err := config.SetConfiguration(); err != nil {
		log.Fatal("while reading from config, PluginCiscoACI got" + err.Error())
	}

	if err := dc.SetConfiguration(config.Data.MessageBusConf.MessageQueueConfigFilePath); err != nil {
		log.Fatal("while trying to set messagebus configuration, PluginCiscoACI got: " + err.Error())
	}

	// CreateJobQueue defines the queue which will act as an infinite buffer
	// In channel is an entry or input channel and the Out channel is an exit or output channel
	caphandler.In, caphandler.Out = common.CreateJobQueue()

	// RunReadWorkers will create a worker pool for doing a specific task
	// which is passed to it as Publish method after reading the data from the channel.
	go common.RunReadWorkers(caphandler.Out, capmessagebus.Publish, 1)
	intializeACIData()
	intializePluginStatus()
	app()
}

func app() {
	app := routers()
	go func() {
		eventsrouters()
	}()
	conf := &lutilconf.HTTPConfig{
		Certificate:   &config.Data.KeyCertConf.Certificate,
		PrivateKey:    &config.Data.KeyCertConf.PrivateKey,
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
		ServerAddress: config.Data.PluginConf.Host,
		ServerPort:    config.Data.PluginConf.Port,
	}
	pluginServer, err := conf.GetHTTPServerObj()
	if err != nil {
		log.Fatal("while initializing plugin server, PluginCiscoACI got: " + err.Error())
	}
	app.Run(iris.Server(pluginServer))
}

func routers() *iris.Application {
	app := iris.New()
	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' && path[len(path)-2] != '/' {
			path = path[:len(path)-1]
			r.RequestURI = path
			r.URL.Path = path
		}
		next(w, r)
	})

	pluginRoutes := app.Party("/ODIM/v1")
	pluginRoutes.Post("/validate", capmiddleware.BasicAuth, caphandler.Validate)
	pluginRoutes.Post("/Sessions", caphandler.CreateSession)
	pluginRoutes.Post("/Subscriptions", capmiddleware.BasicAuth, caphandler.CreateEventSubscription)
	pluginRoutes.Delete("/Subscriptions", capmiddleware.BasicAuth, caphandler.DeleteEventSubscription)
	pluginRoutes.Get("/Status", capmiddleware.BasicAuth, caphandler.GetPluginStatus)
	pluginRoutes.Post("/Startup", capmiddleware.BasicAuth, caphandler.GetPluginStartup)

	pluginRoutes.Get("/Fabrics", caphandler.GetFabricResource)
	pluginRoutes.Get("/Fabrics/{id}", caphandler.GetFabricData)
	pluginRoutes.Get("/Fabrics/{id}/Switches", caphandler.GetSwitchCollection)
	pluginRoutes.Get("/Fabrics/{id}/Switches/{rid}", caphandler.GetSwitchInfo)
	pluginRoutes.Get("/Fabrics/{id}/Switches/{rid}/Ports", caphandler.GetFabricResource)
	pluginRoutes.Get("/Fabrics/{id}/Switches/{id2}/Ports/{rid}", caphandler.GetFabricResource)
	pluginRoutes.Get("/Fabrics/{id}/Zones", caphandler.GetFabricResource)
	pluginRoutes.Post("/Fabrics/{id}/Zones", caphandler.GetFabricResource)
	pluginRoutes.Get("/Fabrics/{id}/Zones/{rid}", caphandler.GetFabricResource)
	pluginRoutes.Delete("/Fabrics/{id}/Zones/{rid}", caphandler.GetFabricResource)
	pluginRoutes.Patch("/Fabrics/{id}/Zones/{rid}", caphandler.GetFabricResource)
	pluginRoutes.Get("/Fabrics/{id}/AddressPools", caphandler.GetFabricResource)
	pluginRoutes.Post("/Fabrics/{id}/AddressPools", caphandler.GetFabricResource)
	pluginRoutes.Get("/Fabrics/{id}/AddressPools/{rid}", caphandler.GetFabricResource)
	pluginRoutes.Delete("/Fabrics/{id}/AddressPools/{rid}", caphandler.GetFabricResource)

	managers := pluginRoutes.Party("/Managers")
	managers.Get("/", caphandler.GetManagersCollection)
	managers.Get("/{id}", caphandler.GetManagersInfo)
	taskmon := pluginRoutes.Party("/taskmon")
	taskmon.Get("/{TaskID}", caphandler.GetTaskMonitor)

	task := pluginRoutes.Party("/TaskService")
	task.Get("/", caphandler.GetTaskService)
	task.Get("/Tasks", caphandler.GetTaskService)
	task.Get("/Tasks/{TaskID}", caphandler.GetTaskService)
	task.Get("/Tasks/{TaskID}/SubTasks", caphandler.GetTaskService)
	task.Get("/Tasks/{TaskID}/SubTasks/{subTaskID}", caphandler.GetTaskService)
	task.Delete("/Tasks/{TaskID}", caphandler.GetTaskService)

	return app
}

func eventsrouters() {
	app := iris.New()
	app.Post(config.Data.EventConf.DestURI, caphandler.RedfishEvents)
	conf := &lutilconf.HTTPConfig{
		Certificate:   &config.Data.KeyCertConf.Certificate,
		PrivateKey:    &config.Data.KeyCertConf.PrivateKey,
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
		ServerAddress: config.Data.EventConf.ListenerHost,
		ServerPort:    config.Data.EventConf.ListenerPort,
	}
	evtServer, err := conf.GetHTTPServerObj()
	if err != nil {
		log.Fatal("while initializing event server, PluginCiscoACI got: " + err.Error())
	}
	app.Run(iris.Server(evtServer))
}

// intializePluginStatus sets plugin status
func intializePluginStatus() {
	caputilities.Status.Available = "no"
	caputilities.Status.Uptime = time.Now().Format(time.RFC3339)
}

// intializeACIData reads required fabric,switch and port data from aci and stored it in the data store
func intializeACIData() {
	capdata.FabricDataStore.Data = make(map[string]*capdata.Fabric)
	capdata.SwitchDataStore.Data = make(map[string]*models.FabricNodeMember, 0)
	aciNodesData, err := caputilities.GetFabricNodeData()
	if err != nil {
		log.Fatal("while intializing ACI Data  PluginCiscoACI got: " + err.Error())
	}
	for _, aciNodeData := range aciNodesData {
		switchID := uuid.NewV4().String() + ":" + aciNodeData.NodeId
		capdata.FabricDataStore.Lock.Lock()
		fabricID := config.Data.RootServiceUUID + ":" + aciNodeData.FabricId
		if data, ok := capdata.FabricDataStore.Data[fabricID]; ok {
			data.SwitchData = append(data.SwitchData, switchID)
			data.PodID = aciNodeData.PodId
		} else {

			capdata.FabricDataStore.Data[fabricID] = &capdata.Fabric{
				SwitchData: []string{
					switchID,
				},
				PodID: aciNodeData.PodId,
			}
		}
		capdata.FabricDataStore.Lock.Unlock()

		capdata.SwitchDataStore.Lock.Lock()
		capdata.SwitchDataStore.Data[switchID] = aciNodeData
		capdata.SwitchDataStore.Lock.Unlock()
	}

	// TODO:
	// adding logic to collect the ports data
	// registering the for the aci events

	//updating the plugin status
	caputilities.Status.Available = "yes"
	// Send resource added event odim
	capdata.FabricDataStore.Lock.RLock()
	for fabricID := range capdata.FabricDataStore.Data {
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
	capdata.FabricDataStore.Lock.RUnlock()

	return
}
