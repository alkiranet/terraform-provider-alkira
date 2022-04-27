package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructVnetRouting construct connector_azure_vnet routing options
func constructVnetRouting(d *schema.ResourceData) (*alkira.ConnectorVnetRouting, error) {

	vnetRouting := alkira.ConnectorVnetRouting{}

	importOptions := alkira.ConnectorVnetImportOptions{}
	importOptions.RouteImportMode = d.Get("routing_options").(string)
	importOptions.PrefixListIds = convertTypeListToIntList(d.Get("routing_prefix_list_ids").([]interface{}))

	serviceRoutes, err := constructVnetServiceRoutes(d.Get("service_route").(*schema.Set))

	if err != nil {
		return nil, err
	}

	vnetRouting.ImportOptions = importOptions
	vnetRouting.ServiceRoutes = *serviceRoutes

	return &vnetRouting, nil
}

// constructVnetServiceRoutes expand service route of connector_azure_vnet
func constructVnetServiceRoutes(in *schema.Set) (*alkira.ConnectorVnetServiceRoutes, error) {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] empty service_route of connector_azure_vnet")
		return nil, nil
	}

	serviceRoutes := alkira.ConnectorVnetServiceRoutes{}

	for _, input := range in.List() {
		serviceRouteInput := input.(map[string]interface{})

		switch serviceRouteType := serviceRouteInput["type"].(string); serviceRouteType {
		case "CIDR":
			{
				cidr := alkira.ConnectorVnetServiceRoute{}

				if v, ok := serviceRouteInput["value"].(string); ok {
					cidr.Value = v
				}

				cidr.ServiceTags = convertTypeListToStringList(serviceRouteInput["service_tags"].([]interface{}))

				serviceRoutes.Cidrs = append(serviceRoutes.Cidrs, cidr)
			}
		case "SUBNET":
			{
				subnet := alkira.ConnectorVnetServiceRoute{}

				if v, ok := serviceRouteInput["subnet_id"].(string); ok {
					subnet.Id = v
				} else {
					return nil, fmt.Errorf("ERROR: subnet_id is required if type of service routes is SUBNET.")
				}

				if v, ok := serviceRouteInput["value"].(string); ok {
					subnet.Value = v
				}

				subnet.ServiceTags = convertTypeListToStringList(serviceRouteInput["service_tags"].([]interface{}))

				serviceRoutes.Subnets = append(serviceRoutes.Subnets, subnet)
			}
		default:
			return nil, fmt.Errorf("ERROR: invalid routing type")
		}
	}

	return &serviceRoutes, nil

}
