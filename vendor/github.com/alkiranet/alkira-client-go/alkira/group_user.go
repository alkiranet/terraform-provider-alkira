// Copyright (C) 2022-2023 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type UserGroup struct {
	Id          json.Number `json:"id,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
}

// NewUserGroup new user group
func NewUserGroup(ac *AlkiraClient) *AlkiraAPI[UserGroup] {
	uri := fmt.Sprintf("%s/user-groups", ac.URI)
	api := &AlkiraAPI[UserGroup]{ac, uri}
	return api
}
