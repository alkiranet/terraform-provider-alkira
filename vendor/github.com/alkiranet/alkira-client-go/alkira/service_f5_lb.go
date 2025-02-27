package alkira

import (
	"encoding/json"
	"fmt"
)

type ServiceF5Lb struct {
	SegmentOptions   F5SegmentOption `json:"segmentOptions"`
	Description      string          `json:"description,omitempty"`
	Cxp              string          `json:"cxp"`
	Size             string          `json:"size"`
	Id               json.Number     `json:"id,omitempty"`
	Name             string          `json:"name"`
	ServiceGroupName string          `json:"serviceGroupName"`
	ImplicitGroupId  int             `json:"implicitGroupId,omitempty"` // response only
	Instances        []F5Instance    `json:"instances"`
	Segments         []string        `json:"segments"`
	BillingTags      []int           `json:"billingTags,omitempty"`
	PrefixListId     int             `json:"prefixListId,omitempty"`
	GlobalCidrListId int             `json:"globalCidrListId"`
}

type F5Instance struct {
	Deployment               F5InstanceDeployment `json:"deployment"`
	Name                     string               `json:"name"`
	RegistrationCredentialId string               `json:"registrationCredentialId,omitempty"`
	CredentialId             string               `json:"credentialId"`
	HostNameFqdn             string               `json:"hostNameFqdn"`
	LicenseType              string               `json:"licenseType"`
	Version                  string               `json:"version"`
	Id                       int                  `json:"id,omitempty"` // RESPONSE ONLY
}

type F5SegmentOption map[string]F5SegmentSubOption
type F5SegmentSubOption struct {
	NatPoolPrefixLength int `json:"natPoolPrefixLength,omitempty"`
	ElbNicCount         int `json:"elbNicCount"`
}

type F5InstanceDeployment struct {
	Option string `json:"option,omitempty"`
	Type   string `json:"type"`
}

func NewServiceF5Lb(ac *AlkiraClient) *AlkiraAPI[ServiceF5Lb] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/f5-lb-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceF5Lb]{ac, uri, true}
	return api

}
