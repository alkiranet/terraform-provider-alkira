package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PolicyRuleListRequest struct {
	Description string               `json:"description"`
	Name        string               `json:"name"`
	Rules       []PolicyRuleListRule `json:"rules"`
}

type PolicyRuleListRule struct {
	Priority int `json:"priority"`
	RuleId   int `json:"ruleId"`
}

type policyRuleListResponse struct {
	Id int `json:"id"`
}

// Create a policy rule
func (ac *AlkiraClient) CreatePolicyRuleList(p *PolicyRuleListRequest) (int, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rulelists", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(p)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreatePolicyRuleList: request failed: %v", err)
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

// Delete a policy rule list
func (ac *AlkiraClient) DeletePolicyRuleList(id int) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rulelists/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeletePolicyRuleList: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
