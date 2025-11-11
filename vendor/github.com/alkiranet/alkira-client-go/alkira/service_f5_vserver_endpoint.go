package alkira

import (
	"encoding/json"
	"fmt"
)

type F5vServerEndpoint struct {
	F5ServiceId          int         `json:"f5ServiceId"`
	Name                 string      `json:"name"`
	Id                   json.Number `json:"id,omitempty"` // RESPONSE ONLY
	Type                 string      `json:"type"`
	Segment              string      `json:"segment"`
	FqdnPrefix           string      `json:"fqdnPrefix"`
	Protocol             string      `json:"protocol"`
	Snat                 string      `json:"snat"`
	F5ServiceInstanceIds []int       `json:"f5ServiceInstanceIds,omitempty"`
	PortRanges           []string    `json:"portRanges"`
}

func NewF5vServerEndpoint(ac *AlkiraClient) *AlkiraAPI[F5vServerEndpoint] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/f5-vserver-endpoints", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[F5vServerEndpoint]{ac, uri, true}
	return api

}
