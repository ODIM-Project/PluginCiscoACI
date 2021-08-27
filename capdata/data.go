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
)

//Fabric ACI data of switch id and pod id
type Fabric struct {
	SwitchData []string
	PodID      string
}

// EndpointData hold the EndpointData data
type EndpointData struct {
	Endpoint           *model.Endpoint
	ACIPolicyGroupData *ACIPolicyGroupData
}

// ACIPolicyGroupData holds info regarding the ACI policy profile
type ACIPolicyGroupData struct {
	PolicyGroupDN             string
	SwitchProfileName         string
	SwitchAssociationName     string
	SwitchProfileSelectorName string
	AccessPortSelectorName    string
	PcVPCPolicyGroupName      string
	PCVPCPolicyGroupDN        string
}

// ACIDomainData hold dn of ACI DOMAIN and DomaineEntity
type ACIDomainData struct {
	DomainDN              string
	DomainEntityProfileDn string
}
