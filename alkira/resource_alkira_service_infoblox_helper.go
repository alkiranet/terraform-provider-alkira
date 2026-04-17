package alkira

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateInfobloxGridMasterCredential always updates the grid master
// credential because username and password are WriteOnly and HasChanges
// cannot detect changes (state is always null for WriteOnly fields).
func updateInfobloxGridMasterCredential(d *schema.ResourceData, c *alkira.AlkiraClient) error {
	log.Printf("[INFO] Updating Infoblox Grid Master Credential")

	// Read credential_id from grid_master block in state
	gmList := d.Get("grid_master").([]interface{})
	if len(gmList) == 0 {
		return nil
	}
	cfg := gmList[0].(map[string]interface{})
	credentialId, ok := cfg["credential_id"].(string)
	if !ok || credentialId == "" {
		log.Printf("[INFO] grid_master credential_id is empty, skipping grid master credential update")
		return nil
	}

	username := getInfobloxNestedWriteOnlyValue(d, "grid_master", 0, "username")
	password := getInfobloxNestedWriteOnlyValue(d, "grid_master", 0, "password")

	cred, err := c.GetCredentialById(credentialId)
	if err != nil {
		return fmt.Errorf("failed to get grid master credential name for ID %s: %w", credentialId, err)
	}

	credential := &alkira.CredentialInfobloxGridMaster{
		Username: username,
		Password: password,
	}

	return c.UpdateCredential(credentialId, cred.Name, alkira.CredentialTypeInfobloxGridMaster, credential, 0)
}

// updateInfobloxInstanceCredentials always updates all instance
// credentials because password is WriteOnly and HasChanges cannot
// detect changes (state is always null for WriteOnly fields).
func updateInfobloxInstanceCredentials(d *schema.ResourceData, c *alkira.AlkiraClient) error {
	log.Printf("[INFO] Updating Infoblox Instance Credentials")

	instances := d.Get("instance").([]interface{})
	for i, inst := range instances {
		cfg := inst.(map[string]interface{})
		credentialId, ok := cfg["credential_id"].(string)
		if !ok || credentialId == "" {
			continue
		}

		password := getInfobloxNestedWriteOnlyValue(d, "instance", i, "password")

		cred, err := c.GetCredentialById(credentialId)
		if err != nil {
			return fmt.Errorf("failed to get instance credential name for ID %s: %w", credentialId, err)
		}

		credential := alkira.CredentialInfobloxInstance{
			Password: password,
		}

		err = c.UpdateCredential(credentialId, cred.Name, alkira.CredentialTypeInfobloxInstance, credential, 0)
		if err != nil {
			return fmt.Errorf("failed to update instance credential %s: %w", credentialId, err)
		}
	}

	return nil
}

// updateInfobloxCredentials updates all infoblox credentials on every apply.
func updateInfobloxCredentials(d *schema.ResourceData, c *alkira.AlkiraClient) error {
	if err := updateInfobloxGridMasterCredential(d, c); err != nil {
		return err
	}
	if err := updateInfobloxInstanceCredentials(d, c); err != nil {
		return err
	}
	return nil
}

// getInfobloxWriteOnlyValue reads a WriteOnly field from raw config.
// WriteOnly fields are never stored in state, so d.GetOk() would always
// return the zero value — only GetRawConfigAt can read the actual value.
func getInfobloxWriteOnlyValue(d *schema.ResourceData, field string) string {
	attrPath := cty.Path{cty.GetAttrStep{Name: field}}
	val, diags := d.GetRawConfigAt(attrPath)

	if !diags.HasError() && !val.IsNull() && val.IsKnown() && val.Type() == cty.String {
		if strVal := val.AsString(); strVal != "" {
			return strVal
		}
	}

	return ""
}

// getInfobloxNestedWriteOnlyValue reads a WriteOnly field from a nested TypeList block
// via GetRawConfigAt with an indexed cty path.
func getInfobloxNestedWriteOnlyValue(d *schema.ResourceData, block string, index int, field string) string {
	attrPath := cty.Path{
		cty.GetAttrStep{Name: block},
		cty.IndexStep{Key: cty.NumberIntVal(int64(index))},
		cty.GetAttrStep{Name: field},
	}
	val, diags := d.GetRawConfigAt(attrPath)

	if !diags.HasError() && !val.IsNull() && val.IsKnown() && val.Type() == cty.String {
		if strVal := val.AsString(); strVal != "" {
			return strVal
		}
	}

	return ""
}

// func expandInfobloxInstances(in *schema.Set, m interface{}) ([]alkira.InfobloxInstance, error) {
func expandInfobloxInstances(d *schema.ResourceData, in []interface{}, m interface{}) ([]alkira.InfobloxInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		return nil, fmt.Errorf("ERROR: Infoblox instances cannot be nil or empty")
	}

	instances := make([]alkira.InfobloxInstance, len(in))
	for i, instance := range in {
		var r alkira.InfobloxInstance
		var nameWithSuffix string

		// Read password from raw config (WriteOnly field)
		password := getInfobloxNestedWriteOnlyValue(d, "instance", i, "password")

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

func expandInfobloxGridMaster(d *schema.ResourceData, in []interface{}, sharedSecretCredentialId string, m interface{}) (*alkira.InfobloxGridMaster, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) > 1 || len(in) < 1 {
		return nil, fmt.Errorf("ERROR: Exactly one object allowed in grid master options")
	}

	im := &alkira.InfobloxGridMaster{}

	// Read WriteOnly fields from raw config
	username := getInfobloxNestedWriteOnlyValue(d, "grid_master", 0, "username")
	password := getInfobloxNestedWriteOnlyValue(d, "grid_master", 0, "password")

	for _, option := range in {

		cfg := option.(map[string]interface{})
		if v, ok := cfg["external"].(bool); ok {
			im.External = v
		}
		if v, ok := cfg["ip"].(string); ok {
			im.Ip = v
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
	d.Set("allow_list_id", in.AllowListId)
	d.Set("service_group_id", in.ServiceGroupId)
	d.Set("service_group_implicit_group_id", in.ServiceGroupImplicitGroupId)
	d.Set("shared_secret_credential_id", in.GridMaster.SharedSecretCredentialId)
}
