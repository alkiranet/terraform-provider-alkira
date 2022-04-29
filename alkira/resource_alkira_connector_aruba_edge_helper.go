package alkira

import (
	"errors"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deflateArubaEdgeInstances(ins []alkira.ArubaEdgeInstance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, instance := range ins {
		i := map[string]interface{}{
			"account_name":  instance.AccountName,
			"credential_id": instance.CredentialId,
			"host_name":     instance.HostName,
			"name":          instance.Name,
			"site_tag":      instance.SiteTag,
		}
		instances = append(instances, i)
	}

	return instances
}

func expandArubaEdgeInstances(in *schema.Set, client *alkira.AlkiraClient) ([]alkira.ArubaEdgeInstance, error) {

	var instances []alkira.ArubaEdgeInstance

	credentialResponse, err := getAllCredentialsAsCredentialResponseDetails(client)
	if err != nil {
		return nil, err
	}

	for _, v := range in.List() {
		var name, accountKey, accountName, hostName, siteTag string
		m := v.(map[string]interface{})

		if v, ok := m["account_key"].(string); ok {
			accountKey = v
		}
		if v, ok := m["name"].(string); ok {
			name = v
		}
		if v, ok := m["host_name"].(string); ok {
			hostName = v
		}
		if v, ok := m["account_name"].(string); ok {
			accountName = v
		}
		if v, ok := m["site_tag"].(string); ok {
			siteTag = v
		}

		credId, err := findOrCreateArubaEdgeInstanceCredentialByName(client, credentialResponse, name, accountKey)
		if err != nil {
			return nil, err
		}

		c := alkira.ArubaEdgeInstance{
			AccountName:  accountName,
			CredentialId: credId,
			HostName:     hostName,
			Name:         name,
			SiteTag:      siteTag,
		}

		instances = append(instances, c)
	}

	return instances, nil
}

func deflateArubaEdgeVrfMapping(vrf []alkira.ArubaEdgeVRFMapping) []map[string]interface{} {
	var mappings []map[string]interface{}

	for _, vrfmapping := range vrf {
		i := map[string]interface{}{
			"advertise_on_prem_routes":        vrfmapping.AdvertiseOnPremRoutes,
			"segment_id":                      strconv.Itoa(vrfmapping.AlkiraSegmentId),
			"aruba_edge_connect_segment_name": vrfmapping.ArubaEdgeConnectSegmentName,
			"disable_internet_exit":           vrfmapping.DisableInternetExit,
			"gateway_gbp_asn":                 vrfmapping.GatewayBgpAsn,
		}
		mappings = append(mappings, i)
	}

	return mappings
}

func expandArubeEdgeVrfMapping(in *schema.Set) ([]alkira.ArubaEdgeVRFMapping, error) {
	var mappings []alkira.ArubaEdgeVRFMapping

	if in == nil || in.Len() == 0 {
		return nil, errors.New("Invalid aruba edge mapping input: Cannot be nil or empty.")
	}

	for _, v := range in.List() {
		var arubaEdgeVRFMapping alkira.ArubaEdgeVRFMapping
		m := v.(map[string]interface{})

		if v, ok := m["advertise_on_prem_routes"].(bool); ok {
			arubaEdgeVRFMapping.AdvertiseOnPremRoutes = v
		}
		if v, ok := m["segment_id"].(string); ok {
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			arubaEdgeVRFMapping.AlkiraSegmentId = i
		}
		if v, ok := m["aruba_edge_connect_segment_name"].(string); ok {
			arubaEdgeVRFMapping.ArubaEdgeConnectSegmentName = v
		}
		if v, ok := m["disable_internet_exit"].(bool); ok {
			arubaEdgeVRFMapping.DisableInternetExit = v
		}
		if v, ok := m["gateway_gbp_asn"].(int); ok {
			arubaEdgeVRFMapping.GatewayBgpAsn = v
		}

		mappings = append(mappings, arubaEdgeVRFMapping)
	}

	return mappings, nil
}

func setArubaEdgeResourceFields(connector *alkira.ConnectorArubaEdge, d *schema.ResourceData) {
	d.Set("aruba_edge_vrf_mapping", deflateArubaEdgeVrfMapping(connector.ArubaEdgeVrfMapping))
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("boost_mode", connector.BoostMode)
	d.Set("cxp", connector.Cxp)
	d.Set("gateway_gbp_asn", connector.GatewayBgpAsn)
	d.Set("group", connector.Group)
	d.Set("instances", deflateArubaEdgeInstances(connector.Instances))
	d.Set("name", connector.Name)
	d.Set("segment_names", connector.Segments)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("version", connector.Version)
}

func findArubaEdgeInstanceResponseDetailByName(credentials []alkira.CredentialResponseDetail, name string) *alkira.CredentialResponseDetail {
	for _, c := range credentials {
		if name == c.Name {
			return &c
		}
	}

	return nil
}

func createArubaEdgeInstanceCredential(client *alkira.AlkiraClient, name, accountKey string) (string, error) {
	return client.CreateCredential(name, alkira.CredentialTypeArubaEdgeConnectInstance, alkira.CredentialArubaEdgeConnectInstance{AccountKey: accountKey})
}

func findOrCreateArubaEdgeInstanceCredentialByName(client *alkira.AlkiraClient, credentials []alkira.CredentialResponseDetail, name, accountKey string) (string, error) {
	credential := findArubaEdgeInstanceResponseDetailByName(credentials, name)

	//If credential is not found in existing set create a new one
	if credential == nil {
		newCredentialId, err := createArubaEdgeInstanceCredential(client, name, accountKey)
		if err != nil {
			return "", err
		}

		return newCredentialId, nil
	}

	return credential.Id, nil
}
