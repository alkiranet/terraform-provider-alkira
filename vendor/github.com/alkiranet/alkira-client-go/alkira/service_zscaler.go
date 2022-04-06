// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type Zscaler struct {
	BillingTags           []int               `json:"billingTags"`
	Cxp                   string              `json:"cxp"`
	Description           string              `json:"description"`
	Id                    json.Number         `json:"id,omitempty"`           // only set on response
	InternalName          string              `json:"internalName,omitempty"` //only set on response
	IpsecConfiguration    *ZscalerIpSecConfig `json:"ipsecConfiguration"`
	Name                  string              `json:"name"`
	PrimaryPublicEdgeIp   string              `json:"primaryPublicEdgeIp"`
	SecondaryPublicEdgeIp string              `json:"secondaryPublicEdgeIp"`
	Segments              []string            `json:"segments"`
	Size                  string              `json:"size"`
	TunnelType            string              `json:"tunnelType"`
}

type ZscalerIpSecConfig struct {
	EspDhGroupNumber       string `json:"espDhGroupNumber"`
	EspEncryptionAlgorithm string `json:"espEncryptionAlgorithm"`
	EspIntegrityAlgorithm  string `json:"espIntegrityAlgorithm"`
	HealthCheckType        string `json:"healthCheckType"`
	HttpProbeUrl           string `json:"httpProbeUrl"`
	IkeDhGroupNumber       string `json:"ikeDhGroupNumber"`
	IkeEncryptionAlgorithm string `json:"ikeEncryptionAlgorithm"`
	IkeIntegrityAlgorithm  string `json:"ikeIntegrityAlgorithm"`
	LocalFqdnId            string `json:"localFqdnId"`
	PreSharedKey           string `json:"preSharedKey"`
	PingProbeIp            string `json:"pingProbeIp"`
}

func (ac *AlkiraClient) CreateZscaler(z *Zscaler) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/zscaler-internet-access-services", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(z)

	if err != nil {
		return "", fmt.Errorf("CreateZscaler: marshal failed: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result Zscaler
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateZscaler: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

func (ac *AlkiraClient) GetZscalers() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/zscaler-internet-access-services", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ac *AlkiraClient) GetZscalerById(id string) (*Zscaler, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/zscaler-internet-access-services/%s", ac.URI, ac.TenantNetworkId, id)

	var zscaler Zscaler

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &zscaler)

	if err != nil {
		return nil, fmt.Errorf("GetZscalerById: failed to unmarshal: %v", err)
	}

	return &zscaler, nil
}

func (ac *AlkiraClient) UpdateZscaler(id string, z *Zscaler) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/zscaler-internet-access-services/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(z)

	if err != nil {
		return fmt.Errorf("UpdateZscaler: failed to marshal request: %v", err)
	}

	return ac.update(uri, body)
}

func (ac *AlkiraClient) DeleteZscaler(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/zscaler-internet-access-services/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}
