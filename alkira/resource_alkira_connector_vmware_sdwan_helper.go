package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setVirtualVedge set virtual edge block values
func setVirtualEdge(d *schema.ResourceData, connector *alkira.ConnectorVmwareSdwan) {
	var vedges []map[string]interface{}

	//
	// Go through all blocks from the config firstly to find a match,
	// virtual vedge's ID should be uniquely identifying a block.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any block IDs yet.
	//
	for _, vedge := range d.Get("virtual_edge").([]interface{}) {
		vedgeConfig := vedge.(map[string]interface{})

		for _, info := range connector.Instances {
			if vedgeConfig["id"].(int) == info.Id || vedgeConfig["hostname"].(string) == info.HostName {
				vedge := map[string]interface{}{
					"credential_id":   info.CredentialId,
					"hostname":        info.HostName,
					"id":              info.Id,
					"activation_code": vedgeConfig["activation_code"].(string),
				}
				vedges = append(vedges, vedge)
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

		// Check if the virtual_vedge already exists in the Terraform config
		for _, vedge := range d.Get("virtual_edge").([]interface{}) {
			vedgeConfig := vedge.(map[string]interface{})

			if vedgeConfig["id"].(int) == info.Id || vedgeConfig["hostname"].(string) == info.HostName {
				new = false
				break
			}
		}

		// If the vedge is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			vedge := map[string]interface{}{
				"credential_id": info.CredentialId,
				"hostname":      info.HostName,
				"id":            info.Id,
			}

			vedges = append(vedges, vedge)
			break
		}
	}

	d.Set("virtual_edge", vedges)
}

// expandVmwareSdwanVrfMappings expand VMWARE SD-WAN VRF segment mapping
func expandVmwareSdwanVrfMappings(in *schema.Set) []alkira.VmwareSdwanVrfMapping {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty target_segment")
		return []alkira.VmwareSdwanVrfMapping{}
	}

	mappings := make([]alkira.VmwareSdwanVrfMapping, in.Len())
	for i, mapping := range in.List() {
		r := alkira.VmwareSdwanVrfMapping{}
		t := mapping.(map[string]interface{})

		if v, ok := t["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		if v, ok := t["allow_nat_exit"].(bool); ok {
			r.DisableInternetExit = !v
		}
		if v, ok := t["gateway_bgp_asn"].(int); ok {
			r.GatewayBgpAsn = v
		}
		if v, ok := t["segment_id"].(int); ok {
			r.SegmentId = v
		}
		if v, ok := t["vmware_sdwan_segment_name"].(string); ok {
			r.VmWareSdWanSegmentName = v
		}

		mappings[i] = r
	}

	return mappings
}

// expandVmwareSdwanVedges expand virtual edges
func expandVmwareSdwanVirtualEdges(ac *alkira.AlkiraClient, in []interface{}) ([]alkira.VmwareSdwanInstance, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Empty virtual_edge")
		return []alkira.VmwareSdwanInstance{}, nil
	}

	mappings := make([]alkira.VmwareSdwanInstance, len(in))

	for i, mapping := range in {
		r := alkira.VmwareSdwanInstance{}
		t := mapping.(map[string]interface{})

		var activationCode string

		if v, ok := t["hostname"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["activation_code"].(string); ok {
			activationCode = v
		}
		if v, ok := t["credential_id"].(string); ok {
			if v == "" {
				log.Printf("[DEBUG] Creating VMWARE-SDWAN Instance Credential")
				credentialName := r.HostName + randomNameSuffix()

				credential := alkira.CredentialVmwareSdwanInstance{
					ActivationCode: activationCode,
				}

				credentialId, err := ac.CreateCredential(
					credentialName,
					alkira.CredentialTypeVmwareSdwanInstance,
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
