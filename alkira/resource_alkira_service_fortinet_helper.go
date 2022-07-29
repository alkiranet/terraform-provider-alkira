package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandFortinetInstances(in *schema.Set) []alkira.FortinetInstance {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid Fortinet instance input")
		return nil
	}

	instances := make([]alkira.FortinetInstance, in.Len())
	for i, instance := range in.List() {
		r := alkira.FortinetInstance{}
		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
			r.HostName = v
		}
		if v, ok := instanceCfg["serial_number"].(string); ok {
			r.SerialNumber = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			r.CredentialId = v
		}
		instances[i] = r
	}

	return instances
}

func expandFortinetSegmentOptions(in *schema.Set, m interface{}) (map[string]alkira.FortinetSegmentName, error) {
	if in == nil || in.Len() == 0 {
		//At the time of this writing segment options is optional we don't care if they don't submit anything.
		return nil, nil
	}

	client := m.(*alkira.AlkiraClient)

	segmentOptions := make(map[string]alkira.FortinetSegmentName)
	for _, options := range in.List() {
		optionsCfg := options.(map[string]interface{})
		z := alkira.FortinetSegmentName{}

		var segment *alkira.Segment
		var zonesToGroups map[string][]string

		if v, ok := optionsCfg["segment_id"].(int); ok {
			sg, err := client.GetSegmentById(strconv.Itoa(v))
			if err != nil {
				return nil, err
			}
			segment = &sg
		}

		if v, ok := optionsCfg["zone"].(*schema.Set); ok {
			zonesToGroups = expandFortinetZone(v)
		}

		//TODO(mac): is this still required?
		//if zoneName == nil || segment == nil || groups == nil {
		//	return nil, errors.New("segment_option fields cannot be nil")
		//}

		z.SegmentId = segment.Id
		z.ZonesToGroups = zonesToGroups
		segmentOptions[segment.Name] = z
	}

	return segmentOptions, nil
}

func deflateFortinetSegmentOptions(c map[string]alkira.FortinetSegmentName) []map[string]interface{} {
	var options []map[string]interface{}

	for _, outerZoneToGroups := range c {
		for zone, groups := range outerZoneToGroups.ZonesToGroups {
			i := map[string]interface{}{
				"segment_id": outerZoneToGroups.SegmentId,
				"zone_name":  zone,
				"groups":     groups,
			}
			options = append(options, i)
		}
	}

	return options
}

func expandFortinetZone(in *schema.Set) map[string][]string {
	zonesToGroups := make(map[string][]string)

	for _, zone := range in.List() {
		zoneCfg := zone.(map[string]interface{})
		var name *string
		var groups []string

		if v, ok := zoneCfg["name"].(string); ok {
			name = &v
		}

		if v, ok := zoneCfg["groups"].([]interface{}); ok {
			groups = convertTypeListToStringList(v)
		}

		zonesToGroups[*name] = groups
	}

	return zonesToGroups
}
