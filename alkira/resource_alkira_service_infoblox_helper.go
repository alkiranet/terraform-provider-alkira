package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandInfobloxInstances(in *schema.Set, createCredential createCredential) ([]alkira.InfobloxInstance, error) {
	if in == nil || in.Len() == 0 {
		return nil, fmt.Errorf("[DEBUG] invalid Infoblox instance input")
	}

	instances := make([]alkira.InfobloxInstance, in.Len())
	for i, instance := range in.List() {
		var r alkira.InfobloxInstance
		var name string
		var password string

		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["any_cast_enabled"].(bool); ok {
			r.AnyCastEnabled = v
		}
		if v, ok := instanceCfg["host_name"].(string); ok {
			r.HostName = v
		}
		if v, ok := instanceCfg["model"].(string); ok {
			r.Model = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			name = v + randomNameSuffix()
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

		credentialInstance := alkira.CredentialInfobloxInstance{password}
		credentialId, err := createCredential(name, alkira.CredentialTypeInfobloxInstance, credentialInstance)
		if err != nil {
			return nil, err
		}

		r.CredentialId = credentialId

		instances[i] = r
	}

	return instances, nil
}

func deflateInfobloxInstances(c []alkira.InfobloxInstance) []map[string]interface{} {
	var m []map[string]interface{}
	for _, v := range c {
		j := map[string]interface{}{
			"any_cast_enabled": v.AnyCastEnabled,
			"host_name":        v.HostName,
			"model":            v.Model,
			"name":             v.Name,
			"type":             v.Type,
			"version":          v.Version,
		}
		m = append(m, j)
	}

	return m
	//var instances []map[string]interface{}

	//for _, instance := range c {
	//	i := map[string]interface{}{
	//		"any_cast_enabled": instance.AnyCastEnabled,
	//		//"credential_id":    instance.CredentialId,
	//		//"host_name":        instance.HostName,
	//		//"model":            instance.Model,
	//		//"type":             instance.Type,
	//		//"version":          instance.Version,
	//	}
	//	instances = append(instances, i)
	//}

	//fmt.Println("len(instances): ", len(instances))
	//return instances
}

func expandInfobloxGridMaster(in *schema.Set, sharedSecretCredentialId string, createCredential createCredential) (*alkira.InfobloxGridMaster, error) {

	if in == nil || in.Len() > 1 || in.Len() < 1 {
		return nil, fmt.Errorf("[DEBUG] Exactly one object allowed in grid master options.")
	}

	im := &alkira.InfobloxGridMaster{}

	var username string
	var password string
	for _, option := range in.List() {

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
	}

	gridMasterCredentialId, err := createCredential(
		im.Name+randomNameSuffix(),
		alkira.CredentialTypeInfobloxGridMaster,
		&alkira.CredentialInfobloxGridMaster{username, password},
	)
	if err != nil {
		return nil, err
	}

	im.GridMasterCredentialId = gridMasterCredentialId
	im.SharedSecretCredentialId = sharedSecretCredentialId

	return im, nil
}

func deflateInfobloxGridMaster(im alkira.InfobloxGridMaster) []map[string]interface{} {
	m := make(map[string]interface{})
	m["external"] = im.External
	m["ip"] = im.Ip
	m["name"] = im.Name

	return []map[string]interface{}{m}
}

func expandInfobloxAnycast(in *schema.Set) (*alkira.InfobloxAnycast, error) {
	if in == nil || in.Len() > 1 || in.Len() < 1 {
		return nil, fmt.Errorf("[DEBUG] Exactly one object allowed in anycast options.")
	}

	ia := &alkira.InfobloxAnycast{}

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})
		//convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
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
	//fmt.Println()
	//fmt.Println("in: ", in)
	//fmt.Println()
	//fmt.Println()
	//fmt.Println("Type of deflateInfobloxInstances", reflect.TypeOf(deflateInfobloxInstances(in.Instances)))
	//fmt.Println()
	//fmt.Println("deflateInfobloxInstances(in.Instances): ", deflateInfobloxInstances(in.Instances))
	//fmt.Println()
	//fmt.Println("Type of deflateInfobloxGridMaster", reflect.TypeOf(deflateInfobloxGridMaster(in.GridMaster)))
	//fmt.Println()
	//fmt.Println("deflateInfobloxGridMaster: ", deflateInfobloxGridMaster(in.GridMaster))
	//fmt.Println()

	d.Set("anycast", deflateInfobloxAnycast(in.AnyCast))
	d.Set("billing_tag_ids", in.BillingTags)
	d.Set("cxp", in.Cxp)
	d.Set("description", in.Description)
	d.Set("global_cidr_list_id", in.GlobalCidrListId)
	d.Set("grid_master", deflateInfobloxGridMaster(in.GridMaster))
	d.Set("instances", deflateInfobloxInstances(in.Instances))
	d.Set("license_type", in.LicenseType)
	d.Set("segment_names", in.Segments)
	d.Set("service_group_name", in.ServiceGroupName)
	d.Set("size", in.Size)

}
