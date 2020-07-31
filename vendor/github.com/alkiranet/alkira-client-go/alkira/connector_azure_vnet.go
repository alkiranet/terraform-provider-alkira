package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type ConnectorAzureVnetRequest struct {
	CXP            string   `json:"cxp"`
	CredentialId   string   `json:"credentialId"`
	CustomerRegion string   `json:"customerRegion"`
	Group          string   `json:"group"`
	Name           string   `json:"name"`
	Segments       []string `json:"segments"`
	Size           string   `json:"size"`
	VnetId         string   `json:"vnetId"`
}

type ConnectorAzureVnetResponse struct {
	Id int `json:"id"`
}

// Create a AZURE-VNET connector
func (ac *AlkiraClient) CreateConnectorAzureVnet(connector *ConnectorAzureVnetRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/azurevnetconnectors"
	id  := 0

	// Construct the request
	body, err := json.Marshal(connector)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, err
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	log.Println(response.StatusCode)
	log.Println(string(data))

	var result ConnectorAzureVnetResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create AZURE-VNET connector")
	}

	id = result.Id

	return id, nil
}

// Delete a AZURE-VNET connector
func (ac *AlkiraClient) DeleteConnectorAzureVnet(connectorId string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/azurevnetconnectors/" + connectorId

	request, err := http.NewRequest("DELETE", url, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	log.Println(response.StatusCode)
	log.Println(string(data))

	if response.StatusCode != 200 {
		return errors.New("Failed to delete AZURE-VNET connector " + connectorId)
	}

	return nil
}
