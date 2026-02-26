package alkira

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deflateArubaEdgeInstances(ins []alkira.ArubaEdgeInstance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, instance := range ins {
		id, _ := instance.Id.Int64()
		i := map[string]interface{}{
			"account_name":  instance.AccountName,
			"credential_id": instance.CredentialId,
			"host_name":     instance.HostName,
			"id":            int(id),
			"name":          instance.Name,
			"site_tag":      instance.SiteTag,
		}
		instances = append(instances, i)
	}

	return instances
}

// setArubaEdgeInstances sets the instances block, preserving write-only fields
// (account_key) from existing state since the API does not return them.
func setArubaEdgeInstances(d *schema.ResourceData, connector *alkira.ConnectorArubaEdge) {
	var instances []map[string]interface{}

	//
	// First pass: match existing state instances to API instances by ID or
	// name, preserving account_key which is not returned by the API.
	//
	for _, inst := range d.Get("instances").([]interface{}) {
		config := inst.(map[string]interface{})

		for _, info := range connector.Instances {
			id, _ := info.Id.Int64()
			if config["id"].(int) == int(id) || config["name"].(string) == info.Name {
				instance := map[string]interface{}{
					"account_key":   config["account_key"].(string),
					"account_name":  info.AccountName,
					"credential_id": info.CredentialId,
					"host_name":     info.HostName,
					"id":            int(id),
					"name":          info.Name,
					"site_tag":      info.SiteTag,
				}
				instances = append(instances, instance)
				break
			}
		}
	}

	//
	// Second pass: find any API instances not present in state (e.g. added
	// outside Terraform) and append them. This will generate a diff.
	//
	for _, info := range connector.Instances {
		isNew := true
		id, _ := info.Id.Int64()

		for _, inst := range d.Get("instances").([]interface{}) {
			config := inst.(map[string]interface{})
			if config["id"].(int) == int(id) || config["name"].(string) == info.Name {
				isNew = false
				break
			}
		}

		if isNew {
			instance := map[string]interface{}{
				"account_name":  info.AccountName,
				"credential_id": info.CredentialId,
				"host_name":     info.HostName,
				"id":            int(id),
				"name":          info.Name,
				"site_tag":      info.SiteTag,
			}
			instances = append(instances, instance)
		}
	}

	d.Set("instances", instances)
}

func expandArubaEdgeInstances(in []interface{}, client *alkira.AlkiraClient) ([]alkira.ArubaEdgeInstance, error) {

	var instances []alkira.ArubaEdgeInstance

	credentialResponse, err := getAllCredentialsAsCredentialResponseDetails(client)
	if err != nil {
		return nil, err
	}

	for _, v := range in {
		var name, accountKey, accountName, hostName, id, siteTag string
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
		if v, ok := m["id"].(int); ok {
			id = strconv.Itoa(v)
		}
		if v, ok := m["account_name"].(string); ok {
			accountName = v
		}
		if v, ok := m["site_tag"].(string); ok {
			siteTag = v
		}

		var credId string
		if existingCredId, ok := m["credential_id"].(string); ok && existingCredId != "" {
			credId = existingCredId
		} else {
			credId, err = findOrCreateArubaEdgeInstanceCredentialByName(client, credentialResponse, name, accountKey)
			if err != nil {
				return nil, err
			}
		}

		c := alkira.ArubaEdgeInstance{
			AccountName:  accountName,
			CredentialId: credId,
			HostName:     hostName,
			Id:           json.Number(id),
			Name:         name,
			SiteTag:      siteTag,
		}

		instances = append(instances, c)
	}

	return instances, nil
}

func deflateArubaEdgeVrfMapping(vrf []alkira.ArubaEdgeVRFMappings) ([]map[string]interface{}, error) {

	var mappings []map[string]interface{}
	for _, vrfmapping := range vrf {

		i := map[string]interface{}{
			"advertise_on_prem_routes":   vrfmapping.AdvertiseOnPremRoutes,
			"segment_id":                 strconv.Itoa(vrfmapping.AlkiraSegmentId),
			"aruba_edge_connect_segment": vrfmapping.ArubaEdgeConnectSegmentName,
			"advertise_default_route":    !vrfmapping.DisableInternetExit,
			"gateway_bgp_asn":            vrfmapping.GatewayBgpAsn,
		}
		mappings = append(mappings, i)
	}

	return mappings, nil
}

func expandArubaEdgeVrfMappings(in *schema.Set) ([]alkira.ArubaEdgeVRFMappings, error) {
	var mappings []alkira.ArubaEdgeVRFMappings
	if in == nil || in.Len() == 0 {
		return nil, errors.New("ERROR: Invalid aruba edge mapping input: Cannot be nil or empty")
	}

	for _, v := range in.List() {
		var arubaEdgeVRFMapping alkira.ArubaEdgeVRFMappings
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
		if v, ok := m["aruba_edge_connect_segment"].(string); ok {
			arubaEdgeVRFMapping.ArubaEdgeConnectSegmentName = v
		}
		if v, ok := m["advertise_default_route"].(bool); ok {
			arubaEdgeVRFMapping.DisableInternetExit = !v
		}
		if v, ok := m["gateway_bgp_asn"].(int); ok {
			arubaEdgeVRFMapping.GatewayBgpAsn = v
		}

		mappings = append(mappings, arubaEdgeVRFMapping)
	}

	return mappings, nil
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
	return client.CreateCredential(name, alkira.CredentialTypeArubaEdgeConnectInstance, alkira.CredentialArubaEdgeConnectInstance{AccountKey: accountKey}, 0)
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
