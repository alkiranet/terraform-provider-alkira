// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type Group struct {
	Id          json.Number `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
}

func NewGroup(ac *AlkiraClient) *AlkiraAPI[Group] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[Group]{ac, uri, true}
	return api
}
