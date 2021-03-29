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

// Packahe caphandler ...
package caphandler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/ODIM-Project/PluginCiscoACI/db"

	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

type mockConnector struct{}

func (d mockConnector) Create(table, resourceID, data string) error {
	return nil
}

func (d mockConnector) Update(table, resourceID, data string) error {
	return nil
}

func (d mockConnector) GetAllMatchingKeys(table, pattern string) ([]string, error) {
	return []string{"validID"}, nil
}

func (d mockConnector) Get(table, resourceID string) (string, error) {
	if resourceID == "validID" {
		switch table {
		case capmodel.TableFabric:
			return `{"SwitchData": ["test"], "PodID": "test"}`, nil
		case capmodel.TableSwitch:
			return `{"Id": "validID", "FabricID": "validID"}`, nil
		case capmodel.TableSwitchPorts:
			return `{"Id": "validID", "FabricID": "validID"}`, nil
		case capmodel.TablePort:
			return `{"Id": "validID", "FabricID": "validID"}`, nil
		case capmodel.TableZone:
			return `{"Id": "validID", "FabricID": "validID"}`, nil
		default:
		}
	}
	return "", fmt.Errorf("not found")
}

func TestGetManagerCollection(t *testing.T) {
	db.Connector = mockConnector{}
	config.SetUpMockConfig(t)
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")

	redfishRoutes.Get("/Managers", GetManagersCollection)

	e := httptest.New(t, mockApp)

	var deviceDetails = capmodel.Device{}
	//Unit Test for success scenario

	e.GET("/ODIM/v1/Managers").WithJSON(deviceDetails).Expect().Status(http.StatusOK)

}

func TestGetManager(t *testing.T) {
	config.SetUpMockConfig(t)
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")

	redfishRoutes.Get("/Managers", GetManagersInfo)
	var deviceDetails = capmodel.Device{}
	e := httptest.New(t, mockApp)
	//Unit Test for success scenario
	e.GET("/ODIM/v1/Managers").WithJSON(deviceDetails).Expect().Status(http.StatusOK)

}
