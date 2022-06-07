package alkira

import (
	"errors"
	"log"

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

func expandFortinetSegmentOptions(in *schema.Set, m interface{}) (map[string]*alkira.GlobalProtectSegmentName, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || in.Len() == 0 {
		return nil, errors.New("invalid Fortinet Segment Options input")
	}

	sgmtOptions := make(map[string]*alkira.GlobalProtectSegmentName)
	for _, sgmtOption := range in.List() {
		r := &alkira.GlobalProtectSegmentName{}
		segmentCfg := sgmtOption.(map[string]interface{})
		var segmentName string

		if v, ok := segmentCfg["segment_id"].(string); ok {
			segment, err := client.GetSegmentById(v)
			if err != nil {
				return nil, err
			}
			segmentName = segment.Name
		}
		if v, ok := segmentCfg["remote_user_zone_name"].(string); ok {
			r.RemoteUserZoneName = v
		}
		if v, ok := segmentCfg["portal_fqdn_prefix"].(string); ok {
			r.PortalFqdnPrefix = v
		}
		if v, ok := segmentCfg["service_group_name"].(string); ok {
			r.ServiceGroupName = v
		}

		sgmtOptions[segmentName] = r
	}

	return sgmtOptions, nil
}
