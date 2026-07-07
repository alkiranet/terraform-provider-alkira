package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func vhubTestSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vhub_routing": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_import_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefix_list_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
		},
	}
}

func TestConstructVhubRouting_EmptyBlock(t *testing.T) {
	d := vhubTestSchema().TestResourceData()
	d.Set("vhub_routing", []interface{}{
		map[string]interface{}{
			"route_import_mode": "",
			"prefix_list_ids":   []interface{}{},
		},
	})

	result := constructVhubRouting(d)

	assert.NotNil(t, result)
	assert.Nil(t, result.ImportOptions, "empty block should produce nil ImportOptions")
}

func TestConstructVhubRouting_WithAdvertiseDefaultRoute(t *testing.T) {
	d := vhubTestSchema().TestResourceData()
	d.Set("vhub_routing", []interface{}{
		map[string]interface{}{
			"route_import_mode": "ADVERTISE_DEFAULT_ROUTE",
			"prefix_list_ids":   []interface{}{},
		},
	})

	result := constructVhubRouting(d)

	assert.NotNil(t, result)
	assert.NotNil(t, result.ImportOptions)
	assert.Equal(t, "ADVERTISE_DEFAULT_ROUTE", result.ImportOptions.RouteImportMode)
	assert.Empty(t, result.ImportOptions.PrefixListIds)
}

func TestConstructVhubRouting_WithCustomPrefix(t *testing.T) {
	d := vhubTestSchema().TestResourceData()
	d.Set("vhub_routing", []interface{}{
		map[string]interface{}{
			"route_import_mode": "ADVERTISE_CUSTOM_PREFIX",
			"prefix_list_ids":   []interface{}{10, 20, 30},
		},
	})

	result := constructVhubRouting(d)

	assert.NotNil(t, result)
	assert.NotNil(t, result.ImportOptions)
	assert.Equal(t, "ADVERTISE_CUSTOM_PREFIX", result.ImportOptions.RouteImportMode)
	assert.Equal(t, []int{10, 20, 30}, result.ImportOptions.PrefixListIds)
}

func TestConstructVhubRouting_EmptyList(t *testing.T) {
	d := vhubTestSchema().TestResourceData()
	d.Set("vhub_routing", []interface{}{})

	result := constructVhubRouting(d)

	assert.Nil(t, result)
}

func TestSetVhubRouting_NilRouting(t *testing.T) {
	d := vhubTestSchema().TestResourceData()

	setVhubRouting(d, nil)

	val := d.Get("vhub_routing").([]interface{})
	assert.Empty(t, val)
}

func TestSetVhubRouting_NilImportOptions(t *testing.T) {
	d := vhubTestSchema().TestResourceData()

	setVhubRouting(d, &alkira.ConnectorVhubRouting{})

	val := d.Get("vhub_routing").([]interface{})
	assert.Len(t, val, 1)

	block := val[0].(map[string]interface{})
	assert.Equal(t, "", block["route_import_mode"])
	assert.Empty(t, block["prefix_list_ids"])
}

func TestSetVhubRouting_WithImportOptions(t *testing.T) {
	d := vhubTestSchema().TestResourceData()

	routing := &alkira.ConnectorVhubRouting{
		ImportOptions: &alkira.ConnectorVhubImportOptions{
			RouteImportMode: "ADVERTISE_CUSTOM_PREFIX",
			PrefixListIds:   []int{5, 10},
		},
	}

	setVhubRouting(d, routing)

	val := d.Get("vhub_routing").([]interface{})
	assert.Len(t, val, 1)

	block := val[0].(map[string]interface{})
	assert.Equal(t, "ADVERTISE_CUSTOM_PREFIX", block["route_import_mode"])
	assert.Equal(t, []interface{}{5, 10}, block["prefix_list_ids"])
}

func TestVhubSchemaImmutableFields(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorAzureVhub().Schema

	immutableFields := []string{
		"virtual_hub_id",
	}

	for _, field := range immutableFields {
		t.Run(field+"_has_ForceNew", func(t *testing.T) {
			s, exists := resourceSchema[field]
			assert.True(t, exists, "field %q must exist in schema", field)
			assert.True(t, s.ForceNew, "field %q must have ForceNew: true", field)
		})
	}
}

func TestVhubSchemaMutableFields(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorAzureVhub().Schema

	mutableFields := []string{
		"name",
		"credential_id",
		"cxp",
		"description",
		"enabled",
		"group",
		"scale_group_id",
		"segment_id",
		"size",
		"billing_tag_ids",
		"vhub_routing",
	}

	for _, field := range mutableFields {
		t.Run(field+"_no_ForceNew", func(t *testing.T) {
			s, exists := resourceSchema[field]
			assert.True(t, exists, "field %q must exist in schema", field)
			assert.False(t, s.ForceNew, "field %q must not have ForceNew", field)
		})
	}
}

func TestVhubSchemaRouteImportModeValues(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorAzureVhub().Schema

	vhubRouting := resourceSchema["vhub_routing"]
	routingSchema := vhubRouting.Elem.(*schema.Resource).Schema
	routeImportMode := routingSchema["route_import_mode"]

	// Validate that OVERRIDE_DEFAULT_ROUTE is not accepted
	validateFunc := routeImportMode.ValidateFunc
	assert.NotNil(t, validateFunc)

	_, errs := validateFunc("ADVERTISE_DEFAULT_ROUTE", "route_import_mode")
	assert.Empty(t, errs)

	_, errs = validateFunc("ADVERTISE_CUSTOM_PREFIX", "route_import_mode")
	assert.Empty(t, errs)

	_, errs = validateFunc("OVERRIDE_DEFAULT_ROUTE", "route_import_mode")
	assert.NotEmpty(t, errs, "OVERRIDE_DEFAULT_ROUTE should be rejected")
}
