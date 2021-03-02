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

// Package capdata ...
package capdata

import (
	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"sync"
)

// FabricData ...
type FabricData struct {
	Data map[string]*Fabric
	Lock sync.RWMutex
}

//Fabric ACI data of switch id and pod id
type Fabric struct {
	SwitchData []string
	PodID      string
}

// FabricDataStore holds the value  aci fabric id and switches
var FabricDataStore FabricData

// SwitchData ...
type SwitchData struct {
	Data map[string]*model.Switch
	Lock sync.RWMutex
}

// SwitchDataStore holds the value  aci switch id and switches info
var SwitchDataStore SwitchData

// SwitchToPortDataStore hold the value of the ports belonging to respective switches
var SwitchToPortDataStore map[string][]string

//PortDataStore hold the value of the ports info of the switch
var PortDataStore map[string]*model.Port
