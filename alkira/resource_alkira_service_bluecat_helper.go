package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"

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
			if v != 0 {
				r.Id = json.Number(strconv.Itoa(v))
			}
		}
		if v, ok := instanceCfg["type"].(string); ok {
			r.Type = v
		}

		// Handle BDDS options
		if bddsOptions, ok := instanceCfg["bdds_options"].(*schema.Set); ok && bddsOptions.Len() > 0 {
			bddsOpt, err := expandBDDSOptions(bddsOptions.List(), m)
			if err != nil {
				return nil, err
			}
			r.BddsOptions = bddsOpt
		}

		// Handle Edge options
		if edgeOptions, ok := instanceCfg["edge_options"].(*schema.Set); ok && edgeOptions.Len() > 0 {
			edgeOpt, err := expandEdgeOptions(edgeOptions.List(), m)
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
				options.HostName+randomNameSuffix(),
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
				options.HostName+randomNameSuffix(),
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

func deflateBluecatInstances(c []alkira.BluecatInstance) []map[string]interface{} {
	var m []map[string]interface{}
	for _, v := range c {
		j := map[string]interface{}{
			"id":   v.Id,
			"name": v.Name,
			"type": v.Type,
		}

		if v.BddsOptions != nil {
			bddsMap := map[string]interface{}{
				"hostname":              v.BddsOptions.HostName,
				"model":                 v.BddsOptions.Model,
				"version":               v.BddsOptions.Version,
				"license_credential_id": v.BddsOptions.LicenseCredentialId,
			}
			j["bdds_options"] = []interface{}{bddsMap}
		}

		if v.EdgeOptions != nil {
			edgeMap := map[string]interface{}{
				"hostname":      v.EdgeOptions.HostName,
				"version":       v.EdgeOptions.Version,
				"credential_id": v.EdgeOptions.CredentialId,
			}
			j["edge_options"] = []interface{}{edgeMap}
		}

		m = append(m, j)
	}

	return m
}

func deflateBluecatAnycast(anycast alkira.BluecatAnycast) []map[string]interface{} {
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
	d.Set("instance", deflateBluecatInstances(in.Instances))
	d.Set("license_type", in.LicenseType)
	d.Set("name", in.Name)
	d.Set("service_group_name", in.ServiceGroupName)
	d.Set("service_group_id", in.ServiceGroupId)
	d.Set("service_group_implicit_group_id", in.ServiceGroupImplicitGroupId)
}
