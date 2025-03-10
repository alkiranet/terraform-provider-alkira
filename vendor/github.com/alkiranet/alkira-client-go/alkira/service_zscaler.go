// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ServiceZscaler struct {
	BillingTags           []int               `json:"billingTags"`
	Cxp                   string              `json:"cxp"`
	Description           string              `json:"description"`
	Id                    json.Number         `json:"id,omitempty"`           // only set on response
	InternalName          string              `json:"internalName,omitempty"` // only set on response
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
	PingProbeIp            string `json:"pingProbeIp,omitempty"`
}

// NewServiceZscaler new service zscaler
func NewServiceZscaler(ac *AlkiraClient) *AlkiraAPI[ServiceZscaler] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/zscaler-internet-access-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceZscaler]{ac, uri, true}
	return api
}
