// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

// ConnectorVnetUdrList is a CIDR row. Value is the match key server-side - TPS
// indexes the cidrs list by it - so it carries no omitempty and there is no Values:
// a CIDR is one prefix. AzureVNETConnector.UDRLists.CIDR on the server does not set
// ignoreUnknown, so an unexpected key here is a 400.
type ConnectorVnetUdrList struct {
	Id         string `json:"id"`
	UdrListIds []int  `json:"udrListIds"`
	Value      string `json:"value"`
}

// AK-74515: an Azure subnet can carry several address prefixes, and is onboarded
// all-or-nothing, so one entry names the subnet and lists every prefix in Values.
//
// A request sets one or the other. A RESPONSE may carry BOTH: the server keeps
// Value as the representative prefix for older clients (AzureVNETConnector.java,
// ImportFromCXPOptions.Subnet). So on the read path do not treat a non-empty Value
// as meaning "single prefix" - check Values first, or a multi-prefix subnet gets
// flattened and written back narrowed. Both carry omitempty because the server
// distinguishes an absent field from an empty one - a "value": "" reaches clients
// that read the field directly and is not the same as the key being missing.
//
// Split from the CIDR row above rather than shared: Values is meaningless on a CIDR
// and the server rejects it, so a separate type makes that unrepresentable instead
// of relying on callers to leave it unset.
type ConnectorVnetUdrListSubnet struct {
	Id         string   `json:"id"`
	UdrListIds []int    `json:"udrListIds"`
	Value      string   `json:"value,omitempty"`
	Values     []string `json:"values,omitempty"`
}

type ConnectorVnetUdrLists struct {
	Cidrs   []ConnectorVnetUdrList       `json:"cidrs"`
	Subnets []ConnectorVnetUdrListSubnet `json:"subnets"`
}

// ConnectorVnetServiceRoute is a CIDR row - see ConnectorVnetUdrList.
type ConnectorVnetServiceRoute struct {
	Id                 string   `json:"id"`
	ServiceTags        []string `json:"serviceTags"`
	NativeServices     []string `json:"nativeServices,omitempty"`
	NativeServiceNames []string `json:"nativeServiceNames,omitempty"`
	Value              string   `json:"value"`
}

// AK-74515: the subnet row - see ConnectorVnetUdrListSubnet.
// No NativeServices: the server's ServiceRoutes.Subnet has only nativeServiceNames
// and sets ignoreUnknown, so a nativeServices key here is silently dropped. The CIDR
// type above carries it for compatibility; this type is new, so it does not.
type ConnectorVnetServiceRouteSubnet struct {
	Id                 string   `json:"id"`
	ServiceTags        []string `json:"serviceTags"`
	NativeServiceNames []string `json:"nativeServiceNames,omitempty"`
	Value              string   `json:"value,omitempty"`
	Values             []string `json:"values,omitempty"`
}

type ConnectorVnetServiceRoutes struct {
	Cidrs   []ConnectorVnetServiceRoute       `json:"cidrs"`
	Subnets []ConnectorVnetServiceRouteSubnet `json:"subnets"`
}

type ConnectorVnetExportOptionUserInputPrefix struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
	// AK-74515: every prefix of a SUBNET entry - see ConnectorVnetUdrListSubnet. The
	// selection is all-or-nothing: the set must name every prefix the subnet has, or
	// the request is rejected.
	//
	// SUBNET entries only. This type is not split by Type because the wire format is
	// one list with a discriminator, so the invalid combination stays representable:
	// on a CIDR entry the server IGNORES Values rather than rejecting it
	// (getSelectedCidrPrefixes maps getValue() alone), and the format check only runs
	// when there is more than one value - so Type: "CIDR" with only Values passes
	// validation and reaches the CIDR branch with an empty prefix. Set Value, not
	// Values, on a CIDR entry.
	Values []string `json:"values,omitempty"`
}

type ConnectorVnetExportOptions struct {
	UserInputPrefixes []ConnectorVnetExportOptionUserInputPrefix `json:"userInputPrefixes"`
}

type ConnectorVnetImportOptionsCidr struct {
	RouteImportMode string `json:"routeImportMode"`
	PrefixListIds   []int  `json:"prefixListIds"`
	Value           string `json:"value"`
}

type ConnectorVnetImportOptionsSubnet struct {
	Id              string `json:"id"`
	RouteImportMode string `json:"routeImportMode"`
	PrefixListIds   []int  `json:"prefixListIds"`
	Value           string `json:"value,omitempty"`
	// AK-74515: see ConnectorVnetUdrList.
	Values []string `json:"values,omitempty"`
}

type ConnectorVnetImportOptions struct {
	Cidrs           []ConnectorVnetImportOptionsCidr   `json:"cidrs,omitempty"`
	PrefixListIds   []int                              `json:"prefixListIds,omitempty"`
	RouteImportMode string                             `json:"routeImportMode"`
	Subnets         []ConnectorVnetImportOptionsSubnet `json:"subnets,omitempty"`
}

type ConnectorVnetRouting struct {
	ExportOptions ConnectorVnetExportOptions `json:"exportToCXPOptions,omitempty"`
	ImportOptions ConnectorVnetImportOptions `json:"importFromCXPOptions"`
	ServiceRoutes ConnectorVnetServiceRoutes `json:"serviceRoutes,omitempty"`
	UdrLists      ConnectorVnetUdrLists      `json:"udrLists,omitempty"`
}

type ConnectorAzureVnet struct {
	BillingTags                       []int                 `json:"billingTags"`
	CXP                               string                `json:"cxp"`
	CredentialId                      string                `json:"credentialId"`
	Group                             string                `json:"group,omitempty"`
	Enabled                           bool                  `json:"enabled"`
	Id                                json.Number           `json:"id,omitempty"`              // RESPONSE ONLY
	ImplicitGroupId                   int                   `json:"implicitGroupId,omitempty"` // RESPONSE ONLY
	Name                              string                `json:"name"`
	NativeServiceNames                []string              `json:"nativeServiceNames,omitempty"`
	ResourceGroupName                 string                `json:"resourceGroupName,omitempty"`
	SecondaryCXPs                     []string              `json:"secondaryCXPs,omitempty"`
	Segments                          []string              `json:"segments"`
	ServiceTags                       []string              `json:"serviceTags,omitempty"`
	Size                              string                `json:"size"`
	VnetId                            string                `json:"vnetId"`
	ConnectionMode                    string                `json:"connectionMode,omitempty"`
	VnetRouting                       *ConnectorVnetRouting `json:"vnetRouting"`
	CustomerASN                       int                   `json:"customerAsn,omitempty"`
	ScaleGroupId                      string                `json:"scaleGroupId,omitempty"`
	PeeringGatewayCxpId               int                   `json:"cxpPeeringGatewayId,omitempty"`
	DirectInterVNETCommunicationGroup string                `json:"directInterVnetCommunicationGroup,omitempty"`
	UdrListIds                        []int                 `json:"udrListIds,omitempty"`
	Description                       string                `json:"description,omitempty"`
}

// NewConnectorAzureVnet initialize a new connector
func NewConnectorAzureVnet(ac *AlkiraClient) *AlkiraAPI[ConnectorAzureVnet] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAzureVnet]{ac, uri, true}
	return api
}
