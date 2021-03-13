// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ConnectorAwsVpcRequest struct {
	BillingTags    []int    `json:"billingTags"`
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
	Id int `json:"id"`
}

// Create a AWS-VPC connector
func (ac *AlkiraClient) CreateConnectorAwsVpc(connector *ConnectorAwsVpcRequest) (int, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/awsvpcconnectors", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(connector)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreateConnectorAwsVpc: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 201 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	var result ConnectorAwsVpcResponse
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return id, fmt.Errorf("CreateConnectorAwsVpc: request failed: %v", err)
	}

	id = result.Id

	return id, nil
}

// Delete an AWS-VPC connector
func (ac *AlkiraClient) DeleteConnectorAwsVpc(id int) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/awsvpcconnectors/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteConnectorAwsVpc: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("DeleteConnectorAwsVpc: (%d) %s", response.StatusCode, string(data))
	}

	return nil
}
