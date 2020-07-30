package alkira

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TenantNetworkId struct {
	Id int `json:"id"`
}

type TenantNetworkState struct {
	State string `json:"state"`
}

// Get the tenant networks of the current tenant
func (ac *AlkiraClient) GetTenantNetworks() {
	url := ac.URI + "tenantnetworks"

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	log.Println(response.StatusCode)
	log.Println(string(data))

	return
}

// Get the tenant networks of the current tenant
func (ac *AlkiraClient) GetTenantNetworksId() (int) {
	var result []TenantNetworkId

	url := ac.URI + "tenantnetworks"

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	//log.Println(response.StatusCode)
	//log.Println(string(data))

	json.Unmarshal([]byte(data), &result)

	return result[0].Id
}

// Get the tenant network state
func (ac *AlkiraClient) GetTenantNetworkState() (string, error) {
	url := ac.URI + "tenantnetworks/" + ac.TenantNetworkId

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result TenantNetworkState
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 200 {
		return result.State, errors.New("Failed to get tenant network")
	}

	return result.State, nil
}

func (ac *AlkiraClient) ProvisionTenantNetwork() (string, error) {
	url := ac.URI + "tenantnetworks/" + ac.TenantNetworkId + "/provision"

	request, err := http.NewRequest("POST", url, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result TenantNetworkState
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return result.State, errors.New("Failed to provision tenant network")
	}
	return result.State, nil
}
