package alkira

import (
	"reflect"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceConnectorGcpVpcReadSetGcpRoutingOptions(t *testing.T) {
	expectedPrefixIds := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	expectedCustomPrefix := "CUSTOM"

	c := &alkira.ConnectorGcpVpcRouting{}
	c.ImportOptions.RouteImportMode = expectedCustomPrefix
	c.ImportOptions.PrefixListIds = expectedPrefixIds

	r := resourceAlkiraConnectorGcpVpc()
	d := r.TestResourceData()

	setGcpRoutingOptions(c, d)
	list := d.Get("gcp_routing").(*schema.Set).List()
	m := list[0].(map[string]interface{})

	actualCustomPrefixIds := m["custom_prefix"].(string)
	assertStrEquals(t, actualCustomPrefixIds, expectedCustomPrefix)

	actualCustomPrefix := convertTypeListToIntList(m["prefix_list_ids"].([]interface{}))
	assertTrue(t, reflect.DeepEqual(actualCustomPrefix, expectedPrefixIds), "Failed: Prefixes do not match")
}
