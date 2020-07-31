package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type ConnectorGcpVpcRequest struct {
	CXP            string   `json:"cxp"`
	CredentialId   string   `json:"credentialId"`
	CustomerRegion string   `json:"customerRegion"`
	Group          string   `json:"group"`
	Name           string   `json:"name"`
	Segments       []string `json:"segments"`
	Size           string   `json:"size"`
	VpcId          string   `json:"vpcId"`
	VpcName        string   `json:"vpcName"`
}

type ConnectorGcpVpcResponse struct {
	Id int `json:"id"`
}

// Create a GCP-VPC connector
func (ac *AlkiraClient) CreateConnectorGcpVpc(connector *ConnectorGcpVpcRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/gcpvpcconnectors"
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

	var result ConnectorGcpVpcResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create GCP-VPC connector")
	}

	id = result.Id

	return id, nil
}

// Delete a GCP-VPC connector
func (ac *AlkiraClient) DeleteConnectorGcpVpc(connectorId string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/gcpvpcconnectors/" + connectorId

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
		return errors.New("Failed to delete GCP-VPC connector" + connectorId)
	}

	return nil
}
