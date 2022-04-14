package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


// constructVnetRouting expand AZURE VNET routing options
func constructVnetRouting(option string, prefixList []interface{}) *alkira.ConnectorVnetRouting {

	routing := alkira.ConnectorVnetImportOptions{}

	routing.RouteImportMode = option
	routing.PrefixListIds = convertTypeListToIntList(prefixList)

	return &alkira.ConnectorVnetRouting{routing}
}


// expandConnectorAzureVnetServiceRoute expand service route of connector_azure_vnet
func expandConnectorAzureVnetServiceRoute(in *schema.Set) []*alkira.ConnectorAzureVnet {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] empty IPSec endpoint input")
		return nil
	}
