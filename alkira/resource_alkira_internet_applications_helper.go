package alkira

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getInternetApplicationGroup(client *alkira.AlkiraClient) int {
	api := alkira.NewGroup(client)
	groups, err := api.GetAll()

	if err != nil {
		log.Printf("[ERROR] failed to get groups")
		return 0
	}

	var result []alkira.Group
	json.Unmarshal([]byte(groups), &result)

	for _, group := range result {
		if group.Name == "ALK-INB-INT-GROUP" {

			groupId, _ := strconv.Atoi(string(group.Id))
			return groupId
		}
	}

	return 0
}

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
