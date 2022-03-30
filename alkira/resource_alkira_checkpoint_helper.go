package alkira

import (
	"errors"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CheckpointGetSegById func(string) (alkira.Segment, error)

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

func expandCheckpointSegmentOptions(in *schema.Set, fn CheckpointGetSegById) (alkira.CheckpointSegmentNameToZone, error) {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid Checkpoint segment options input")
		return nil, nil
	}

	if in.Len() < 1 {
		return nil, nil
	}

	segmentName, zoneName, groups, err := convertCheckpointSegmentOptions(in, fn)
	if err != nil {
		return nil, err
	}

	z := make(alkira.CheckpointZoneToGroups)
	z[zoneName] = groups

	c := make(alkira.CheckpointSegmentNameToZone)
	c[segmentName] = z

	return c, nil
}

func convertCheckpointSegmentOptions(in *schema.Set, fn CheckpointGetSegById) (string, string, []string, error) {
	var segmentName string
	var zoneName string
	var groups []string

	if in == nil {
		return "", "", nil, errors.New("Checkpoint segment options cannot be nil")
	}

	for _, options := range in.List() {
		optionsCfg := options.(map[string]interface{})

		if v, ok := optionsCfg["segment_id"].(int); ok {
			if fn == nil {
				return "", "", nil, errors.New("Checkpoint's get segment by id (CheckpointGetSegById) cannot be nil")
			}

			sg, err := fn(strconv.Itoa(v))
			if err != nil {
				return "", "", nil, err
			}
			segmentName = sg.Name
		}

		if v, ok := optionsCfg["zone_name"].(string); ok {
			zoneName = v
		}

		if v, ok := optionsCfg["groups"].([]interface{}); ok {
			groups = convertTypeListToStringList(v)
		}
	}

	return segmentName, zoneName, groups, nil
}

func expandCheckpointManagementServer(in *schema.Set, fn CheckpointGetSegById) (*alkira.CheckpointManagementServer, error) {
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
			sg, err := fn(strconv.Itoa(v))
			if err != nil {
				return nil, err
			}

			mg.SegmentId = v
			mg.Segment = sg.Name
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

func deflateCheckpointSegmentOptions(c alkira.CheckpointSegmentNameToZone, fn CheckpointGetSegById) ([]map[string]interface{}, error) {

	var options []map[string]interface{}

	for segmentName, checkpointZoneToGroups := range c {
		seg, err := fn(segmentName)
		if err != nil {
			return nil, err
		}

		for zone, groups := range checkpointZoneToGroups {
			i := map[string]interface{}{
				"segment_id": seg.Id,
				"zone_name":  zone,
				"groups":     groups,
			}
			options = append(options, i)
		}
	}

	return options, nil
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
