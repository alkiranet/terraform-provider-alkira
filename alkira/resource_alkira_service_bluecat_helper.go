package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandBluecatInstances(in []interface{}, m interface{}) ([]alkira.BluecatInstance, error) {
	if in == nil || len(in) == 0 {
		return nil, fmt.Errorf("[ERROR]: Bluecat instances cannot be nil or empty")
	}

	instances := make([]alkira.BluecatInstance, len(in))
	for i, instance := range in {
		var r alkira.BluecatInstance

		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["type"].(string); ok {
			r.Type = v
		}

		// Handle BDDS options
		if bddsOptions, ok := instanceCfg["bdds_options"].([]interface{}); ok && len(bddsOptions) > 0 {
			bddsOpt, err := expandBDDSOptions(bddsOptions, m)
			if err != nil {
				return nil, err
			}
			r.BddsOptions = bddsOpt
		}

		// Handle Edge options
		if edgeOptions, ok := instanceCfg["edge_options"].([]interface{}); ok && len(edgeOptions) > 0 {
			edgeOpt, err := expandEdgeOptions(edgeOptions, m)
			if err != nil {
				return nil, err
			}
			r.EdgeOptions = edgeOpt
		}

		instances[i] = r
	}

	return instances, nil
}

func expandBDDSOptions(in []interface{}, m interface{}) (*alkira.BDDSOptions, error) {
	client := m.(*alkira.AlkiraClient)

	if len(in) == 0 {
		return nil, nil
	}

	cfg := in[0].(map[string]interface{})
	options := &alkira.BDDSOptions{}

	var clientId string
	var activationKey string
	if v, ok := cfg["hostname"].(string); ok {
		options.HostName = v
	}
	if v, ok := cfg["model"].(string); ok {
		options.Model = v
	}
	if v, ok := cfg["version"].(string); ok {
		options.Version = v
	}
	if v, ok := cfg["client_id"].(string); ok {
		clientId = v
	}
	if v, ok := cfg["activation_key"].(string); ok {
		activationKey = v
	}
	if v, ok := cfg["license_credential_id"].(string); ok {
		if v == "" {
			licenseCredentialId, err := client.CreateCredential(
				"bluecat-bdds-" + randomNameSuffix(),
				alkira.CredentialTypeBluecatBDDSInstanceLicense,
				&alkira.CredentialBluecatBDDSInstanceLicense{
					ClientId:      clientId,
					ActivationKey: activationKey,
				},
				0,
			)

			if err != nil {
				return nil, err
			}

			options.LicenseCredentialId = licenseCredentialId
		}

		if v != "" {
			options.LicenseCredentialId = v
		}
	}

	return options, nil
}

func expandEdgeOptions(in []interface{}, m interface{}) (*alkira.EdgeOptions, error) {
	client := m.(*alkira.AlkiraClient)

	if len(in) == 0 {
		return nil, nil
	}

	cfg := in[0].(map[string]interface{})
	options := &alkira.EdgeOptions{}

	var configData string
	if v, ok := cfg["hostname"].(string); ok {
		options.HostName = v
	}
	if v, ok := cfg["version"].(string); ok {
		options.Version = v
	}
	if v, ok := cfg["config_data"].(string); ok {
		configData = v
	}
	if v, ok := cfg["credential_id"].(string); ok {
		if v == "" {
			credentialId, err := client.CreateCredential(
				"bluecat-edge-" + randomNameSuffix(),
				alkira.CredentialTypeBluecatEdgeInstance,
				&alkira.CredentialBluecatEdgeInstance{
					ConfigData: configData,
				},
				0,
			)

			if err != nil {
				return nil, err
			}

			options.CredentialId = credentialId
		}

		if v != "" {
			options.CredentialId = v
		}
	}

	return options, nil
}

func expandBluecatAnycast(in *schema.Set) (*alkira.BluecatAnycast, error) {
	if in == nil || in.Len() == 0 {
		return &alkira.BluecatAnycast{}, nil
	}

	anycast := &alkira.BluecatAnycast{}

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})
		if v, ok := cfg["ips"].([]interface{}); ok {
			anycast.Ips = convertTypeListToStringList(v)
		}
		if v, ok := cfg["backup_cxps"].([]interface{}); ok {
			anycast.BackupCxps = convertTypeListToStringList(v)
		}
	}
	return anycast, nil
}

func deflateBluecatInstances(c []alkira.BluecatInstance, d *schema.ResourceData) []map[string]interface{} {
	var m []map[string]interface{}

	// Read existing instances from state to preserve sensitive fields
	// not returned by the API.
	oldInstances := d.Get("instance").([]interface{})

	for _, v := range c {
		j := map[string]interface{}{
			"id":   v.Id,
			"name": v.Name,
			"type": v.Type,
		}

		// Find matching instance in state by id or name to preserve
		// sensitive fields not returned by the API.
		var oldInstance map[string]interface{}
		for _, value := range oldInstances {
			cfg := value.(map[string]interface{})

			if cfg["id"].(int) == v.Id && v.Id != 0 {
				oldInstance = cfg
				break
			}

			if cfg["name"].(string) == v.Name && v.Name != "" {
				oldInstance = cfg
				break
			}

			// When id and name are not yet set (first apply),
			// match by hostname from bdds_options or edge_options.
			oldHostname := getHostnameFromInstance(cfg)
			newHostname := getHostnameFromBluecatInstance(v)
			if oldHostname != "" && oldHostname == newHostname {
				oldInstance = cfg
				break
			}
		}

		if v.BddsOptions != nil {
			bddsMap := map[string]interface{}{
				"hostname":              v.BddsOptions.HostName,
				"model":                 v.BddsOptions.Model,
				"version":               v.BddsOptions.Version,
				"license_credential_id": v.BddsOptions.LicenseCredentialId,
			}

			// Preserve client_id and activation_key from state since
			// the API does not return these sensitive fields.
			if oldInstance != nil {
				if oldBdds, ok := oldInstance["bdds_options"]; ok {
					oldBddsList := oldBdds.([]interface{})
					if len(oldBddsList) > 0 {
						oldBddsMap := oldBddsList[0].(map[string]interface{})
						bddsMap["client_id"] = oldBddsMap["client_id"]
						bddsMap["activation_key"] = oldBddsMap["activation_key"]
					}
				}
			}

			j["bdds_options"] = []interface{}{bddsMap}
		}

		if v.EdgeOptions != nil {
			edgeMap := map[string]interface{}{
				"hostname":      v.EdgeOptions.HostName,
				"version":       v.EdgeOptions.Version,
				"credential_id": v.EdgeOptions.CredentialId,
			}

			// Preserve config_data from state since the API does not
			// return this field.
			if oldInstance != nil {
				if oldEdge, ok := oldInstance["edge_options"]; ok {
					oldEdgeList := oldEdge.([]interface{})
					if len(oldEdgeList) > 0 {
						oldEdgeMap := oldEdgeList[0].(map[string]interface{})
						edgeMap["config_data"] = oldEdgeMap["config_data"]
					}
				}
			}

			j["edge_options"] = []interface{}{edgeMap}
		}

		m = append(m, j)
	}

	return m
}

// getHostnameFromBluecatInstance extracts the hostname from either
// BddsOptions or EdgeOptions of an API-returned instance.
func getHostnameFromBluecatInstance(instance alkira.BluecatInstance) string {
	if instance.BddsOptions != nil {
		return instance.BddsOptions.HostName
	}
	if instance.EdgeOptions != nil {
		return instance.EdgeOptions.HostName
	}
	return ""
}

// getHostnameFromInstance extracts the hostname from either
// bdds_options or edge_options of a state instance.
func getHostnameFromInstance(cfg map[string]interface{}) string {
	if bdds, ok := cfg["bdds_options"].([]interface{}); ok && len(bdds) > 0 {
		if opts, ok := bdds[0].(map[string]interface{}); ok {
			if h, ok := opts["hostname"].(string); ok {
				return h
			}
		}
	}

	if edge, ok := cfg["edge_options"].([]interface{}); ok && len(edge) > 0 {
		if opts, ok := edge[0].(map[string]interface{}); ok {
			if h, ok := opts["hostname"].(string); ok {
				return h
			}
		}
	}

	return ""
}

func deflateBluecatAnycast(anycast alkira.BluecatAnycast) []map[string]interface{} {
	// Return nil if anycast is empty to avoid spurious diffs
	// when the user hasn't configured anycast but API returns empty struct
	if len(anycast.Ips) == 0 && len(anycast.BackupCxps) == 0 {
		return nil
	}

	m := make(map[string]interface{})
	m["ips"] = anycast.Ips
	m["backup_cxps"] = anycast.BackupCxps

	return []map[string]interface{}{m}
}

func setAllBluecatResourceFields(d *schema.ResourceData, in *alkira.ServiceBluecat) {
	d.Set("bdds_anycast", deflateBluecatAnycast(in.BddsAnycast))
	d.Set("edge_anycast", deflateBluecatAnycast(in.EdgeAnycast))
	d.Set("billing_tag_ids", in.BillingTags)
	d.Set("cxp", in.Cxp)
	d.Set("description", in.Description)
	d.Set("global_cidr_list_id", in.GlobalCidrListId)
	d.Set("instance", deflateBluecatInstances(in.Instances, d))
	d.Set("license_type", in.LicenseType)
	d.Set("name", in.Name)
	d.Set("service_group_name", in.ServiceGroupName)
	d.Set("service_group_id", in.ServiceGroupId)
	d.Set("service_group_implicit_group_id", in.ServiceGroupImplicitGroupId)
}
