package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setJuniperInstances(d *schema.ResourceData, connector *alkira.ConnectorJuniperSdwan) {
	var instances []map[string]interface{}

	//
	// Go through all blocks from the config firstly to find a match,
	// Juniper Instance's ID should be uniquely identifying a block.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any block IDs yet.
	//
	for _, instance := range d.Get("instance").([]interface{}) {
		config := instance.(map[string]interface{})

		for _, info := range connector.Instances {
			if config["id"].(int) == info.Id || config["hostname"].(string) == info.HostName {
				instance := map[string]interface{}{
					"hostname":                       info.HostName,
					"id":                             info.Id,
					"registration_key":               config["registration_key"].(string),
					"registration_key_credential_id": info.RegistrationKeyCredentialId,
				}
				instances = append(instances, instance)
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

		// Check if the instance already exists in the Terraform config
		for _, instance := range d.Get("instance").([]interface{}) {
			config := instance.(map[string]interface{})

			if config["id"].(int) == info.Id || config["hostname"].(string) == info.HostName {
				new = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			instance := map[string]interface{}{
				"hostname":                       info.HostName,
				"id":                             info.Id,
				"registration_key_credential_id": info.RegistrationKeyCredentialId,
			}

			instances = append(instances, instance)
			break
		}
	}

	d.Set("instance", instances)
}

// expandJuniperSdwanVrfMappings expand Juniper SD-WAN VRF Mapping
func expandJuniperSdwanVrfMappings(in *schema.Set) []alkira.ConnectorJuniperSsrVrfMapping {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty SSR Vrf Mapping")
		return []alkira.ConnectorJuniperSsrVrfMapping{}
	}

	mappings := make([]alkira.ConnectorJuniperSsrVrfMapping, in.Len())
	for i, mapping := range in.List() {
		r := alkira.ConnectorJuniperSsrVrfMapping{}
		t := mapping.(map[string]interface{})

		if v, ok := t["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		if v, ok := t["advertise_default_route"].(bool); ok {
			r.DisableInternetExit = !v
		}
		if v, ok := t["segment_id"].(int); ok {
			r.SegmentId = v
		}
		r.JuniperSsrBgpAsn = 65000
		r.JuniperSsrVrfName = "default"
		mappings[i] = r
	}

	return mappings
}

// expandJuniperSdwanWanInstances expand Juniper SD-WAN Instances
func expandJuniperSdwanInstances(ac *alkira.AlkiraClient, in []interface{}) ([]alkira.ConnectorJuniperSdwanInstance, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Empty instance")
		return []alkira.ConnectorJuniperSdwanInstance{}, nil
	}

	mappings := make([]alkira.ConnectorJuniperSdwanInstance, len(in))

	for i, mapping := range in {
		r := alkira.ConnectorJuniperSdwanInstance{}
		t := mapping.(map[string]interface{})

		var registrationKey string

		if v, ok := t["hostname"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["registration_key"].(string); ok {
			registrationKey = v
		}
		if v, ok := t["registration_key_credential_id"].(string); ok {
			if v == "" {
				log.Printf("[DEBUG] Creating Juniper SD-WAN instance registration key credential")
				credentialName := r.HostName + randomNameSuffix()

				credential := alkira.CredentialApiKey{
					ApiKey: registrationKey,
				}

				credentialId, err := ac.CreateSingleUseCredential(
					credentialName,
					alkira.CredentialTypeApiKey,
					credential,
					0,
				)

				if err != nil {
					return nil, err
				}

				r.RegistrationKeyCredentialId = credentialId
			} else {
				r.RegistrationKeyCredentialId = v
			}
		}

		if v, ok := t["id"].(int); ok {
			r.Id = v
		}

		mappings[i] = r
	}

	return mappings, nil
}
