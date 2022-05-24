package alkira

import (
	"encoding/json"
	"fmt"
)

type Infoblox struct {
	AnyCast          InfobloxAnycast    `json:"anycast"`
	BillingTags      []int              `json:"billingTags"`
	Cxp              string             `json:"cxp"`
	Description      string             `json:"description,omitempty"`
	GlobalCidrListId int                `json:"globalCidrListId"`
	GridMaster       InfobloxGridMaster `json:"gridMaster"`
	Id               json.Number        `json:"id,omitempty"`
	Instances        []InfobloxInstance `json:"instances"`
	InternalName     string             `json:"internalName,omitempty"`
	LicenseType      string             `json:"licenseType,omitempty"`
	Name             string             `json:"name"`
	Segments         []string           `json:"segments"`
	ServiceGroupId   int                `json:"serviceGroupId,omitempty"`
	ServiceGroupName string             `json:"serviceGroupName"`
	Size             string             `json:"size,omitempty"`
}

type InfobloxAnycast struct {
	BackupCxps []string `json:"backupCxps,omitempty"`
	Enabled    bool     `json:"enabled"`
	Ips        []string `json:"ips,omitempty"`
}

type InfobloxGridMaster struct {
	External                 bool   `json:"external,omitempty"`
	GridMasterCredentialId   string `json:"gridMasterCredentialId"`
	Ip                       string `json:"ip"`
	Name                     string `json:"name"`
	SharedSecretCredentialId string `json:"sharedSecretCredentialId"`
}

type InfobloxInstance struct {
	AnyCastEnabled     bool        `json:"anyCastEnabled"`
	ConfiguredMasterIp string      `json:"configuredMasterIp,omitempty"`
	CredentialId       string      `json:"credentialId"`
	HostName           string      `json:"hostName"`
	Id                 json.Number `json:"id,omitempty"`
	InternalName       string      `json:"internalName,omitempty"`
	LanPrefix          string      `json:"lanPrefix,omitempty"`
	ManagementPrefix   string      `json:"managementPrefix,omitempty"`
	Model              string      `json:"model"`
	Name               string      `json:"name,omitempty"`
	ProductId          string      `json:"productId,omitempty"`
	PublicIp           string      `json:"publicIp,omitempty"`
	Type               string      `json:"type"`
	Version            string      `json:"version,omitempty"`
}

func (ac *AlkiraClient) CreateInfoblox(in *Infoblox) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/infoblox-services", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(in)

	if err != nil {
		return "", fmt.Errorf("CreateInfoblox: marshal failed: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result Infoblox
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateInfoblox: failed to unmarshal: %v", err)
	}

	return result.Id.String(), nil
}

func (ac *AlkiraClient) GetAllInfoblox() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/infoblox-services", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ac *AlkiraClient) GetInfobloxById(id string) (*Infoblox, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/infoblox-services/%s", ac.URI, ac.TenantNetworkId, id)

	var infoblox Infoblox

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &infoblox)

	if err != nil {
		return nil, fmt.Errorf("GetInfobloxById: failed to unmarshal: %v", err)
	}

	return &infoblox, nil
}

func (ac *AlkiraClient) UpdateInfoblox(id string, in *Infoblox) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/infoblox-services/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(in)

	if err != nil {
		return fmt.Errorf("UpdateInfoblox: failed to marshal request: %v", err)
	}

	return ac.update(uri, body)
}

func (ac *AlkiraClient) DeleteInfoblox(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/infoblox-services/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}
