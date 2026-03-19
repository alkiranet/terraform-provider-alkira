package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/stretchr/testify/assert"
)

func ociVcnTestSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vcn_cidr": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vcn_subnet": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   {Type: schema.TypeString, Optional: true},
						"cidr": {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"vcn_route_table": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {Type: schema.TypeString, Optional: true},
						"prefix_list_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"options": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ADVERTISE_DEFAULT_ROUTE",
								"OVERRIDE_DEFAULT_ROUTE",
								"ADVERTISE_CUSTOM_PREFIX",
							}, false),
						},
					},
				},
			},
		},
	}
}

func TestSetConnectorOciVcnRouting_Nil(t *testing.T) {
	d := ociVcnTestSchema().TestResourceData()
	setConnectorOciVcnRouting(d, nil)

	assert.Empty(t, d.Get("vcn_cidr"))
	assert.Empty(t, d.Get("vcn_subnet"))
	assert.Empty(t, d.Get("vcn_route_table"))
}

func TestSetConnectorOciVcnRouting_InvalidType(t *testing.T) {
	d := ociVcnTestSchema().TestResourceData()
	// Pass a non-map type — should log error and return without panic
	setConnectorOciVcnRouting(d, "not-a-map")

	assert.Empty(t, d.Get("vcn_cidr"))
}

func TestSetConnectorOciVcnRouting_EmptyRouting(t *testing.T) {
	d := ociVcnTestSchema().TestResourceData()

	vcnRouting := map[string]interface{}{
		"exportToCXPOptions": map[string]interface{}{
			"userInputPrefixes": []interface{}{},
			"routeExportMode":   "USER_INPUT_PREFIXES",
		},
		"importFromCXPOptions": map[string]interface{}{
			"routeTables": []interface{}{},
		},
	}

	setConnectorOciVcnRouting(d, vcnRouting)

	assert.Empty(t, d.Get("vcn_cidr"))
	assert.Empty(t, d.Get("vcn_subnet"))
	assert.Empty(t, d.Get("vcn_route_table"))
}

func TestSetConnectorOciVcnRouting_WithCidr(t *testing.T) {
	d := ociVcnTestSchema().TestResourceData()

	vcnRouting := map[string]interface{}{
		"exportToCXPOptions": map[string]interface{}{
			"userInputPrefixes": []interface{}{
				map[string]interface{}{"type": "CIDR", "value": "10.0.0.0/16"},
				map[string]interface{}{"type": "CIDR", "value": "192.168.0.0/24"},
			},
			"routeExportMode": "USER_INPUT_PREFIXES",
		},
		"importFromCXPOptions": map[string]interface{}{
			"routeTables": []interface{}{},
		},
	}

	setConnectorOciVcnRouting(d, vcnRouting)

	cidr := d.Get("vcn_cidr").([]interface{})
	assert.Len(t, cidr, 2)
	assert.Contains(t, cidr, "10.0.0.0/16")
	assert.Contains(t, cidr, "192.168.0.0/24")
	assert.Empty(t, d.Get("vcn_subnet"))
}

func TestSetConnectorOciVcnRouting_WithSubnets(t *testing.T) {
	d := ociVcnTestSchema().TestResourceData()

	vcnRouting := map[string]interface{}{
		"exportToCXPOptions": map[string]interface{}{
			"userInputPrefixes": []interface{}{
				map[string]interface{}{
					"id":    "ocid1.subnet.oc1..aaa",
					"type":  "SUBNET",
					"value": "172.16.0.0/29",
				},
			},
			"routeExportMode": "USER_INPUT_PREFIXES",
		},
		"importFromCXPOptions": map[string]interface{}{
			"routeTables": []interface{}{},
		},
	}

	setConnectorOciVcnRouting(d, vcnRouting)

	assert.Empty(t, d.Get("vcn_cidr"))
	subnets := d.Get("vcn_subnet").(*schema.Set).List()
	assert.Len(t, subnets, 1)
	subnet := subnets[0].(map[string]interface{})
	assert.Equal(t, "ocid1.subnet.oc1..aaa", subnet["id"])
	assert.Equal(t, "172.16.0.0/29", subnet["cidr"])
}

func TestSetConnectorOciVcnRouting_WithRouteTables(t *testing.T) {
	d := ociVcnTestSchema().TestResourceData()

	vcnRouting := map[string]interface{}{
		"exportToCXPOptions": map[string]interface{}{
			"userInputPrefixes": []interface{}{},
			"routeExportMode":   "USER_INPUT_PREFIXES",
		},
		"importFromCXPOptions": map[string]interface{}{
			"routeTables": []interface{}{
				map[string]interface{}{
					"id":              "ocid1.routetable.oc1..aaa",
					"routeImportMode": "ADVERTISE_CUSTOM_PREFIX",
					"prefixListIds":   []interface{}{float64(1), float64(2)},
				},
			},
		},
	}

	setConnectorOciVcnRouting(d, vcnRouting)

	tables := d.Get("vcn_route_table").(*schema.Set).List()
	assert.Len(t, tables, 1)
	table := tables[0].(map[string]interface{})
	assert.Equal(t, "ocid1.routetable.oc1..aaa", table["id"])
	assert.Equal(t, "ADVERTISE_CUSTOM_PREFIX", table["options"])
	assert.ElementsMatch(t, []int{1, 2}, table["prefix_list_ids"])
}

func TestSetConnectorOciVcnRouting_CidrAndRouteTables(t *testing.T) {
	d := ociVcnTestSchema().TestResourceData()

	vcnRouting := map[string]interface{}{
		"exportToCXPOptions": map[string]interface{}{
			"userInputPrefixes": []interface{}{
				map[string]interface{}{"type": "CIDR", "value": "10.0.0.0/16"},
			},
			"routeExportMode": "USER_INPUT_PREFIXES",
		},
		"importFromCXPOptions": map[string]interface{}{
			"routeTables": []interface{}{
				map[string]interface{}{
					"id":              "ocid1.routetable.oc1..bbb",
					"routeImportMode": "ADVERTISE_DEFAULT_ROUTE",
					"prefixListIds":   []interface{}{},
				},
			},
		},
	}

	setConnectorOciVcnRouting(d, vcnRouting)

	cidr := d.Get("vcn_cidr").([]interface{})
	assert.Len(t, cidr, 1)
	assert.Equal(t, "10.0.0.0/16", cidr[0])

	tables := d.Get("vcn_route_table").(*schema.Set).List()
	assert.Len(t, tables, 1)
	table := tables[0].(map[string]interface{})
	assert.Equal(t, "ocid1.routetable.oc1..bbb", table["id"])
	assert.Equal(t, "ADVERTISE_DEFAULT_ROUTE", table["options"])
}
