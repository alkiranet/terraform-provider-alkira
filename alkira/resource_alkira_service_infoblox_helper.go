package alkira

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandInfobloxInstances(in []interface{}, oldInstances []interface{}, m interface{}) ([]alkira.InfobloxInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		return nil, fmt.Errorf("ERROR: Infoblox instances cannot be nil or empty")
	}

	// Build hostname -> id/credential_id maps from old state. When an instance is
	// inserted or removed in the middle of a TypeList, Terraform's positional diff
	// shifts id/credential_id values to the wrong element — the request then tells
	// the API to delete the wrong instance (observed live: removing the middle
	// instance destroyed the tail one). Matching by hostname keeps each existing
	// instance bound to its own ids regardless of list position (same fix as Bluecat).
	oldIdByHostname := make(map[string]int)
	oldCredentialIdByHostname := make(map[string]string)
	oldCredentialIdOwner := make(map[string]string)
	for _, old := range oldInstances {
		cfg, ok := old.(map[string]interface{})
		if !ok {
			log.Printf("[WARN] Infoblox: skipping malformed instance entry in old state: %v", old)
			continue
		}
		hostname, _ := cfg["hostname"].(string)
		if hostname == "" {
			continue
		}
		if id, ok := cfg["id"].(int); ok && id != 0 {
			oldIdByHostname[hostname] = id
		}
		if cid, ok := cfg["credential_id"].(string); ok && cid != "" {
			oldCredentialIdByHostname[hostname] = cid
			oldCredentialIdOwner[cid] = hostname
		}
	}

	instances := make([]alkira.InfobloxInstance, len(in))
	for i, instance := range in {
		var r alkira.InfobloxInstance
		var nameWithSuffix string
		var password string
		var joinToken string

		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["anycast_enabled"].(bool); ok {
			r.AnyCastEnabled = v
		}
		if v, ok := instanceCfg["hostname"].(string); ok {
			//Note: Name is required but not used in the API. So rather than make our user input an
			//extra field that we just ignore anyway r.Name is set to hostname and the credential
			//name is based off the hostname as well.
			r.Name = v
			r.HostName = v
			nameWithSuffix = v + randomNameSuffix()
		}
		// id by hostname from old state — never trust the positionally-diffed value.
		// Not found => new instance, id stays unset.
		if id, found := oldIdByHostname[r.HostName]; found {
			r.Id = json.Number(strconv.Itoa(id))
		}
		if v, ok := instanceCfg["model"].(string); ok {
			r.Model = v
		}
		if v, ok := instanceCfg["password"].(string); ok {
			password = v
		}
		if v, ok := instanceCfg["platform"].(string); ok {
			r.Platform = v
		}
		if v, ok := instanceCfg["join_token"].(string); ok {
			joinToken = v
		}
		if v, ok := instanceCfg["type"].(string); ok {
			r.Type = v
		}
		if v, ok := instanceCfg["version"].(string); ok {
			r.Version = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			// credential_id by hostname from old state — same positional-shift hazard as id.
			if cid, found := oldCredentialIdByHostname[r.HostName]; found {
				v = cid
			} else if v != "" && oldCredentialIdOwner[v] != "" {
				// A new instance whose list position previously belonged to another instance
				// inherits that instance's credential_id from the positional diff. Never reuse
				// it — mint a fresh credential for the new instance.
				log.Printf("[WARN] Infoblox: new instance %q inherited credential_id of instance %q from its list position; creating a fresh credential",
					r.HostName, oldCredentialIdOwner[v])
				v = ""
			}
			if v == "" {
				var credentialId string
				var err error

				// NIOS-X uses a dedicated INFOBLOX_JOIN_TOKEN credential (join token in the
				// #cloud-config user-data); NIOS uses an INFOBLOX_INSTANCE credential (admin password).
				if r.Platform == "NIOS_X" {
					if joinToken == "" {
						return nil, fmt.Errorf("ERROR: 'join_token' is required for NIOS_X instance %q", r.HostName)
					}
					credentialId, err = client.CreateCredential(
						nameWithSuffix,
						alkira.CredentialTypeInfobloxJoinToken,
						alkira.CredentialInfobloxJoinToken{
							JoinToken: joinToken,
						},
						0,
					)
				} else {
					credentialId, err = client.CreateCredential(
						nameWithSuffix,
						alkira.CredentialTypeInfobloxInstance,
						alkira.CredentialInfobloxInstance{
							Password: password,
						},
						0,
					)
				}

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

// validateInfobloxInstanceHostnames returns an error if any two instances share a
// hostname. Hostnames are the unique keys used to match instances across list
// reorders (id/credential_id remap in expandInfobloxInstances); duplicates make
// that lookup unreliable. Same guard as Bluecat's validateBluecatInstanceHostnames.
func validateInfobloxInstanceHostnames(instances []interface{}) error {
	seen := make(map[string]int, len(instances))
	for i, inst := range instances {
		cfg, ok := inst.(map[string]interface{})
		if !ok {
			continue
		}
		hostname, _ := cfg["hostname"].(string)
		if hostname == "" {
			continue
		}
		if prev, exists := seen[hostname]; exists {
			return fmt.Errorf(
				"instance[%d] and instance[%d] both use hostname %q; hostnames must be unique across all instances",
				prev, i, hostname,
			)
		}
		seen[hostname] = i
	}
	return nil
}

func deflateInfobloxInstances(c []alkira.InfobloxInstance) []map[string]interface{} {
	var m []map[string]interface{}
	for _, v := range c {
		id, err := v.Id.Int64()
		if err != nil {
			log.Printf("[WARN] failed to convert infoblox instance id %q to int: %s", v.Id, err)
		}
		j := map[string]interface{}{
			"anycast_enabled": v.AnyCastEnabled,
			"hostname":        v.HostName,
			"model":           v.Model,
			"platform":   v.Platform,
			"type":            v.Type,
			"version":         v.Version,
			"id":              int(id),
			"credential_id":   v.CredentialId,
		}
		m = append(m, j)
	}

	return m
}

func expandInfobloxGridMaster(in []interface{}, sharedSecretCredentialId string, m interface{}) (*alkira.InfobloxGridMaster, error) {
	client := m.(*alkira.AlkiraClient)

	if len(in) > 1 {
		return nil, fmt.Errorf("ERROR: at most one grid_master object allowed")
	}

	im := &alkira.InfobloxGridMaster{}

	// Only a NIOS-X-only service may omit grid_master, and generateInfobloxRequest handles
	// that case before calling here — reaching this with no block means NIOS instances exist.
	if len(in) == 0 {
		return nil, fmt.Errorf("ERROR: grid_master is required unless all instances are platform NIOS_X")
	}

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
		return nil, fmt.Errorf("ERROR: Exactly one object allowed in anycast options")
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
	if in == nil {
		return
	}
	d.Set("name", in.Name)
	d.Set("anycast", deflateInfobloxAnycast(in.AnyCast))
	d.Set("billing_tag_ids", in.BillingTags)
	d.Set("cxp", in.Cxp)
	d.Set("description", in.Description)
	d.Set("global_cidr_list_id", in.GlobalCidrListId)
	d.Set("grid_master", preserveInfobloxGridMasterSecrets(d, deflateInfobloxGridMaster(in.GridMaster), in.GridMaster))
	d.Set("instance", preserveInfobloxInstanceSecrets(d, deflateInfobloxInstances(in.Instances)))
	d.Set("license_type", in.LicenseType)
	d.Set("service_group_name", in.ServiceGroupName)
	d.Set("size", in.Size)
	d.Set("allow_list_id", in.AllowListId)
	d.Set("service_group_id", in.ServiceGroupId)
	d.Set("service_group_implicit_group_id", in.ServiceGroupImplicitGroupId)
}

// The API never returns the write-only secrets (instance join_token/password, grid_master
// username/password). Rewriting state without them stores "" while the config holds the real
// value, so every plan shows an in-place change and every apply re-triggers a provisioning
// run. Preserve them from the prior state (which mirrors the config after apply).
func preserveInfobloxInstanceSecrets(d *schema.ResourceData, instances []map[string]interface{}) []map[string]interface{} {
	prior := map[string]map[string]interface{}{}
	if v, ok := d.Get("instance").([]interface{}); ok {
		for _, p := range v {
			if pm, ok := p.(map[string]interface{}); ok {
				if hn, ok := pm["hostname"].(string); ok && hn != "" {
					prior[hn] = pm
				}
			}
		}
	}
	for _, j := range instances {
		if hn, ok := j["hostname"].(string); ok {
			if pm, ok := prior[hn]; ok {
				j["join_token"] = pm["join_token"]
				j["password"] = pm["password"]
			}
		}
	}
	return instances
}

func preserveInfobloxGridMasterSecrets(d *schema.ResourceData, gridMaster []map[string]interface{}, api alkira.InfobloxGridMaster) []map[string]interface{} {
	prior, _ := d.Get("grid_master").([]interface{})
	// NIOS-X-only service: the API returns an empty gridMaster and the config carries no
	// grid_master block — keep state empty, else refresh adds a phantom block that every
	// subsequent plan shows as a removal.
	if len(prior) == 0 && api.Name == "" && api.Ip == "" && api.GridMasterCredentialId == "" {
		return []map[string]interface{}{}
	}
	if len(prior) == 1 && len(gridMaster) == 1 {
		if pm, ok := prior[0].(map[string]interface{}); ok {
			gridMaster[0]["username"] = pm["username"]
			gridMaster[0]["password"] = pm["password"]
		}
	}
	return gridMaster
}
