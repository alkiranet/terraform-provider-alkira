// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
)

type UserGroup struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewUserGroup new user group
func NewUserGroup(ac *AlkiraClient) *AlkiraAPI[UserGroup] {
	uri := fmt.Sprintf("%s/user-groups", ac.URI)
	api := &AlkiraAPI[UserGroup]{ac, uri, false}
	return api
}
