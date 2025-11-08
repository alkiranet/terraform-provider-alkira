package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UNUSED: Commented out to suppress linter warnings
// func getInternetApplicationGroup(client *alkira.AlkiraClient) int {
// 	api := alkira.NewGroup(client)
// 	groups, err := api.GetAll()
//
// 	if err != nil {
// 		log.Printf("[ERROR] failed to get groups")
// 		return 0
// 	}
//
// 	var result []alkira.Group
// 	json.Unmarshal([]byte(groups), &result)
//
// 	for _, group := range result {
// 		if group.Name == "ALK-INB-INT-GROUP" {
//
// 			groupId, _ := strconv.Atoi(string(group.Id))
// 			return groupId
// 		}
// 	}
//
// 	return 0
// }

func expandInternetApplicationTargets(in *schema.Set) []alkira.InternetApplicationTargets {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid internet application targets")
		return nil
	}

	targets := make([]alkira.InternetApplicationTargets, in.Len())

	for i, target := range in.List() {
		r := alkira.InternetApplicationTargets{}
		content := target.(map[string]interface{})
		if v, ok := content["type"].(string); ok {
			r.Type = v
		}
		if v, ok := content["value"].(string); ok {
			r.Value = v
		}
		if v, ok := content["port_ranges"].([]interface{}); ok {
			r.PortRanges = convertTypeListToStringList(v)
		}
		targets[i] = r
	}

	return targets
}

func expandInternetApplicationSourceNatPool(in *schema.Set) []*alkira.InternetApplicationSnatIpv4 {
	if in == nil {
		return nil
	}

	pool := make([]*alkira.InternetApplicationSnatIpv4, in.Len())

	for i, ips := range in.List() {
		r := alkira.InternetApplicationSnatIpv4{}
		content := ips.(map[string]interface{})
		if v, ok := content["start_ip"].(string); ok {
			r.StartIp = v
		}
		if v, ok := content["end_ip"].(string); ok {
			r.EndIp = v
		}
		pool[i] = &r
	}

	return pool
}
