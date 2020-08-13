package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type ConnectorIPSecRequest struct {
	CXP            string      `json:"cxp"`
	Group          string      `json:"group"`
	Name           string      `json:"name"`
	SegmentOptions interface{} `json:"segmentOptions"`
	Segments       []string    `json:"segments"`
	Sites          interface{} `json:"sites"`
	Size           string      `json:"size"`
}

type ConnectorIPSecResponse struct {
	Id int `json:"id"`
}

type ConnectorIPSecSite struct {
	Name           string      `json:"name"`
	CustomerGwAsn  string      `json:"customerGwAsn"`
	CustomerGwIp   string      `json:"customerGwIp"`
    PresharedKeys  []string    `json:"presharedKeys"`
}

// Create a IPSEC connector
func (ac *AlkiraClient) CreateConnectorIPSec(connector *ConnectorIPSecRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/ipsecconnectors"
	id  := 0

	// Construct the request
	body, err := json.Marshal(connector)

	log.Println(string(body))
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

	var result ConnectorIPSecResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create IPSEC connector")
	}

	id = result.Id

	return id, nil
}

// Delete a IPSEC connector
func (ac *AlkiraClient) DeleteConnectorIPSec(connectorId string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/ipsecconnectors/" + connectorId

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
		return errors.New("Failed to delete IPSEC connector " + connectorId)
	}

	return nil
}
