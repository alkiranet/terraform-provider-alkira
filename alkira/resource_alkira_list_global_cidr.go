package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraListGlobalCidr() *schema.Resource {
	return &schema.Resource{
		Description:   "A list of CIDRs to be used for services.",
		CreateContext: resourceListGlobalCidr,
		ReadContext:   resourceListGlobalCidrRead,
		UpdateContext: resourceListGlobalCidrUpdate,
		DeleteContext: resourceListGlobalCidrDelete,
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
			"cxp": {
				Description: "CXP the list belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"values": {
				Description: "A list of CIDRs, The CIDR must be `/24` and a " +
					"subnet of the following: `10.0.0.0/18`, `172.16.0.0/12`, " +
					"`192.168.0.0/16`, `100.64.0.0/10`.",
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Description: "A list of associated service types.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceListGlobalCidr(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	// Construct requst
	request := generateListGlobalCidrRequest(d, m)

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
				Summary:  "PROVISION FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceListGlobalCidrRead(ctx, d, m)
}

func resourceListGlobalCidrRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	list, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("cxp", list.CXP)
	d.Set("values", list.Values)
	d.Set("tags", list.Tags)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceListGlobalCidrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateListGlobalCidrRequest(d, m)

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
				Summary:  "PROVISION FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceListGlobalCidrRead(ctx, d, m)
}

func resourceListGlobalCidrDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision == true && provState != "SUCCESS" {
		return diag.FromErr(provErr)
	}

	d.SetId("")
	return nil
}

func generateListGlobalCidrRequest(d *schema.ResourceData, m interface{}) *alkira.GlobalCidrList {

	request := &alkira.GlobalCidrList{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CXP:         d.Get("cxp").(string),
		Values:      convertTypeListToStringList(d.Get("values").([]interface{})),
		Tags:        convertTypeListToStringList(d.Get("tags").([]interface{})),
	}

	return request
}
