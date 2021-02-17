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
	"net/http"

	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	iris "github.com/kataras/iris/v12"
)

//GetManagersCollection Fetches details of the manager collection
func GetManagersCollection(ctx iris.Context) {
	// TODO: implementation pending
	ctx.StatusCode(http.StatusNotImplemented)

}

//GetManagersInfo Fetches details of the given manager info
func GetManagersInfo(ctx iris.Context) {
	// TODO: implementation pending
	ctx.StatusCode(http.StatusNotImplemented)

}

func getInfoFromDevice(uri string, deviceDetails capmodel.Device, ctx iris.Context) {
	// TODO: implementation pending
}
