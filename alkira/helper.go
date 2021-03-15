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
		log.Printf("[DEBUG] empty input")
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
		log.Printf("[DEBUG] empty input")
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

func expandPanInstances(in *schema.Set) []alkira.ServicePanInstance {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid IPSec site input")
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
		instances[i] = r
	}

	return instances
}

func expandIPSecSites(in *schema.Set) []alkira.ConnectorIPSecSite {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid IPSec site input")
		return nil
	}

	sites := make([]alkira.ConnectorIPSecSite, in.Len())
	for i, site := range in.List() {
		r := alkira.ConnectorIPSecSite{}
		siteCfg := site.(map[string]interface{})
		if v, ok := siteCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := siteCfg["customer_gateway_asn"].(string); ok {
			r.CustomerGwAsn = v
		}
		if v, ok := siteCfg["customer_gateway_ip"].(string); ok {
			r.CustomerGwIp = v
		}
		if v, ok := siteCfg["preshared_keys"].([]string); ok {
			if len(v) != 0 {
				r.PresharedKeys = v
			} else {
				r.PresharedKeys[0] = "[],[]"
			}
		}

		sites[i] = r
	}

	return sites
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
func generateUserInputPrefixes(cidr string, subnets []interface{}) ([]alkira.InputPrefixes, error) {

	if cidr == "" && subnets == nil {
		return nil, fmt.Errorf("ERROR: either vpc_subnets or vpc_cidr must be specified.")
	}

	if cidr != "" && len(subnets) > 0 {
		return nil, fmt.Errorf("ERROR: vpc_subnets and vpc_cidr can't be specified at the same time.")
	}

	// Processing VPC CIDR
	if cidr != "" {
		log.Printf("[DEBUG] Processing VPC CIDR")
		inputPrefix := alkira.InputPrefixes{
			Id:    "",
			Type:  "CIDR",
			Value: cidr,
		}
		return []alkira.InputPrefixes{inputPrefix}, nil
	}

	// Processing VPC subnets
	log.Printf("[DEBUG] Processing VPC Subnets")
	prefixes := make([]alkira.InputPrefixes, len(subnets))
	for i, subnet := range subnets {
		r := alkira.InputPrefixes{}
		r.Id = subnet.(string)
		r.Type = "SUBNET"
		prefixes[i] = r
	}

	return prefixes, nil
}
