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

//Package caputilities ...
package caputilities

import (
	"encoding/json"
	"fmt"
	lutilconf "github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/PluginCiscoACI/capmodel"
	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"io/ioutil"
	"net/http"
)

//ACIHTTPClient struct definition of HTTP wraper clinet used to communicate with ACI
type ACIHTTPClient struct {
	httpClient *http.Client
}

var aciClient *client.Client
var aciServiceManager *client.ServiceManager

// GetFabricNodeData collects the all switch and fabric  details from the aci
func GetFabricNodeData() ([]*models.FabricNodeMember, error) {
	aciClient = client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	aciServiceManager = client.NewServiceManager(client.DefaultMOURL, aciClient)
	return aciServiceManager.ListFabricNodeMember()

}

//GetPortData collects the all port data for the given switch
func GetPortData(podID, switchID string) (*capmodel.PortCollectionResponse, error) {
	endpoint := fmt.Sprintf("https://%s/api/node/class/topology/pod-%s/node-%s/l1PhysIf.json", config.Data.APICConf.APICHost, podID, switchID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	newClient := ACIHTTPClient{}
	httpConf := &lutilconf.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	if newClient.httpClient, err = httpConf.GetHTTPClientObj(); err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "APIC-Cookie",
		Value: aciClient.AuthToken.Token,
	})
	req.Close = true

	resp, err := newClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Get on the URL %s is giving response with status code %d with response body %s", endpoint, resp.StatusCode, string(body))
		return nil, fmt.Errorf(errMsg)
	}

	var portResponseData capmodel.PortCollectionResponse
	json.Unmarshal(body, &portResponseData)
	return &portResponseData, nil

}

//GetFabricHelath queries the fabric for it's helath from ACI
func GetFabricHelath(podID string) (*capmodel.FabricHelath, error) {
	aciClient := client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	// Get the port data for given switch using the uri /api/node/mo/topology/{pod_id}/health.json
	err := aciClient.Authenticate()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/health.json", config.Data.APICConf.APICHost, podID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	newClient := ACIHTTPClient{}
	httpConf := &lutilconf.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	if newClient.httpClient, err = httpConf.GetHTTPClientObj(); err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "APIC-Cookie",
		Value: aciClient.AuthToken.Token,
	})
	req.Close = true

	resp, err := newClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Get on the URL %s is giving response with status code %d with response body %s", endpoint, resp.StatusCode, string(body))
		return nil, fmt.Errorf(errMsg)
	}

	var fabricHelathData capmodel.FabricHelath
	json.Unmarshal(body, &fabricHelathData)
	return &fabricHelathData, nil

}

// GetSwitchInfo collects the given switch data from the aci
func GetSwitchInfo(podID, switchID int) (*models.System, error) {
	return aciServiceManager.ReadSystem(podID, switchID)

}

// GetSwitchChassisInfo collects the given switch chassis data from the aci
func GetSwitchChassisInfo(podID, switchID string) (*capmodel.SwitchChassis, error) {
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/ch.json", config.Data.APICConf.APICHost, podID, switchID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	newClient := ACIHTTPClient{}
	httpConf := &lutilconf.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	if newClient.httpClient, err = httpConf.GetHTTPClientObj(); err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "APIC-Cookie",
		Value: aciClient.AuthToken.Token,
	})
	req.Close = true

	resp, err := newClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Get on the URL %s is giving response with status code %d with response body %s", endpoint, resp.StatusCode, string(body))
		return nil, fmt.Errorf(errMsg)
	}

	var switchChassisData capmodel.SwitchChassis
	json.Unmarshal(body, &switchChassisData)
	return &switchChassisData, nil
}

//GetSwitchHelath queries the switch for it's helath from ACI
func GetSwitchHelath(podID, switchID string) (*capmodel.Helath, error) {
	aciClient := client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	// Get the port data for given switch using the uri /api/node/mo/topology/{pod_id}/health.json
	err := aciClient.Authenticate()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/health.json", config.Data.APICConf.APICHost, podID, switchID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	newClient := ACIHTTPClient{}
	httpConf := &lutilconf.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	if newClient.httpClient, err = httpConf.GetHTTPClientObj(); err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "APIC-Cookie",
		Value: aciClient.AuthToken.Token,
	})
	req.Close = true

	resp, err := newClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Get on the URL %s is giving response with status code %d with response body %s", endpoint, resp.StatusCode, string(body))
		return nil, fmt.Errorf(errMsg)
	}

	var switchHelathData capmodel.Helath
	json.Unmarshal(body, &switchHelathData)
	return &switchHelathData, nil

}

//GetPortInfo collects the dat for  given port
func GetPortInfo(podID, switchID, portID string) (*capmodel.PortInfoResponse, error) {
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/phys-[%s]/phys.json", config.Data.APICConf.APICHost, podID, switchID, portID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	newClient := ACIHTTPClient{}
	httpConf := &lutilconf.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	if newClient.httpClient, err = httpConf.GetHTTPClientObj(); err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "APIC-Cookie",
		Value: aciClient.AuthToken.Token,
	})
	req.Close = true

	resp, err := newClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Get on the URL %s is giving response with status code %d with response body %s", endpoint, resp.StatusCode, string(body))
		return nil, fmt.Errorf(errMsg)
	}

	var portResponseData capmodel.PortInfoResponse
	json.Unmarshal(body, &portResponseData)
	return &portResponseData, nil

}

//GetPortHealth collects the helath  for  given port
func GetPortHealth(podID, switchID, portID string) (*capmodel.Helath, error) {
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/phys-[%s]/phys/health.json", config.Data.APICConf.APICHost, podID, switchID, portID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	newClient := ACIHTTPClient{}
	httpConf := &lutilconf.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	if newClient.httpClient, err = httpConf.GetHTTPClientObj(); err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "APIC-Cookie",
		Value: aciClient.AuthToken.Token,
	})
	req.Close = true

	resp, err := newClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Get on the URL %s is giving response with status code %d with response body %s", endpoint, resp.StatusCode, string(body))
		return nil, fmt.Errorf(errMsg)
	}

	var portResponseData capmodel.Helath
	json.Unmarshal(body, &portResponseData)
	return &portResponseData, nil

}
