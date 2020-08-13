package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type PolicyRuleListRequest struct {
	Description    string               `json:"description"`
	Name           string               `json:"name"`
	Rules          []PolicyRuleListRule `json:"rules"`
}

type PolicyRuleListRule struct {
	Priority    int   `json:"priority"`
	RuleId      int   `json:"ruleId"`
}

type policyRuleListResponse struct {
	Id int `json:"id"`
}

// Create a policy rule
func (ac *AlkiraClient) CreatePolicyRuleList(p *PolicyRuleListRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/policy/rulelists"
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
		return id, errors.New("Failed to create policy rule list")
	}

	id = result.Id
	return id, nil
}

// Delete a policy rule list
func (ac *AlkiraClient) DeletePolicyRuleList(id string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/policy/rulelists/" + id

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
		return errors.New("Failed to delete policy rule list" + id)
	}

	return nil
}
