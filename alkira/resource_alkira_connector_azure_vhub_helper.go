package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setVhubRouting sets the vhub_routing block value in Terraform state
func setVhubRouting(d *schema.ResourceData, routing *alkira.ConnectorVhubRouting) {
	if routing == nil {
		return
	}

	vhubRouting := map[string]interface{}{
		"route_import_mode": "",
		"prefix_list_ids":   nil,
	}

	if routing.ImportOptions != nil {
		vhubRouting["route_import_mode"] = routing.ImportOptions.RouteImportMode
		vhubRouting["prefix_list_ids"] = routing.ImportOptions.PrefixListIds
	}

	d.Set("vhub_routing", []interface{}{vhubRouting})
}

// constructVhubRouting constructs the ConnectorAzureVhubRouting from Terraform resource data
func constructVhubRouting(d *schema.ResourceData) *alkira.ConnectorVhubRouting {
	routingList := d.Get("vhub_routing").([]interface{})

	if len(routingList) == 0 || routingList[0] == nil {
		return nil
	}

	routing := routingList[0].(map[string]interface{})

	if routing["route_import_mode"].(string) == "" {
		return &alkira.ConnectorVhubRouting{}
	}

	importOptions := &alkira.ConnectorVhubImportOptions{
		RouteImportMode: routing["route_import_mode"].(string),
		PrefixListIds:   convertTypeListToIntList(routing["prefix_list_ids"].([]interface{})),
	}

	return &alkira.ConnectorVhubRouting{
		ImportOptions: importOptions,
	}
}
