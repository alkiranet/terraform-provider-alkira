package alkira

import (
	"log"

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

func expandArubaEdgeInstances(in *schema.Set) []alkira.ArubaEdgeInstance {
	var instances []alkira.ArubaEdgeInstance

	for _, v := range in.List() {
		var instance alkira.ArubaEdgeInstance
		m := v.(map[string]interface{})

		if v, ok := m["account_name"].(string); ok {
			instance.AccountName = v
		}
		if v, ok := m["credential_id"].(string); ok {
			instance.CredentialId = v
		}
		if v, ok := m["host_name"].(string); ok {
			instance.HostName = v
		}
		if v, ok := m["name"].(string); ok {
			instance.Name = v
		}
		if v, ok := m["site_tag"].(string); ok {
			instance.SiteTag = v
		}

		instances = append(instances, instance)
	}

	return instances
}

func deflateArubaEdgeVrfMapping(vrf []alkira.ArubaEdgeVRFMapping) []map[string]interface{} {
	var mappings []map[string]interface{}

	for _, vrfmapping := range vrf {
		i := map[string]interface{}{
			"advertise_on_prem_routes":        vrfmapping.AdvertiseOnPremRoutes,
			"alkira_segment_id":               vrfmapping.AlkiraSegmentId,
			"aruba_edge_connect_segment_name": vrfmapping.ArubaEdgeConnectSegmentName,
			"disable_internet_exit":           vrfmapping.DisableInternetExit,
			"gateway_gbp_asn":                 vrfmapping.GatewayBgpAsn,
		}
		mappings = append(mappings, i)
	}

	return mappings
}

func expandArubeEdgeVrfMapping(in *schema.Set) []alkira.ArubaEdgeVRFMapping {
	var mappings []alkira.ArubaEdgeVRFMapping

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid aruba edge mapping input")
		return nil
	}

	for _, v := range in.List() {
		var arubaEdgeVRFMapping alkira.ArubaEdgeVRFMapping
		m := v.(map[string]interface{})

		if v, ok := m["advertise_on_prem_routes"].(bool); ok {
			arubaEdgeVRFMapping.AdvertiseOnPremRoutes = v
		}
		if v, ok := m["alkira_segment_id"].(int); ok {
			arubaEdgeVRFMapping.AlkiraSegmentId = v
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

	return mappings
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
