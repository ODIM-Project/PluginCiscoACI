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

//Package capmodel ...
package capmodel

//FabricHealth ...
type FabricHealth struct {
	TotalCount string               `json:"totalCount"`
	IMData     []FabricHealthIMData `json:"imdata"`
}

// FabricHealthIMData ...
type FabricHealthIMData struct {
	FabricHealthData FabricHealthData `json:"fabricHealthTotal"`
}

// FabricHealthData ...
type FabricHealthData struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// ErrorResponse struct ...
type ErrorResponse struct {
	IMData     []ErrorIMData `json:"imdata"`
	TotalCount string        `json:"totalCount"`
}

// ErrorIMData strcut ...
type ErrorIMData struct {
	Error AttStruct `json:"error"`
}

// AttStruct strcut ...
type AttStruct struct {
	Attributes Attribute `json:"attributes"`
}

// Attribute struct ...
type Attribute struct {
	code string `json:"code"`
	text string `json:"text"`
}
