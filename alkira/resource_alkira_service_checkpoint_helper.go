package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CheckpointGetSegById func(string) (alkira.Segment, error)

func checkpointRespDetailsToCheckpointInstance(details []alkira.CredentialResponseDetail) []alkira.CheckpointInstance {
	var instances []alkira.CheckpointInstance

	for _, v := range details {
		instances = append(instances, alkira.CheckpointInstance{
			CredentialId: v.Id,
			Name:         v.Name,
		})
	}

	return instances
}

func expandCheckpointInstances(in *schema.Set) []alkira.CheckpointInstance {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid Checkpoint instance input")
		return nil
	}

	instances := make([]alkira.CheckpointInstance, in.Len())
	for i, instance := range in.List() {
		r := alkira.CheckpointInstance{}
		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			r.CredentialId = v
		}
		instances[i] = r
	}

	return instances
}

func expandCheckpointManagementServer(in *schema.Set, m interface{}) (*alkira.CheckpointManagementServer, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || in.Len() > 1 {
		log.Printf("[DEBUG] Only one object allowed in managment server options")
		return nil, nil
	}

	if in.Len() < 1 {
		return nil, nil
	}

	mg := &alkira.CheckpointManagementServer{}

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})
		if v, ok := cfg["configuration_mode"].(string); ok {
			mg.ConfigurationMode = v
		}
		if v, ok := cfg["credential_id"].(string); ok {
			mg.CredentialId = v
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
