package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraListAsPath() *schema.Resource {
	return &schema.Resource{
		Description: "This list could be used in a policy rule, a route " +
			"will match successfully if any one value from the list is " +
			"included within the AS-PATH of the route.",
		CreateContext: resourceListAsPath,
		ReadContext:   resourceListAsPathRead,
		UpdateContext: resourceListAsPathUpdate,
		DeleteContext: resourceListAsPathDelete,
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
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"values": {
				Description: "Value can be regular expression of AS PATH " +
					"or space sparated AS numbers. BGP regular expressions" +
					"are based on POSIX 1003.2 regular expressions.",
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceListAsPath(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListAsPath(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateListAsPathRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send request
	resource, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

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

	d.SetId(string(resource.Id))
	return resourceListAsPathRead(ctx, d, m)
}

func resourceListAsPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListAsPath(m.(*alkira.AlkiraClient))

	list, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("values", list.Values)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceListAsPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListAsPath(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateListAsPathRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
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

	return resourceListAsPathRead(ctx, d, m)
}

func resourceListAsPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListAsPath(m.(*alkira.AlkiraClient))

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

func generateListAsPathRequest(d *schema.ResourceData, m interface{}) (*alkira.List, error) {

	values := convertTypeListToStringList(d.Get("values").([]interface{}))

	request := &alkira.List{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      values,
	}

	return request, nil
}
