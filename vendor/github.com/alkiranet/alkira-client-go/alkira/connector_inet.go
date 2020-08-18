package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ConnectorInternet struct {
	CXP         string   `json:"cxp"`
	Description string   `json:"description"`
	Group       string   `json:"group"`
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Segments    []string `json:"segments"`
	Size        string   `json:"size"`
}

// Create an internet connector
func (ac *AlkiraClient) CreateConnectorInternet(connector *ConnectorInternet) (int, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internetconnectors", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(connector)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreateConnectorInternet: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 201 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	var result ConnectorInternet
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return id, fmt.Errorf("CreateConnectorInternet: parse failed: %v", err)
	}

	id = result.Id

	return id, nil
}

// Get an internet connector by id
func (ac *AlkiraClient) GetConnectorInternet(id int) (ConnectorInternet, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internetconnectors/%d", ac.URI, ac.TenantNetworkId, id)
	var result ConnectorInternet

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return result, fmt.Errorf("GetConnectorInternet: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return result, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return result, fmt.Errorf("GetConnectorInternet: parse failed: %v", err)
	}

	return result, nil
}

// Delete an internet connector
func (ac *AlkiraClient) DeleteConnectorInet(id int) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internetconnectors/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteConnectorInternet: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
