package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setVedge set virtual edge block values
func setVirtualEdge(d *schema.ResourceData, connector *alkira.ConnectorVmwareSdwan) {
	var vedges []map[string]interface{}

	//
	// Go through all vedge blocks from the config firstly to find a
	// match, vedge's ID should be uniquely identifying an vedge
	// block.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any vedge IDs yet.
	//
	for _, vedge := range d.Get("virtual_edge").([]interface{}) {
		vedgeConfig := vedge.(map[string]interface{})

		for _, info := range connector.Instance {
			if vedgeConfig["id"].(int) == info.Id || vedgeConfig["hostname"].(string) == info.HostName {
				vedge := map[string]interface{}{
					"credential_id": info.CredentialId,
					"hostname":      info.HostName,
					"id":            info.Id,
					"username":      vedgeConfig["username"].(string),
					"password":      vedgeConfig["password"].(string),
				}
				vedges = append(vedges, vedge)
				break
			}
		}
	}

	//
	// Go through all CiscoEdgeInfo from the API response one more
	// time to find any vedge that has not been tracked from Terraform
	// config.
	//
	for _, info := range connector.Instance {
		new := true

		// Check if the vedge already exists in the Terraform config
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
