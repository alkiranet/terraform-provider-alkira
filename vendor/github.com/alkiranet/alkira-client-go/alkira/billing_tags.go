// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type BillingTag struct {
	Id          json.Number `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
}

// NewBillingTag
func NewBillingTag(ac *AlkiraClient) *AlkiraAPI[BillingTag] {
	uri := fmt.Sprintf("%s/tags", ac.URI)
	api := &AlkiraAPI[BillingTag]{ac, uri, false}

	return api
}
