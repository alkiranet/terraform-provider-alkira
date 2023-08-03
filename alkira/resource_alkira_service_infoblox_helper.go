package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// func expandInfobloxInstances(in *schema.Set, m interface{}) ([]alkira.InfobloxInstance, error) {
func expandInfobloxInstances(in []interface{}, m interface{}) ([]alkira.InfobloxInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		return nil, fmt.Errorf("infoblox instances cannot be nil or empty")
	}

	instances := make([]alkira.InfobloxInstance, len(in))
	for i, instance := range in {
		var r alkira.InfobloxInstance
		var nameWithSuffix string
		var password string

		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["anycast_enabled"].(bool); ok {
			r.AnyCastEnabled = v
		}
		if v, ok := instanceCfg["id"].(int); ok {
			if v != 0 {
				r.Id = json.Number(strconv.Itoa(v))
			}
		}
		if v, ok := instanceCfg["hostname"].(string); ok {
			//Note: Name is required but not used in the API. So rather than make our user input an
			//extra field that we just ignore anyway r.Name is set to hostname and the credential
			//name is based off the hostname as well.
			r.Name = v
			r.HostName = v
			nameWithSuffix = v + randomNameSuffix()
		}
		if v, ok := instanceCfg["model"].(string); ok {
			r.Model = v
		}
		if v, ok := instanceCfg["password"].(string); ok {
			password = v
		}
		if v, ok := instanceCfg["type"].(string); ok {
			r.Type = v
		}
		if v, ok := instanceCfg["version"].(string); ok {
			r.Version = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			if v == "" {
				credentialInstance := alkira.CredentialInfobloxInstance{
					Password: password,
				}

				credentialId, err := client.CreateCredential(
					nameWithSuffix,
					alkira.CredentialTypeInfobloxInstance,
					credentialInstance,
					0,
				)

				if err != nil {
					return nil, err
				}

				r.CredentialId = credentialId
			}

			if v != "" {
				r.CredentialId = v
			}
		}

		instances[i] = r
	}

	return instances, nil
}

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

func expandInfobloxGridMaster(in []interface{}, sharedSecretCredentialId string, m interface{}) (*alkira.InfobloxGridMaster, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) > 1 || len(in) < 1 {
		return nil, fmt.Errorf("[DEBUG] Exactly one object allowed in grid master options.")
	}

	im := &alkira.InfobloxGridMaster{}

	var username string
	var password string
	for _, option := range in {

		cfg := option.(map[string]interface{})
		if v, ok := cfg["external"].(bool); ok {
			im.External = v
		}
		if v, ok := cfg["ip"].(string); ok {
			im.Ip = v
		}
		if v, ok := cfg["username"].(string); ok {
			username = v
		}
		if v, ok := cfg["password"].(string); ok {
			password = v
		}
		if v, ok := cfg["name"].(string); ok {
			im.Name = v
		}
		if v, ok := cfg["credential_id"].(string); ok {
			if v == "" {
				gridMasterCredentialId, err := client.CreateCredential(
					im.Name+randomNameSuffix(),
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

				im.GridMasterCredentialId = gridMasterCredentialId
			}

			if v != "" {
				im.GridMasterCredentialId = v
			}
		}
	}

	im.SharedSecretCredentialId = sharedSecretCredentialId

	return im, nil
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

func setAllInfobloxResourceFields(d *schema.ResourceData, in *alkira.ServiceInfoblox) {
	d.Set("anycast", deflateInfobloxAnycast(in.AnyCast))
	d.Set("billing_tag_ids", in.BillingTags)
	d.Set("cxp", in.Cxp)
	d.Set("description", in.Description)
	d.Set("global_cidr_list_id", in.GlobalCidrListId)
	d.Set("grid_master", deflateInfobloxGridMaster(in.GridMaster))
	d.Set("instance", deflateInfobloxInstances(in.Instances))
	d.Set("license_type", in.LicenseType)
	d.Set("segment_ids", in.Segments)
	d.Set("service_group_name", in.ServiceGroupName)
	d.Set("size", in.Size)
}
