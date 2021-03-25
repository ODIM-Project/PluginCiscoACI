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
	"encoding/base64"
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

// GetClient returns a new connection client to APIC
func GetClient() *client.Client {
	aciClient = client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	return aciClient
}

// GetConnection returns a new connection to APIC
func GetConnection() *client.ServiceManager {
	aciClient = client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	aciServiceManager = client.NewServiceManager(client.DefaultMOURL, aciClient)
	return aciServiceManager
}

// GetFabricNodeData collects the all switch and fabric  details from the aci
func GetFabricNodeData() ([]*models.FabricNodeMember, error) {
	aciClient = client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	aciServiceManager = client.NewServiceManager(client.DefaultMOURL, aciClient)
	return aciServiceManager.ListFabricNodeMember()

}

//GetPortData collects the all port data for the given switch
func GetPortData(podID, ACISwitchID string) (*capmodel.PortCollectionResponse, error) {
	endpoint := fmt.Sprintf("https://%s/api/node/class/topology/pod-%s/node-%s/l1PhysIf.json", config.Data.APICConf.APICHost, podID, ACISwitchID)

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

//GetFabricHealth queries the fabric for it's Health from ACI
func GetFabricHealth(podID string) (*capmodel.FabricHealth, error) {
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

	var fabricHealthData capmodel.FabricHealth
	json.Unmarshal(body, &fabricHealthData)
	return &fabricHealthData, nil

}

// GetSwitchInfo collects the given switch data from the aci
func GetSwitchInfo(podID, ACISwitchID int) (*models.System, error) {
	return aciServiceManager.ReadSystem(podID, ACISwitchID)

}

// GetSwitchChassisInfo collects the given switch chassis data from the aci
func GetSwitchChassisInfo(podID, ACISwitchID string) (*capmodel.SwitchChassis, *capmodel.Health, error) {
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/ch.json", config.Data.APICConf.APICHost, podID, ACISwitchID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	newClient := ACIHTTPClient{}
	httpConf := &lutilconf.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	if newClient.httpClient, err = httpConf.GetHTTPClientObj(); err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Get on the URL %s is giving response with status code %d with response body %s", endpoint, resp.StatusCode, string(body))
		return nil, nil, fmt.Errorf(errMsg)
	}

	var switchChassisData capmodel.SwitchChassis
	var chassisHealth capmodel.Health
	json.Unmarshal(body, &switchChassisData)
	healthEndpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/ch/health.json", config.Data.APICConf.APICHost, podID, ACISwitchID)

	healthReq, err := http.NewRequest("GET", healthEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	healthReq.Close = true
	healthReq.Header.Set("Accept", "application/json")
	healthReq.AddCookie(&http.Cookie{
		Name:  "APIC-Cookie",
		Value: aciClient.AuthToken.Token,
	})
	healthReq.Close = true

	healthResp, err := newClient.httpClient.Do(healthReq)
	if err != nil {
		return nil, nil, err
	}
	defer healthResp.Body.Close()
	healthBody, err := ioutil.ReadAll(healthResp.Body)
	if err != nil {
		return nil, nil, err
	}
	json.Unmarshal(healthBody, &chassisHealth)
	return &switchChassisData, &chassisHealth, nil
}

//GetSwitchHealth queries the switch for it's Health from ACI
func GetSwitchHealth(podID, ACISwitchID string) (*capmodel.Health, error) {
	aciClient := client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	// Get the port data for given switch using the uri /api/node/mo/topology/{pod_id}/health.json
	err := aciClient.Authenticate()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/health.json", config.Data.APICConf.APICHost, podID, ACISwitchID)

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

	var switchHealthData capmodel.Health
	json.Unmarshal(body, &switchHealthData)
	return &switchHealthData, nil

}

//GetPortInfo collects the dat for  given port
func GetPortInfo(podID, ACISwitchID, portID string) (*capmodel.PortInfoResponse, error) {
	aciClient := client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	// Get the port data for given switch using the uri /api/node/mo/topology/{pod_id}/health.json
	err := aciClient.Authenticate()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/phys-[%s]/phys.json", config.Data.APICConf.APICHost, podID, ACISwitchID, portID)

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

//GetPortHealth collects the Health  for  given port
func GetPortHealth(podID, ACISwitchID, portID string) (*capmodel.Health, error) {
	aciClient := client.NewClient("https://"+config.Data.APICConf.APICHost, config.Data.APICConf.UserName, client.Password(config.Data.APICConf.Password), client.Insecure(true))
	// Get the port data for given switch using the uri /api/node/mo/topology/{pod_id}/health.json
	err := aciClient.Authenticate()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("https://%s/api/node/mo/topology/pod-%s/node-%s/sys/phys-[%s]/phys/health.json", config.Data.APICConf.APICHost, podID, ACISwitchID, portID)

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

	var portResponseData capmodel.Health
	json.Unmarshal(body, &portResponseData)
	return &portResponseData, nil

}

// GetPortPolicyGroup collects all policy group for given fabric and  switch
func GetPortPolicyGroup(podID, switchPath string) ([]*models.FabricPathEndpoint, error) {
	serviceManager := GetConnection()
	endPointUrL := fmt.Sprintf("/api/node/class/topology/pod-%s/protpaths%s/fabricPathEp.json", podID, switchPath)

	cont, err := serviceManager.GetViaURL(endPointUrL)
	list := models.FabricPathEndpointListFromContainer(cont)

	return list, err
}

// CheckValidityOfEthernet check if provided Ethernet is available in ODIM
func CheckValidityOfEthernet(reqURL string, odimUsername string, odimPassword string) (bool, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return false, err
	}
	newClient, err := GetRedfishClient()
	if err != nil {
		return false, err
	}
	auth := odimUsername + ":" + odimPassword
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	resp, err := newClient.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}
