package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type PolicyRuleRequest struct {
	Description    string                   `json:"description"`
	MatchCondition PolicyRuleMatchCondition `json:"matchCondition"`
	Name           string                   `json:"name"`
	RuleAction     PolicyRuleAction         `json:"ruleAction"`
}

type PolicyRuleMatchCondition struct {
	SrcIp       string   `json:"srcIp"`
	DstIp       string   `json:"dstIp"`
	Dscp        string   `json:"dscp"`
	Protocol    string   `json:"protocol"`
	SrcPortList []string `json:"srcPortList"`
	DstPortList []string `json:"dstPortList"`
}

type PolicyRuleAction struct {
	Action       string   `json:"action"`
}

type policyRuleResponse struct {
	Id int `json:"id"`
}

// Create a policy rule
func (ac *AlkiraClient) CreatePolicyRule(p *PolicyRuleRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/policy/rules"
	id  := 0

	// Construct the request
	body, err := json.Marshal(p)

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

	var result policyResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create policy rule")
	}

	id = result.Id
	return id, nil
}

// Delete a policy rule
func (ac *AlkiraClient) DeletePolicyRule(id string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/policy/rules/" + id

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
		return errors.New("Failed to delete policy rule" + id)
	}

	return nil
}
