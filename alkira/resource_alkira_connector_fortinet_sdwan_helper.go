package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setWanEdge set wan_edge block values
func setWanEdge(d *schema.ResourceData, connector *alkira.ConnectorFortinetSdwan) {
	var wanEdges []map[string]interface{}

	//
	// Go through all blocks from the config firstly to find a match,
	// WAN edge's ID should be uniquely identifying a block.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any block IDs yet.
	//
	for _, wanEdge := range d.Get("wan_edge").([]interface{}) {
		config := wanEdge.(map[string]interface{})

		for _, info := range connector.Instances {
			if config["id"].(int) == info.Id || config["hostname"].(string) == info.HostName {
				instance := map[string]interface{}{
					"credential_id": info.CredentialId,
					"hostname":      info.HostName,
					"id":            info.Id,
					"license_type":  info.LicenseType,
					"username":      config["username"].(string),
					"password":      config["password"].(string),
					"serial_number": config["serial_number"].(string),
					"version":       config["version"].(string),
				}
				wanEdges = append(wanEdges, instance)
				break
			}
		}
	}

	//
	// Go through all instances from the API response one more time to
	// find any instances that has not been tracked from Terraform
	// config.
	//
	for _, info := range connector.Instances {
		new := true

		// Check if the wan_edge already exists in the Terraform config
		for _, edge := range d.Get("wan_edge").([]interface{}) {
			config := edge.(map[string]interface{})

			if config["id"].(int) == info.Id || config["hostname"].(string) == info.HostName {
				new = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			instance := map[string]interface{}{
				"credential_id": info.CredentialId,
				"hostname":      info.HostName,
				"license_type":  info.LicenseType,
				"serial_number": info.SerialNumber,
				"id":            info.Id,
			}

			wanEdges = append(wanEdges, instance)
			break
		}
	}

	d.Set("wan_edge", wanEdges)
}

// expandFortinetSdwanVrfMappings expand Fortinet SD-WAN VRF segment mapping
func expandFortinetSdwanVrfMappings(in *schema.Set) []alkira.ConnectorFortinetSdwanVrfMapping {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty target_segment")
		return []alkira.ConnectorFortinetSdwanVrfMapping{}
	}

	mappings := make([]alkira.ConnectorFortinetSdwanVrfMapping, in.Len())
	for i, mapping := range in.List() {
		r := alkira.ConnectorFortinetSdwanVrfMapping{}
		t := mapping.(map[string]interface{})

		if v, ok := t["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		if v, ok := t["advertise_default_route"].(bool); ok {
			r.DisableInternetExit = !v
		}
		if v, ok := t["gateway_bgp_asn"].(int); ok {
			r.GatewayBgpAsn = v
		}
		if v, ok := t["segment_id"].(int); ok {
			r.SegmentId = v
		}
		if v, ok := t["vrf_id"].(int); ok {
			r.Vrf = v
		}

		mappings[i] = r
	}

	return mappings
}

// expandFortinetSdwanWanEedges expand WAN edge instances
func expandFortinetSdwanWanEdges(ac *alkira.AlkiraClient, in []interface{}) ([]alkira.ConnectorFortinetSdwanInstance, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Empty wan_edge")
		return []alkira.ConnectorFortinetSdwanInstance{}, nil
	}

	mappings := make([]alkira.ConnectorFortinetSdwanInstance, len(in))

	for i, mapping := range in {
		r := alkira.ConnectorFortinetSdwanInstance{}
		t := mapping.(map[string]interface{})

		var username string
		var password string
		var licenseType string
		var licenseKey string

		if v, ok := t["hostname"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["username"].(string); ok {
			username = v
		}
		if v, ok := t["password"].(string); ok {
			password = v
		}
		if v, ok := t["license_type"].(string); ok {
			licenseType = v
			r.LicenseType = v
		}
		if v, ok := t["license_key"].(string); ok {
			licenseKey = v
		}
		if v, ok := t["serial_number"].(string); ok {
			r.SerialNumber = v
		}
		if v, ok := t["version"].(string); ok {
			r.Version = v
		}
		if v, ok := t["credential_id"].(string); ok {
			if v == "" {
				log.Printf("[DEBUG] Creating Fortinet SD-WAN instance credential")
				credentialName := r.HostName + randomNameSuffix()

				credential := alkira.CredentialFortinetSdwanInstance{
					Username:    username,
					Password:    password,
					LicenseType: licenseType,
					LicenseKey:  licenseKey,
				}

				credentialId, err := ac.CreateCredential(
					credentialName,
					alkira.CredentialTypeFortinetSdwanInstance,
					credential,
					0,
				)

				if err != nil {
					return nil, err
				}

				r.CredentialId = credentialId
			} else {
				r.CredentialId = v
			}
		}
		if v, ok := t["id"].(int); ok {
			r.Id = v
		}

		mappings[i] = r
	}

	return mappings, nil
}
