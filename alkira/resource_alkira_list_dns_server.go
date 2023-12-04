package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraListDnsServer() *schema.Resource {
	return &schema.Resource{
		Description:   "A list of DNS servers.",
		CreateContext: resourceListDnsServer,
		ReadContext:   resourceListDnsServerRead,
		UpdateContext: resourceListDnsServerUpdate,
		DeleteContext: resourceListDnsServerDelete,
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
			"dns_server_ips": {
				Description: "DNS server IPs. The IP can't be `any` and can't " +
					"be an IP from the following CIDRs: `0.0.0.0/8`, " +
					"`127.0.0.0/8`, `169.254.0.0/16`, `224.0.0.0/4`, " +
					"`240.0.0.0/4`, `255.255.255.255/32`.",
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"segment_id": {
				Description: "The segment that is associated with the list.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceListDnsServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewDnsServerList(m.(*alkira.AlkiraClient))

	// Construct requst
	request, reqErr := generateListDnsServerRequest(d, m)

	if reqErr != nil {
		return diag.FromErr(reqErr)
	}

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

	return resourceListDnsServerRead(ctx, d, m)
}

func resourceListDnsServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewDnsServerList(m.(*alkira.AlkiraClient))

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
	d.Set("dns_server_ips", list.DnsServerIps)

	//
	// Get segemnt
	//
	segmentId, err := getSegmentIdByName(list.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("segment_id", segmentId)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceListDnsServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewDnsServerList(m.(*alkira.AlkiraClient))

	// Construct request
	request, reqErr := generateListDnsServerRequest(d, m)

	if reqErr != nil {
		return diag.FromErr(reqErr)
	}

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

	return resourceListDnsServerRead(ctx, d, m)
}

func resourceListDnsServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewDnsServerList(m.(*alkira.AlkiraClient))

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

func generateListDnsServerRequest(d *schema.ResourceData, m interface{}) (*alkira.DnsServerList, error) {
	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	request := &alkira.DnsServerList{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		DnsServerIps: convertTypeSetToStringList(d.Get("dns_server_ips").(*schema.Set)),
		Segment:      segmentName,
	}

	return request, nil
}
