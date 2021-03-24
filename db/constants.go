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

package db

const (
	// scanPaginationSize defines the size of DB keys to be scanned on single query
	scanPaginationSize = 100
	// TableFabric is the table for storing switch and pod ids
	TableFabric = "ACI-Fabric"
	// TableSwitch is the table for storing switch information
	TableSwitch = "ACI-Switch"
	// TableSwitchChassis is the table for storing switch chassis information
	TableSwitchChassis = "ACI-SwitchChassis"
	// TableSwitchPorts is the table for storing ports of each switch
	TableSwitchPorts = "ACI-SwitchPorts"
	// TablePort is the table for storing port information
	TablePort = "ACI-Port"
	// TableZone is the table for storing zone information
	TableZone = "ACI-Zone"
)
