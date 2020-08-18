package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PolicyRuleRequest struct {
	Description    string                   `json:"description"`
	MatchCondition PolicyRuleMatchCondition `json:"matchCondition"`
	Name           string                   `json:"name"`
	RuleAction     PolicyRuleAction         `json:"ruleAction"`
}

type PolicyRuleMatchCondition struct {
	SrcIp                 string   `json:"srcIp"`
	DstIp                 string   `json:"dstIp"`
	Dscp                  string   `json:"dscp"`
	Protocol              string   `json:"protocol"`
	SrcPortList           []string `json:"srcPortList"`
	DstPortList           []string `json:"dstPortList"`
	ApplicationList       []string `json:"applicationList"`
	ApplicationFamilyList []string `json:"applicationFamilyList"`
	InternetApplicationId int      `json:"internetApplicationId"`
}

type PolicyRuleAction struct {
	Action          string   `json:"action"`
	ServiceTypeList []string `json:"serviceTypeList"`
}

type policyRuleResponse struct {
	Id int `json:"id"`
}

// Create a policy rule
func (ac *AlkiraClient) CreatePolicyRule(p *PolicyRuleRequest) (int, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(p)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreatePolicyRule: request failed: %v", err)
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

// Delete a policy rule
func (ac *AlkiraClient) DeletePolicyRule(id int) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeletePolicyRule: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
