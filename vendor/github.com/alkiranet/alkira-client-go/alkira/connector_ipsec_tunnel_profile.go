// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorIPSecTunnelProfileIpSecConfiguration struct {
	EncryptionAlgorithm string `json:"encryptionAlgorithm"`
	IntegrityAlgorithm  string `json:"integrityAlgorithm"`
	DhGroup             string `json:"dhGroup"`
}

type ConnectorIPSecTunnelProfileIkeConfiguration struct {
	EncryptionAlgorithm string `json:"encryptionAlgorithm"`
	IntegrityAlgorithm  string `json:"integrityAlgorithm"`
	DhGroup             string `json:"dhGroup"`
}

type ConnectorIPSecTunnelProfile struct {
	Id                 json.Number                                   `json:"id,omitempty"` // response only
	Name               string                                        `json:"name"`
	Description        string                                        `json:"description"`
	IpSecConfiguration ConnectorIPSecTunnelProfileIpSecConfiguration `json:"ipSecConfiguration"`
	IkeConfiguration   ConnectorIPSecTunnelProfileIkeConfiguration   `json:"ikeConfiguration"`
}

// NewConnectorIPSec initialize a new connector
func NewConnectorIPSecTunnelProfile(ac *AlkiraClient) *AlkiraAPI[ConnectorIPSecTunnelProfile] {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ipsec-tunnel-profiles", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorIPSecTunnelProfile]{ac, uri, true}
	return api
}
