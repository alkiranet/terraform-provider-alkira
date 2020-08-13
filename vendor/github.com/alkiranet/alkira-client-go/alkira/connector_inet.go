package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type ConnectorInetRequest struct {
	CXP            string   `json:"cxp"`
	Group          string   `json:"group"`
	Name           string   `json:"name"`
	Segments       []string `json:"segments"`
}

type ConnectorInetResponse struct {
	Id              int         `json:"id"`
}

// Create a INET connector
func (ac *AlkiraClient) CreateConnectorInet(connector *ConnectorInetRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/internetconnectors"
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

	var result ConnectorInetResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create INET connector")
	}

	id = result.Id

	return id, nil
}

// Delete an INET connector
func (ac *AlkiraClient) DeleteConnectorInet(connectorId string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/internetconnectors/" + connectorId

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
		return errors.New("Failed to delete INET connector " + connectorId)
	}

	return nil
}
