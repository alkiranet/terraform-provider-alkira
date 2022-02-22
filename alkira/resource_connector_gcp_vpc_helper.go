package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	//"github.com/zclconf/go-cty/cty"
	"github.com/hashicorp/go-cty/cty"
)

func convertGcpRouting(in *schema.Set) *alkira.ConnectorGcpVpcRouting {
	if in == nil || in.Len() > 1 {
		log.Printf("[DEBUG] Only one object allowed in gcp routing options")
		return nil
	}

	if in.Len() < 1 {
		return nil
	}

	gcp := &alkira.ConnectorGcpVpcRouting{
		ImportOptions: alkira.ConnectorGcpVpcImportOptions{},
	}

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})

		if v, ok := cfg["prefix_list_ids"].([]interface{}); ok {
			gcp.ImportOptions.PrefixListIds = convertTypeListToIntList(v)
		}

		if v, ok := cfg["custom_prefix"].(string); ok {
			gcp.ImportOptions.RouteImportMode = v
		}
	}

	return gcp
}

func validateCustomPrefix(v interface{}, p cty.Path) diag.Diagnostics {
	value := v.(string)
	var diags diag.Diagnostics

	if value == "ADVERTISE_CUSTOM_PREFIX" {
		return diags
	}

	if value == "ADVERTISE_DEFAULT_ROUTE" {
		return diags
	}

	diag := diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Invalid custom_prefix value. Valid custom_prefix options are ADVERTISE_CUSTOM_PREFIX ADVERTISE_DEFAULT_ROUTE.",
		Detail:   fmt.Sprintf("Invalid custom_prefix value: %q. Only ADVERTISE_CUSTOM_PREFIX or ADVERTISE_DEFAULT_ROUTE are accepted values.", value),
	}
	diags = append(diags, diag)

	return diags
}
