package alkira

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setPrismaSDWANInstances set instances block values, preserving
// write-only fields from existing state since the API may not return them.
func setPrismaSDWANInstances(d *schema.ResourceData, connector *alkira.ConnectorPrismaSDWAN) {
	var instances []map[string]interface{}

	//
	// Go through all blocks from the config firstly to find a match,
	// instance ID should be uniquely identifying a block.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any block IDs yet.
	//
	for _, inst := range d.Get("instances").([]interface{}) {
		config := inst.(map[string]interface{})

		for _, info := range connector.Instances {
			id, _ := info.Id.Int64()
			if config["id"].(int) == int(id) || config["host_name"].(string) == info.HostName {
				instance := map[string]interface{}{
					"credential_id": info.CredentialId,
					"host_name":     info.HostName,
					"id":            info.Id,
					"ion_model":     info.IonModel,
					"version":       info.Version,
				}
				instances = append(instances, instance)
				break
			}
		}
	}

	//
	// Go through all instances from the API response one more time to
	// find any instances that have not been tracked from Terraform
	// config.
	//
	for _, info := range connector.Instances {
		isNew := true

		// Check if the instance already exists in the Terraform config
		for _, inst := range d.Get("instances").([]interface{}) {
			config := inst.(map[string]interface{})

			id, _ := info.Id.Int64()
			if config["id"].(int) == int(id) || config["host_name"].(string) == info.HostName {
				isNew = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if isNew {
			instance := map[string]interface{}{
				"credential_id": info.CredentialId,
				"host_name":     info.HostName,
				"id":            info.Id,
				"ion_model":     info.IonModel,
				"version":       info.Version,
			}

			instances = append(instances, instance)
			break
		}
	}

	d.Set("instances", instances)
}

// expandPrismaSDWANVrfMappings expand Prisma SD-WAN VRF segment mapping
func expandPrismaSDWANVrfMappings(in *schema.Set) []alkira.ConnectorPrismaSDWANVRFMapping {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty target_segment")
		return []alkira.ConnectorPrismaSDWANVRFMapping{}
	}

	mappings := make([]alkira.ConnectorPrismaSDWANVRFMapping, in.Len())
	for i, mapping := range in.List() {
		r := alkira.ConnectorPrismaSDWANVRFMapping{}
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
		if v, ok := t["vrf_name"].(string); ok {
			r.VrfName = v
		}

		mappings[i] = r
	}

	return mappings
}

// expandPrismaSDWANInstances expand Prisma SD-WAN connector instances
func expandPrismaSDWANInstances(in []interface{}) []alkira.ConnectorPrismaSDWANInstance {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Empty instances")
		return []alkira.ConnectorPrismaSDWANInstance{}
	}

	instances := make([]alkira.ConnectorPrismaSDWANInstance, len(in))

	for i, mapping := range in {
		r := alkira.ConnectorPrismaSDWANInstance{}
		t := mapping.(map[string]interface{})

		if v, ok := t["host_name"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["credential_id"].(string); ok {
			r.CredentialId = v
		}
		if v, ok := t["ion_model"].(string); ok {
			r.IonModel = v
		}
		if v, ok := t["version"].(string); ok {
			r.Version = v
		}
		if v, ok := t["id"].(int); ok {
			r.Id = json.Number(fmt.Sprintf("%d", v))
		} else if v, ok := t["id"].(json.Number); ok {
			r.Id = v
		}

		instances[i] = r
	}

	return instances
}
