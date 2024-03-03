package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// generateConnectorVersaSdwanRequest generate request for Versa SD-WAN connector
func generateConnectorVersaSdwanRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorVersaSdwan, error) {

	// Expand Versa SDWAN VOS devices block
	instances, err := expandVersaSdwanVosDevices(m.(*alkira.AlkiraClient),
		d.Get("versa_vos_device").([]interface{}))

	if err != nil {
		return nil, err
	}

	// Construct the request payload
	connector := &alkira.ConnectorVersaSdwan{
		BillingTags:           convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Cxp:                   d.Get("cxp").(string),
		Group:                 d.Get("group").(string),
		Enabled:               d.Get("enabled").(bool),
		Instances:             instances,
		Name:                  d.Get("name").(string),
		GlobalTenantId:        d.Get("global_tenant_id").(int),
		LocalId:               d.Get("local_id").(string),
		LocalPublicSharedKey:  d.Get("local_public_shared_key").(string),
		RemoteId:              d.Get("remote_id").(string),
		RemotePublicSharedKey: d.Get("remote_public_shared_key").(string),
		Size:                  d.Get("size").(string),
		TunnelProtocol:        d.Get("tunnel_protocol").(string),
		VersaControllerHost:   d.Get("versa_controller_host").(string),
		VersaSdWanVRFMappings: expandVersaSdwanVrfMappings(d.Get("vrf_segment_mapping").(*schema.Set)),
	}

	return connector, nil
}

// expandVersaSdwanVrfMappings expand Versa SD-WAN VRF segment mapping
func expandVersaSdwanVrfMappings(in *schema.Set) []alkira.VersaSdwanVrfMapping {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty vrf_segment_mapping")
		return []alkira.VersaSdwanVrfMapping{}
	}

	mappings := make([]alkira.VersaSdwanVrfMapping, in.Len())
	for i, mapping := range in.List() {
		r := alkira.VersaSdwanVrfMapping{}
		t := mapping.(map[string]interface{})

		if v, ok := t["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		if v, ok := t["advertise_default_route"].(bool); ok {
			r.DisableInternetExit = !v
		}
		if v, ok := t["versa_bgp_asn"].(int); ok {
			r.GatewayBgpAsn = v
		}
		if v, ok := t["segment_id"].(int); ok {
			r.SegmentId = v
		}
		if v, ok := t["vrf_name"].(string); ok {
			r.VrfName = v
		}

		mappings[i] = r
	}

	return mappings
}

// expandVersaSdwanVosDevices expand Versa SD-WAN VOS devices
func expandVersaSdwanVosDevices(ac *alkira.AlkiraClient, in []interface{}) ([]alkira.VersaSdwanInstance, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Empty VOS Devices")
		return []alkira.VersaSdwanInstance{}, nil
	}

	mappings := make([]alkira.VersaSdwanInstance, len(in))

	for i, mapping := range in {
		r := alkira.VersaSdwanInstance{}
		t := mapping.(map[string]interface{})

		if v, ok := t["hostname"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["id"].(int); ok {
			r.Id = v
		}
		if v, ok := t["local_device_serial_number"].(string); ok {
			r.SerialNumber = v
		}
		if v, ok := t["version"].(string); ok {
			r.Version = v
		}

		mappings[i] = r
	}

	return mappings, nil
}

// setVersaSdwanInstance set Versa SDWAN instance block values
func setVersaSdwanInstance(d *schema.ResourceData, connector *alkira.ConnectorVersaSdwan) {
	var vosDevices []map[string]interface{}

	//
	// Go through all vedge blocks from the config firstly to find a
	// match, vedge's ID should be uniquely identifying an vedge
	// block.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any vedge IDs yet.
	//
	for _, vos := range d.Get("versa_vos_device").([]interface{}) {
		vosConfig := vos.(map[string]interface{})

		for _, info := range connector.Instances {
			if vosConfig["id"].(int) == info.Id || vosConfig["hostname"].(string) == info.HostName {
				instance := map[string]interface{}{
					"hostname":                   info.HostName,
					"id":                         info.Id,
					"local_device_serial_number": info.SerialNumber,
					"vesion":                     info.Version,
				}
				vosDevices = append(vosDevices, instance)
				break
			}
		}
	}

	//
	// Go through all VersaSdwanInstance from the API response one more
	// time to find any instance that has not been tracked from Terraform
	// config.
	//
	for _, info := range connector.Instances {
		new := true

		// Check if the instance already exists in the Terraform config
		for _, instance := range d.Get("versa_vos_device").([]interface{}) {
			cfg := instance.(map[string]interface{})

			if cfg["id"].(int) == info.Id || cfg["hostname"].(string) == info.HostName {
				new = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			instance := map[string]interface{}{
				"hostname":                   info.HostName,
				"id":                         info.Id,
				"local_device_serial_number": info.SerialNumber,
				"version":                    info.Version,
			}

			vosDevices = append(vosDevices, instance)
			break
		}
	}

	d.Set("versa_vos_device", vosDevices)
}
