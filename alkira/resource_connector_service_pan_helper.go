package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type panZone struct {
	Segment string
	Zone    string
	Groups  interface{}
}

func expandGlobalProtectSegmentOptions(in *schema.Set) map[string]*alkira.GlobalProtectSegmentName {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid Global Protect Segment Options input")
		return nil
	}

	sgmtOptions := make(map[string]*alkira.GlobalProtectSegmentName)
	for _, sgmtOption := range in.List() {
		r := &alkira.GlobalProtectSegmentName{}
		segmentCfg := sgmtOption.(map[string]interface{})
		segmentName := ""

		if v, ok := segmentCfg["segment_name"].(string); ok {
			segmentName = v
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

	return sgmtOptions
}

func expandGlobalProtectSegmentOptionsInstance(in *schema.Set) map[string]*alkira.GlobalProtectSegmentNameInstance {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid input for Global Pan Protect Options for service PAN instance")
		return nil
	}

	sgmtOptions := make(map[string]*alkira.GlobalProtectSegmentNameInstance)
	for _, sgmtOption := range in.List() {
		r := &alkira.GlobalProtectSegmentNameInstance{}
		segmentCfg := sgmtOption.(map[string]interface{})
		segmentName := ""

		if v, ok := segmentCfg["segment_name"].(string); ok {
			segmentName = v
		}
		if v, ok := segmentCfg["portal_enabled"].(bool); ok {
			r.PortalEnabled = v
		}
		if v, ok := segmentCfg["gateway_enabled"].(bool); ok {
			r.GatewayEnabled = v
		}
		if v, ok := segmentCfg["prefix_list_id"].(int); ok {
			r.PrefixListId = v
		}

		sgmtOptions[segmentName] = r
	}

	return sgmtOptions
}

func expandPanSegmentOptions(in *schema.Set) map[string]interface{} {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid SegmentOptions input")
		return nil
	}

	zoneMap := make([]panZone, in.Len())

	for i, option := range in.List() {
		r := panZone{}
		cfg := option.(map[string]interface{})
		if v, ok := cfg["segment_name"].(string); ok {
			r.Segment = v
		}
		if v, ok := cfg["zone_name"].(string); ok {
			r.Zone = v
		}

		r.Groups = cfg["groups"]

		zoneMap[i] = r
	}

	segmentOptions := make(map[string]interface{})

	for _, x := range zoneMap {
		zone := make(map[string]interface{})
		zone[x.Zone] = x.Groups

		for _, y := range zoneMap {
			if x.Segment == y.Segment {
				zone[y.Zone] = y.Groups
			}
		}

		zonesToGroups := make(map[string]interface{})
		zonesToGroups["zonesToGroups"] = zone

		segmentOptions[x.Segment] = zonesToGroups
	}

	return segmentOptions
}
func expandPanInstances(in *schema.Set) []alkira.ServicePanInstance {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid PAN instance input")
		return nil
	}

	instances := make([]alkira.ServicePanInstance, in.Len())
	for i, instance := range in.List() {
		r := alkira.ServicePanInstance{}
		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			r.CredentialId = v
		}
		if v, ok := instanceCfg["global_protect_segment_options"].(*schema.Set); ok {
			r.GlobalProtectSegmentOptions = expandGlobalProtectSegmentOptionsInstance(v)
		}
		instances[i] = r
	}

	return instances
}
