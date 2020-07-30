package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type ConnectorAwsVpcRequest struct {
	CXP            string   `json:"cxp"`
	CredentialId   string   `json:"credentialId"`
	CustomerName   string   `json:"customerName"`
	CustomerRegion string   `json:"customerRegion"`
	Group          string   `json:"group"`
	Name           string   `json:"name"`
	Segments       []string `json:"segments"`
	Size           string   `json:"size"`
	VpcId          string   `json:"vpcId"`
	VpcOwnerId     string   `json:"vpcOwnerId"`
}

type ConnectorAwsVpcResponse struct {
	Id              int         `json:"id"`
}

// Create a AWS-VPC connector
func (ac *AlkiraClient) CreateConnectorAwsVpc(connector *ConnectorAwsVpcRequest) (int, error) {
	url := ac.URI + "api/v1/tenantnetworks/" + ac.TenantNetworkId + "/awsvpcconnectors"
	id  := 0

	// Construct the request
	body, err := json.Marshal(connector)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		log.Printf("Error : %s", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	log.Println(response.StatusCode)
	log.Println(string(data))

	var result ConnectorAwsVpcResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create AWS-VPC connector")
	}

	id = result.Id

	return id, nil
}

// Delete a AWS-VPC connector
func (ac *AlkiraClient) DeleteConnectorAwsVpc(connectorId string) (error) {
	url := ac.URI + "api/v1/tenantnetworks/" + ac.TenantNetworkId + "/awsvpcconnectors/" + connectorId

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
		return errors.New("Failed to delete AWS-VPC connector " + connectorId)
	}

	return nil
}
