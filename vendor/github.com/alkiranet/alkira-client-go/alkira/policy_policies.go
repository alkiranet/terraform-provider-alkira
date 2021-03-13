// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PolicyRequest struct {
	Description string   `json:"description"`
	Enabled     string   `json:"enabled"`
	FromGroups  []string `json:"fromGroups"`
	Name        string   `json:"name"`
	RuleListId  string   `json:"ruleListId"`
	SegmentIds  []string `json:"segmentIds"`
	ToGroups    []string `json:"toGroups"`
}

type policyResponse struct {
	Id int `json:"id"`
}

// Create a policy
func (ac *AlkiraClient) CreatePolicy(p *PolicyRequest) (int, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/policies", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(p)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreatePolicy: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result policyResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	id = result.Id

	return id, nil
}

// Delete a policy
func (ac *AlkiraClient) DeletePolicy(id int) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/policies/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeletePolicy: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
