// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
)

type ZtaProfile struct {
	Id   string `json:"id,omitempty"` // only set on response
	Name string `json:"name"`
}

// NewZtaProfile new Ztna profile
func NewZtaProfile(ac *AlkiraClient) *AlkiraAPI[ZtaProfile] {
	uri := fmt.Sprintf("%s/zero-trust-access-profiles", ac.URI)
	api := &AlkiraAPI[ZtaProfile]{ac, uri, true}
	return api
}
