// Copyright (C) 2022-2023 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ByoipExtraAttributes struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

type Byoip struct {
	Prefix          string               `json:"prefix"`
	Cxp             string               `json:"cxp"`
	Description     string               `json:"description"`
	ExtraAttributes ByoipExtraAttributes `json:"extraAttributes"`
	DoNotAdvertise  bool                 `json:"doNotAdvertise"`
	Id              json.Number          `json:"id,omitempty"` // response only
}

// NewByoip new BYOIP
func NewByoip(ac *AlkiraClient) *AlkiraAPI[Byoip] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/byoips", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[Byoip]{ac, uri}
	return api
}
