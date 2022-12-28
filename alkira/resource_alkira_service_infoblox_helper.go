package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deflateInfobloxInstances(c []alkira.InfobloxInstance) []map[string]interface{} {
	var m []map[string]interface{}
	for _, v := range c {
		j := map[string]interface{}{
			"anycast_enabled": v.AnyCastEnabled,
			"hostname":        v.HostName,
			"model":           v.Model,
			"type":            v.Type,
			"version":         v.Version,
			"id":              v.Id,
			"credential_id":   v.CredentialId,
		}
		m = append(m, j)
	}

	return m
}

func deflateInfobloxGridMaster(im alkira.InfobloxGridMaster) []map[string]interface{} {
	m := make(map[string]interface{})
	m["external"] = im.External
	m["ip"] = im.Ip
	m["name"] = im.Name
	m["credential_id"] = im.GridMasterCredentialId

	return []map[string]interface{}{m}
}

func expandInfobloxAnycast(in *schema.Set) (*alkira.InfobloxAnycast, error) {
	if in == nil || in.Len() > 1 || in.Len() < 1 {
		return nil, fmt.Errorf("[DEBUG] Exactly one object allowed in anycast options.")
	}

	ia := &alkira.InfobloxAnycast{}

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})
		if v, ok := cfg["enabled"].(bool); ok {
			ia.Enabled = v
		}
		if v, ok := cfg["ips"].([]interface{}); ok {
			ia.Ips = convertTypeListToStringList(v)
		}
		if v, ok := cfg["backup_cxps"].([]interface{}); ok {
			ia.BackupCxps = convertTypeListToStringList(v)
		}
	}
	return ia, nil

}

func deflateInfobloxAnycast(ia alkira.InfobloxAnycast) []map[string]interface{} {
	m := make(map[string]interface{})
	m["enabled"] = ia.Enabled
	m["ips"] = ia.Ips
	m["backup_cxps"] = ia.BackupCxps

	return []map[string]interface{}{m}
}

func setAllInfobloxResourceFields(d *schema.ResourceData, in *alkira.Infoblox) {
	d.Set("anycast", deflateInfobloxAnycast(in.AnyCast))
	d.Set("billing_tag_ids", in.BillingTags)
	d.Set("cxp", in.Cxp)
	d.Set("global_cidr_list_id", in.GlobalCidrListId)
	d.Set("grid_master", deflateInfobloxGridMaster(in.GridMaster))
	d.Set("instance", deflateInfobloxInstances(in.Instances))
	d.Set("segment_ids", in.Segments)
	d.Set("service_group_name", in.ServiceGroupName)
	d.Set("license_type", in.LicenseType)
}

func expandGridMaster(in []interface{}, m interface{}) (*alkira.InfobloxGridMaster, error) {
	client := m.(*alkira.AlkiraClient)

	if len(in) != 1 {
		return nil, fmt.Errorf("[DEBUG] Exactly one object allowed in grid_master options.")
	}

	ig := &alkira.InfobloxGridMaster{}
	var username string
	var password string
	var sharedSecret string
	var exists bool

	for _, option := range in {
		cfg := option.(map[string]interface{})
		if v, ok := cfg["existing"].(bool); ok {
			exists = v
		}
		if v, ok := cfg["name"].(string); ok {
			ig.Name = v
		}
		if v, ok := cfg["password"].(string); ok {
			password = v
		}
		if v, ok := cfg["username"].(string); ok {
			username = v
		}
		if v, ok := cfg["grid_master_credential_id"].(string); ok {
			if v == "" {
				gridMasterCredentialId, err := client.CreateCredential(
					ig.Name+randomNameSuffix(),
					alkira.CredentialTypeInfobloxGridMaster,
					&alkira.CredentialInfobloxGridMaster{
						Username: username,
						Password: password,
					},
					0,
				)

				if err != nil {
					return nil, err
				}

				ig.GridMasterCredentialId = gridMasterCredentialId
			}

			if v != "" {
				ig.GridMasterCredentialId = v
			}
		}
		if v, ok := cfg["shared_secret"].(string); ok {
			sharedSecret = v
		}
		if v, ok := cfg["shared_secret_credential_id"].(string); ok {
			if v == "" {
				sharedSecretCredentialId, err := client.CreateCredential(
					ig.Name+randomNameSuffix(),
					alkira.CredentialTypeInfoblox,
					&alkira.CredentialInfoblox{SharedSecret: sharedSecret},
					0,
				)
				if err != nil {
					return nil, err
				}

				ig.SharedSecretCredentialId = sharedSecretCredentialId
			}

			if v != "" {
				ig.SharedSecretCredentialId = v
			}
		}
		if v, ok := cfg["ip"].(string); ok {
			if exists {
				ig.Ip = v
			}
		}
	}
	return ig, nil
}

func expandInfobloxInstances(in []interface{}, m interface{}) ([]alkira.InfobloxInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		return nil, fmt.Errorf("infoblox instances cannot be nil or empty")
	}
	var password string

	instances := make([]alkira.InfobloxInstance, len(in))
	for i, instance := range in {
		cfg := instance.(map[string]interface{})

		if v, ok := cfg["password"].(string); ok {
			password = v
		}
		if v, ok := cfg["hostname"].(string); ok {
			instances[i].HostName = v
		}
		if v, ok := cfg["credential_id"].(string); ok {
			if v == "" {
				credentialId, err := client.CreateCredential(
					instances[i].HostName+randomNameSuffix(),
					alkira.CredentialTypeInfobloxInstance,
					&alkira.CredentialInfobloxInstance{Password: password},
					0,
				)
				if err != nil {
					return nil, err
				}

				instances[i].CredentialId = credentialId
			}

			if v != "" {
				instances[i].CredentialId = v
			}
		}
		if v, ok := cfg["model"].(string); ok {
			instances[i].Model = v
		}
		if v, ok := cfg["type"].(string); ok {
			instances[i].Type = v
		}
		if v, ok := cfg["version"].(string); ok {
			instances[i].Version = v
		}
		if v, ok := cfg["anycast_enabled"].(bool); ok {
			instances[i].AnyCastEnabled = v
		}
	}

	return instances, nil
}
