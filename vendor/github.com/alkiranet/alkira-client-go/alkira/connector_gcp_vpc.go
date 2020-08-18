package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ConnectorGcpVpc struct {
	CXP            string   `json:"cxp"`
	CredentialId   string   `json:"credentialId"`
	CustomerRegion string   `json:"customerRegion"`
	Group          string   `json:"group"`
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Segments       []string `json:"segments"`
	Size           string   `json:"size"`
	VpcId          string   `json:"vpcId"`
	VpcName        string   `json:"vpcName"`
}

// Create a GCP-VPC connector
func (ac *AlkiraClient) CreateConnectorGcpVpc(connector *ConnectorGcpVpc) (int, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/gcpvpcconnectors", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(connector)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreateConnectorGcpVpc: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 && response.StatusCode != 201 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	var result ConnectorGcpVpc
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return id, fmt.Errorf("CreateConnectorGcpVpc: parse failed: %v", err)
	}

	id = result.Id

	return id, nil
}

// Get a GCP-VPC connector
func (ac *AlkiraClient) GetConnectorGcpVpc(id int) (ConnectorGcpVpc, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/gcpvpcconnectors/%d", ac.URI, ac.TenantNetworkId, id)
	var result ConnectorGcpVpc

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return result, fmt.Errorf("GetConnectorGcpVpc: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return result, fmt.Errorf("GetConnectorGcpVpc: (%d) %s", response.StatusCode, string(data))
	}

	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return result, fmt.Errorf("GetConnectorGcpVpc: parse failed: %v", err)
	}

	return result, nil
}

// Delete a GCP-VPC connector
func (ac *AlkiraClient) DeleteConnectorGcpVpc(id int) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/gcpvpcconnectors/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteConnectorGcpVpc: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("DeleteConnectorGcpVpc: (%d) %s", response.StatusCode, string(data))
	}

	return nil
}
