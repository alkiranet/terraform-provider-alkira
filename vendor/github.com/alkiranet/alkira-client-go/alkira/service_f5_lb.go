package alkira

import (
	"encoding/json"
	"fmt"
)

type ServiceF5Lb struct {
	SegmentOptions      F5SegmentOption `json:"segmentOptions"`
	Description         string          `json:"description,omitempty"`
	Cxp                 string          `json:"cxp"`
	Size                string          `json:"size"`
	Id                  json.Number     `json:"id,omitempty"`
	Name                string          `json:"name"`
	ServiceGroupName    string          `json:"serviceGroupName"`
	IlbServiceGroupName string          `json:"ilbServiceGroupName"`
	ImplicitGroupId     int             `json:"implicitGroupId,omitempty"`    // response only
	IlbImplicitGroupId  int             `json:"ilbImplicitGroupId,omitempty"` // response only
	Instances           []F5Instance    `json:"instances"`
	Segments            []string        `json:"segments"`
	BillingTags         []int           `json:"billingTags,omitempty"`
	PrefixListId        int             `json:"prefixListId,omitempty"`
	GlobalCidrListId    int             `json:"globalCidrListId"`
}

type F5Instance struct {
	Deployment               F5InstanceDeployment         `json:"deployment"`
	Name                     string                       `json:"name"`
	RegistrationCredentialId string                       `json:"registrationCredentialId,omitempty"`
	CredentialId             string                       `json:"credentialId"`
	HostNameFqdn             string                       `json:"hostNameFqdn"`
	LicenseType              string                       `json:"licenseType"`
	Version                  string                       `json:"version"`
	AvailabilityZone         json.Number                  `json:"availabilityZone,omitempty"`
	Id                       int                          `json:"id,omitempty"`       // RESPONSE ONLY
	Metadata                 *F5LBServiceInstanceMetadata `json:"metadata,omitempty"` // RESPONSE ONLY
}

type F5LBServiceInstanceMetadata struct {
	F5MgmtPublicIp    string                                        `json:"f5MgmtPublicIp,omitempty"`
	SegmentToMetadata map[string]F5LBServiceInstanceSegmentMetadata `json:"segmentToMetadata,omitempty"`
}

type F5LBServiceInstanceSegmentMetadata struct {
	SegmentId      int64                       `json:"segmentId,omitempty"`
	RouteDomainId  int64                       `json:"routeDomainId,omitempty"`
	RoutingType    string                      `json:"routingType,omitempty"`
	F5MgmtPublicIp string                      `json:"f5MgmtPublicIp,omitempty"`
	Vlans          []string                    `json:"vlans,omitempty"`
	Tunnels        []F5LBServiceInstanceTunnel `json:"tunnels,omitempty"`
}

type F5LBServiceInstanceTunnel struct {
	TunnelProtocol     string `json:"tunnelProtocol,omitempty"`
	TunnelUUID         string `json:"tunnelUUID,omitempty"`
	TunnelId           string `json:"tunnelId,omitempty"`
	CustomerTunnelName string `json:"customerTunnelName,omitempty"`
	CxpTunnelName      string `json:"cxpTunnelName,omitempty"`
	TunnelInternalName string `json:"tunnelInternalName,omitempty"`
	InfraNodeName      string `json:"infraNodeName,omitempty"`
	CustomerOuterIp    string `json:"customerOuterIp,omitempty"`
	CxpOuterIp         string `json:"cxpOuterIp,omitempty"`
	CxpInnerIp         string `json:"cxpInnerIp,omitempty"`
	CustomerInnerIp    string `json:"customerInnerIp,omitempty"`
	LbType             string `json:"lbType,omitempty"`
	BgpEnabled         *bool  `json:"bgpEnabled,omitempty"`
}

type F5SegmentOption map[string]F5SegmentSubOption
type F5SegmentSubOption struct {
	NatPoolPrefixLength int            `json:"natPoolPrefixLength,omitempty"`
	ElbNicCount         int            `json:"elbNicCount"`
	ElbBgpOptions       *ElbBgpOptions `json:"elbBgpOptions,omitempty"`
	LbType              []string       `json:"lbType,omitempty"`
}

type ElbBgpOptions struct {
	AdvertiseToCXPPrefixListId int `json:"advertiseToCXPPrefixListId,omitempty"`
}

type F5InstanceDeployment struct {
	Option string `json:"option,omitempty"`
	Type   string `json:"type"`
}

func NewServiceF5Lb(ac *AlkiraClient) *AlkiraAPI[ServiceF5Lb] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/f5-lb-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceF5Lb]{ac, uri, true}
	return api
}
