package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraListUdr() *schema.Resource {
	return &schema.Resource{
		Description:   "User Defined Routes (UDR) list.",
		CreateContext: resourceListUdr,
		ReadContext:   resourceListUdrRead,
		UpdateContext: resourceListUdrUpdate,
		DeleteContext: resourceListUdrDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description for the list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cloud_provider": {
				Description: "Cloud provider. Only `AZURE` is supported for " +
					"now.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "AZURE",
				ValidateFunc: validation.StringInSlice([]string{"AZURE"}, false),
			},
			"route": {
				Description: "ID of `list_dns_server` resource.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Description: "Description for the route.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"prefix": {
							Description: "The prefix of the route. This " +
								"prefix must be in the CIDR format " +
								"(`x.x.x.x/mask`). The mask can be between " +
								"`8-32`.",
							Type:     schema.TypeString,
							Required: true,
						},
						// "next_hop_type": {
						// 	Description: "The next hop type. Value could " +
						// 		"`INTERNET` or empty.",
						// 	Type:     schema.TypeString,
						// 	Optional: true,
						// },
						// "next_hop_value": {
						// 	Description: "The next hop value of the route. " +
						// 		"When `next_hope_type` is defined as " +
						// 		"`INTERNET`, then this must be " +
						// 		"left empty.",
						// 	Type:     schema.TypeString,
						// 	Required: true,
						// },
					},
				},
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceListUdr(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewUdrList(m.(*alkira.AlkiraClient))

	// Construct requst
	request := generateListUdrRequest(d, m)

	// Send request
	resource, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(resource.Id))

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceListUdrRead(ctx, d, m)
}

func resourceListUdrRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewUdrList(m.(*alkira.AlkiraClient))

	list, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("cloud_provider", list.CloudProvider)
	d.Set("route", list.Udrs)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceListUdrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewUdrList(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateListUdrRequest(d, m)

	// Send request
	provState, err, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceListUdrRead(ctx, d, m)
}

func resourceListUdrDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewUdrList(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateListUdrRequest(d *schema.ResourceData, m interface{}) *alkira.UdrList {

	routes := expandListUdrRoutes(d.Get("route").(*schema.Set))

	request := &alkira.UdrList{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		CloudProvider: d.Get("cloud_provider").(string),
		Udrs:          routes,
	}

	return request
}

func expandListUdrRoutes(in *schema.Set) []alkira.UdrListUdrs {

	if in == nil || in.Len() == 0 {
		return nil
	}

	routes := make([]alkira.UdrListUdrs, in.Len())
	for i, route := range in.List() {
		r := alkira.UdrListUdrs{}
		routeCfg := route.(map[string]interface{})
		if v, ok := routeCfg["prefix"].(string); ok {
			r.Prefix = v
		}
		if v, ok := routeCfg["description"].(string); ok {
			r.Description = v
		}

		//
		// For now, those two fields have fixed value
		//
		r.NextHopType = "INTERNET"
		r.NextHopValue = ""

		routes[i] = r
	}

	return routes
}
