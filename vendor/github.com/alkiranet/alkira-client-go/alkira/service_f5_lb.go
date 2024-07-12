package alkira

import (
	"encoding/json"
	"fmt"
)

type ServiceF5Lb struct {
	Instances      []F5Instances `json:"instances"`
	Name           string        `json:"name"`
	Description    string        `json:"description,omitempty"`
	Cxp            string        `json:"cxp"`
	Size           string        `json:"size"`
	Id             json.Number   `json:"id"`              //RESPONSE ONLY
	State          string        `json:"state,omitempty"` //RESPONSE ONLY
	Segments       []string      `json:"segments"`
	BillingTags    []int         `json:"billingTags"`
	ElbCidrs       []string      `json:"elbCidrs"`
	BigIpAllowList []string      `json:"bigIpAllowList,omitempty"`
}

type F5Instances struct {
	Deployment               *InstanceDeployment `json:"deployment"`
	Name                     string              `json:"name"`
	RegistrationCredentialId string              `json:"registrationCredentialId"`
	CredentialId             string              `json:"credentialId"`
	LicenseType              string              `json:"licenseType"`
	Version                  string              `json:"version"`
}

type InstanceDeployment struct {
	Option string `json:"option"`
	Type   string `json:"type"`
}

func NewServiceF5Lb(ac *AlkiraClient) *AlkiraAPI[ServiceF5Lb] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/f5-lb-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceF5Lb]{ac, uri, true}
	return api

}
