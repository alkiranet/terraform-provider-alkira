package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Name          string   `json:"name"`
	CustomerGwAsn string   `json:"customerGwAsn"`
	CustomerGwIp  string   `json:"customerGwIp"`
	PresharedKeys []string `json:"presharedKeys"`
}

// Create an IPSEC connector
func (ac *AlkiraClient) CreateConnectorIPSec(connector *ConnectorIPSecRequest) (int, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ipsecconnectors", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(connector)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreateConnectorIpSec: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result ConnectorIPSecResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	id = result.Id

	return id, nil
}

// Delete an IPSEC connector by Id
func (ac *AlkiraClient) DeleteConnectorIPSec(id int) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ipsecconnectors/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteConnectorIpSec: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
