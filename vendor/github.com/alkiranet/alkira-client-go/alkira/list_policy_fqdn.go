// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PolicyFqdnList struct {
	Id              json.Number `json:"id,omitempty"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Fqdns           []string    `json:"fqdns"`
	DnsServerListId int         `json:"dnsServerListId"`
}

// NewPolicyFqdnList new global cidr list
func NewPolicyFqdnList(ac *AlkiraClient) *AlkiraAPI[PolicyFqdnList] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy-fqdn-lists", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[PolicyFqdnList]{ac, uri, true}
	return api
}
