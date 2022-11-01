package alkira

import (
	"errors"
	"log"
	"strconv"

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
		if v, ok := cfg["management_server_password"].(string); ok {
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
		if v, ok := cfg["segment_id"].(int); ok {

			var sg alkira.Segment
			var err error

			// 0 is an invalid ID but also the default value of int
			if v != 0 {
				sg, err = client.GetSegmentById(strconv.Itoa(v))
				if err != nil {
					return nil, err
				}
				mg.SegmentId = v
				mg.Segment = sg.Name
			}

		}
		if v, ok := cfg["type"].(string); ok {
			mg.Type = v
		}
		if v, ok := cfg["user_name"].(string); ok {
			mg.UserName = v
		}
	}
	return mg, nil
}

func expandCheckpointInstances(name string, in *schema.Set, m interface{}) ([]alkira.CheckpointInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid Checkpoint Firewall instance input.")
		return nil, errors.New("Invalid Checkpoint Firewall instance input.")
	}

	var chkpfwInstanceKey string
	instances := make([]alkira.CheckpointInstance, in.Len())
	for i, instance := range in.List() {
		r := alkira.CheckpointInstance{}
		instanceCfg := instance.(map[string]interface{})
		r.Name = name + "-instance-" + strconv.Itoa(i+1)
		if v, ok := instanceCfg["sic_key"].(string); ok {
			chkpfwInstanceKey = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			if v == "" {
				credentialName := r.Name + "-" + randomNameSuffix()
				c := &alkira.CredentialCheckPointFwServiceInstance{SicKey: chkpfwInstanceKey}
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
			}
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

func convertCheckpointSegmentNameToSegmentId(names []string, m interface{}) (int, error) {
	client := m.(*alkira.AlkiraClient)

	if len(names) != 1 {
		log.Printf("[DEBUG] invalid number of segments in Checkpoint Firewall instance.")
		return 0, errors.New("Invalid number of segments in Checkpoint Firewall instance.")
	}

	seg, err := client.GetSegmentByName(names[0])
	if err != nil {
		log.Printf("[DEBUG] failed to get segment. %s does not exist: ", names[0])
		return 0, err
	}
	return seg.Id, nil
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
	m["segment_id"] = mg.SegmentId
	m["type"] = mg.Type
	m["user_name"] = mg.UserName

	return []map[string]interface{}{m}
}

func deflateCheckpointInstances(c []alkira.CheckpointInstance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, instance := range c {
		i := map[string]interface{}{
			"name":          instance.Name,
			"credential_id": instance.CredentialId,
		}
		instances = append(instances, i)
	}

	return instances
}
