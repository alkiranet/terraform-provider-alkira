package alkira

import (
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandCheckpointManagementServer(name string, in *schema.Set, m interface{}) (*alkira.CheckpointManagementServer, error) {

	client := m.(*alkira.AlkiraClient)

	if in == nil || in.Len() > 1 {
		log.Printf("[DEBUG] Invalid Checkpoint Firewall Management Server input.")
		return nil, errors.New("Invalid Checkpoint Firewall Management Server input.")
	}

	if in.Len() < 1 {
		return nil, nil
	}

	mg := &alkira.CheckpointManagementServer{}
	var manServerPass string

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})
		if v, ok := cfg["configuration_mode"].(string); ok {
			mg.ConfigurationMode = v
		}
		if v, ok := cfg["password"].(string); ok {
			manServerPass = v
		}
		if v, ok := cfg["credential_id"].(string); ok {
			if v == "" && mg.ConfigurationMode == "AUTOMATED" {
				manServerCredName := name + "-" + randomNameSuffix()
				c := &alkira.CredentialCheckPointFwManagementServer{Password: manServerPass}
				credentialId, err := client.CreateCredential(manServerCredName, alkira.CredentialTypeChkpFwManagement, c, 0)
				if err != nil {
					return nil, err
				}
				mg.CredentialId = credentialId
			}

			if v != "" {
				mg.CredentialId = v
			}
		}
		if v, ok := cfg["domain"].(string); ok {
			mg.Domain = v
		}
		if v, ok := cfg["global_cidr_list_id"].(int); ok {
			mg.GlobalCidrListId = v
		}
		if v, ok := cfg["ips"].([]interface{}); ok {
			mg.Ips = convertTypeListToStringList(v)
		}
		if v, ok := cfg["reachability"].(string); ok {
			mg.Reachability = v
		}
		if v, ok := cfg["segment_id"].(string); ok {
			if v != "" {
				segment, err := getSegmentNameById(v, m)

				if err != nil {
					return nil, err
				}

				mg.Segment = segment
			}
		}
		if v, ok := cfg["type"].(string); ok {
			mg.Type = v
		}
		if v, ok := cfg["username"].(string); ok {
			mg.UserName = v
		}
	}
	return mg, nil
}

func expandCheckpointInstances(in []interface{}, m interface{}) ([]alkira.CheckpointInstance, error) {

	if in == nil || len(in) == 0 {
		return nil, errors.New("Invalid Checkpoint Firewall instance input.")
	}

	client := m.(*alkira.AlkiraClient)

	instances := make([]alkira.CheckpointInstance, len(in))
	for i, instance := range in {
		r := alkira.CheckpointInstance{}
		instanceCfg := instance.(map[string]interface{})

		var sicKey string

		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["sic_key"].(string); ok {
			sicKey = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			if v == "" {
				credentialName := r.Name + "-" + randomNameSuffix()
				c := &alkira.CredentialCheckPointFwServiceInstance{SicKey: sicKey}

				log.Printf("[INFO] Creating Credential CheckpointInstance.")
				credentialId, err := client.CreateCredential(
					credentialName,
					alkira.CredentialTypeChkpFwInstance,
					c,
					0)

				if err != nil {
					return nil, err
				}

				r.CredentialId = credentialId
			} else {
				r.CredentialId = v
			}
		}
		if v, ok := instanceCfg["enable_traffic"].(bool); ok {
			r.TrafficEnabled = v
		}
		instances[i] = r
	}

	return instances, nil
}

/*
Checkpoint expects segment_options to not be empty.
If segment_options is not defined in the TF file, this function adds the default expected data.
If segment_options is included, populates it normally.
*/
func expandCheckpointSegmentOptions(segmentName string, in *schema.Set, m interface{}) (alkira.SegmentNameToZone, error) {

	if in == nil || in.Len() == 0 {

		segmentOptions := make(alkira.SegmentNameToZone)
		zonestoGroups := make(alkira.ZoneToGroups)

		z := alkira.OuterZoneToGroups{}
		j := []string{}

		zonestoGroups["DEFAULT"] = j
		z.ZonesToGroups = zonestoGroups

		segmentOptions[segmentName] = z

		return segmentOptions, nil
	}

	return expandSegmentOptions(in, m)

}

func deflateCheckpointManagementServer(mg alkira.CheckpointManagementServer) []map[string]interface{} {
	m := make(map[string]interface{})
	m["configuration_mode"] = mg.ConfigurationMode
	m["credential_id"] = mg.CredentialId
	m["domain"] = mg.Domain
	m["global_cidr_list_id"] = mg.GlobalCidrListId
	m["ips"] = convertStringArrToInterfaceArr(mg.Ips)
	m["reachability"] = mg.Reachability
	m["segment"] = mg.Segment
	m["type"] = mg.Type
	m["user_name"] = mg.UserName

	return []map[string]interface{}{m}
}

func setCheckpointInstances(d *schema.ResourceData, c []alkira.CheckpointInstance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, value := range d.Get("instance").([]interface{}) {
		cfg := value.(map[string]interface{})

		for _, ins := range c {
			if cfg["id"].(int) == ins.Id || cfg["name"].(string) == ins.Name {
				instance := map[string]interface{}{
					"credential_id":  ins.CredentialId,
					"name":           ins.Name,
					"id":             ins.Id,
					"sic_key":        cfg["sic_key"].(string),
					"enable_traffic": ins.TrafficEnabled,
				}
				instances = append(instances, instance)
				break
			}
		}
	}

	for _, instance := range c {
		new := true

		// Check if the instance already exists in the Terraform config
		for _, ins := range d.Get("instance").([]interface{}) {
			cfg := ins.(map[string]interface{})

			if cfg["id"].(int) == instance.Id || cfg["name"].(string) == instance.Name {
				new = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			instance := map[string]interface{}{
				"credential_id": instance.CredentialId,
				"name":          instance.Name,
				"id":            instance.Id,
			}

			instances = append(instances, instance)
			break
		}
	}

	return instances
}
