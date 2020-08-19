package alkira

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

type panZone struct {
	Segment string
	Zone    string
	Groups  interface{}
}


func getInternetApplicationGroup(client *alkira.AlkiraClient) (int) {
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

		zonesToGroups  := make(map[string]interface{})
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


// func expandIPSecSegmentOptions(in *schema.Set) []map[string]interface{} {
//     if in == nil || in.Len() == 0 {
// 		log.Printf("[DEBUG] invalid IPSec segmentOption input")
//         return nil
//     }

//     sites := make([]alkira.ConnectorIPSecSite, in.Len())
//     for i, site := range in.List() {
//         r := alkira.ConnectorIPSecSite{}
// 		siteCfg := site.(map[string]interface{})
//         if v, ok := siteCfg["name"].(string); ok {
// 			r.Name = v
//         }
//         if v, ok := siteCfg["customer_gateway_asn"].(string); ok {
//             r.CustomerGwAsn = v
//         }
//         if v, ok := siteCfg["customer_gateway_ip"].(string); ok {
//             r.CustomerGwIp = v
//         }
//         if v, ok := siteCfg["preshared_keys"].([]string); ok {
// 			if len(v) != 0 {
// 				r.PresharedKeys = v
// 			} else {
// 				r.PresharedKeys[0] = "[],[]"
// 			}
//         }

//         sites[i] = r
//     }

//     return sites
// }
