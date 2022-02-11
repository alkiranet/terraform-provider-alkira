package alkira

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type panZone struct {
	Segment string
	Zone    string
	Groups  interface{}
}

func getInternetApplicationGroup(client *alkira.AlkiraClient) int {
	groups, err := client.GetGroups()

	if err != nil {
		log.Printf("[ERROR] failed to get groups")
		return 0
	}

	var result []alkira.Group
	json.Unmarshal([]byte(groups), &result)

	for _, group := range result {
		if group.Name == "ALK-INB-INT-GROUP" {
			return group.Id
		}
	}

	return 0
}

func convertTypeListToIntList(in []interface{}) []int {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty TypeList to convert to IntList")
		return nil
	}

	intList := make([]int, len(in))

	for i, value := range in {
		intList[i] = value.(int)
	}

	return intList
}

func convertTypeListToStringList(in []interface{}) []string {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty TypeList to convert to StringList")
		return nil
	}

	strList := make([]string, len(in))

	for i, value := range in {
		strList[i] = value.(string)
	}

	return strList
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

func expandPolicyRuleListRules(in *schema.Set) []alkira.PolicyRuleListRule {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid policy rule")
		return nil
	}

	rules := make([]alkira.PolicyRuleListRule, in.Len())
	for i, rule := range in.List() {
		r := alkira.PolicyRuleListRule{}
		ruleCfg := rule.(map[string]interface{})
		if v, ok := ruleCfg["priority"].(int); ok {
			r.Priority = v
		}
		if v, ok := ruleCfg["rule_id"].(int); ok {
			r.RuleId = v
		}
		rules[i] = r
	}

	return rules
}

// expandAwsVpcRouteTables expand AWS-VPC route tables
func expandAwsVpcRouteTables(in *schema.Set) []alkira.RouteTables {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty VPC route table input")
		return []alkira.RouteTables{}
	}

	tables := make([]alkira.RouteTables, in.Len())
	for i, table := range in.List() {
		r := alkira.RouteTables{}
		t := table.(map[string]interface{})
		if v, ok := t["id"].(string); ok {
			r.Id = v
		}
		if v, ok := t["options"].(string); ok {
			r.Mode = v
		}

		r.PrefixListIds = convertTypeListToIntList(t["prefix_list_ids"].([]interface{}))
		tables[i] = r
	}

	return tables
}

// generateUserInputPrefixes generate UserInputPrefixes used in AWS-VPC connector
func generateUserInputPrefixes(cidr []interface{}, subnets *schema.Set) ([]alkira.InputPrefixes, error) {

	if len(cidr) == 0 && subnets == nil {
		return nil, fmt.Errorf("ERROR: either \"vpc_subnet\" or \"vpc_cidr\" must be specified.")
	}

	// Processing "vpc_cidr"
	if len(cidr) > 0 {
		log.Printf("[DEBUG] Processing vpc_cidr %v", cidr)
		cidrList := make([]alkira.InputPrefixes, len(cidr))

		for i, value := range cidr {
			cidrList[i].Value = value.(string)
			cidrList[i].Type = "CIDR"
		}

		return cidrList, nil
	}

	// Processing VPC subnets
	log.Printf("[DEBUG] Processing vpc_subnet")
	if subnets == nil || subnets.Len() == 0 {
		log.Printf("[DEBUG] Empty vpc_subnet")
		return nil, fmt.Errorf("ERROR: Invalid vpc_subnet.")
	}

	prefixes := make([]alkira.InputPrefixes, subnets.Len())
	for i, subnet := range subnets.List() {
		r := alkira.InputPrefixes{}
		t := subnet.(map[string]interface{})
		if v, ok := t["id"].(string); ok {
			r.Id = v
		}
		if v, ok := t["cidr"].(string); ok {
			r.Value = v
		}

		r.Type = "SUBNET"
		prefixes[i] = r
	}

	return prefixes, nil
}
