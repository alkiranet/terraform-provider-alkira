package alkira

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraF5vServerEndpoint() *schema.Resource {
	return &schema.Resource{
		Description:   "",
		CreateContext: resourceF5vServerEndpointCreate,
		ReadContext:   resourceF5vServerEndpointRead,
		UpdateContext: resourceF5vServerEndpointUpdate,
		DeleteContext: resourceF5vServerEndpointDelete,
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
				Description: "Name of F5 vServer Endpoint.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"f5_service_id": {
				Description: "ID of the F5 service associated with the" +
					"F5 vServer Endpoint.",
				Type:     schema.TypeString,
				Required: true,
			},
			"f5_service_instance_ids": {
				Description: "An array of F5 service instance IDs" +
					"when not provided an F5 vServer endpoint is " +
					"associated with all instances of the F5 Service",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"type": {
				Description: "The type of endpoint." +
					"Can be one of `ELB`, `BOTH`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ELB", "BOTH"}, false),
			},
			"segment_id": {
				Description: "ID of the segment associated with" +
					" the endpoint",
				Type:     schema.TypeString,
				Required: true,
			},
			"fqdn": {
				Description: "The FQDN of the endpoint",
				Type:        schema.TypeString,
				Required:    true,
			},
			"protocol": {
				Description: "The portocol used for the endpoint" +
					"Can be one of `TCP`, `UDP` or `ICMP`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
			},
			"port_ranges": {
				Description: "An array of ports or port ranges." +
					" Values can be mixed i.e. ['20', '100-200']." +
					" An array with only the value '-1' means any port.",
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"snat": {
				Description:  "SNAT for the endpoint.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTOMAP"}, false),
			},
		},
	}
}

func resourceF5vServerEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewF5vServerEndpoint(client)

	request, err := generateRequestF5vServerEndpoint(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	if client.Provision {
		d.Set("provision_state", provState)
		if provState == "FAILED" {
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "PROVISION (CREATE) FAILED",
					Detail:   fmt.Sprintf("%s", provErr),
				},
			}
		}
	}
	return resourceF5vServerEndpointRead(ctx, d, m)
}

func resourceF5vServerEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewF5vServerEndpoint(m.(*alkira.AlkiraClient))

	f5, provState, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}
	d.Set("name", f5.Name)
	d.Set("type", f5.Type)
	d.Set("fqdn", f5.Fqdn)
	d.Set("protocol", f5.Protocol)
	d.Set("port_ranges", f5.PortRanges)
	d.Set("snat", f5.Snat)
	d.Set("f5_service_id", f5.F5ServiceId)
	d.Set("f5_service_instance_ids", f5.F5ServiceInstanceIds)

	segmentId, err := getSegmentIdByName(f5.Segment, m)

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

func resourceF5vServerEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewF5vServerEndpoint(m.(*alkira.AlkiraClient))

	request, err := generateRequestF5vServerEndpoint(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	provState, err, provErr := api.Update(d.Id(), request)

	if client.Provision == true {
		d.Set("provision_state", provState)
		if provState == "FAILED" {
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "PROVISION (UPDATE) FAILED",
					Detail:   fmt.Sprintf("%s", provErr),
				},
			}
		}
	}
	return nil
}

func resourceF5vServerEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewF5vServerEndpoint(m.(*alkira.AlkiraClient))

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

func generateRequestF5vServerEndpoint(d *schema.ResourceData, m interface{}) (*alkira.F5vServerEndpoint, error) {

	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)
	if err != nil {
		return nil, err
	}
	request := &alkira.F5vServerEndpoint{
		Name:                 d.Get("name").(string),
		Segment:              segmentName,
		Type:                 d.Get("type").(string),
		F5ServiceId:          json.Number(d.Get("f5_service_id").(string)),
		F5ServiceInstanceIds: convertTypeSetToIntList(d.Get("f5_service_instance_ids").(*schema.Set)),
		Fqdn:                 d.Get("fqdn").(string),
		Protocol:             d.Get("protocol").(string),
		PortRanges:           convertTypeSetToStringList(d.Get("port_ranges").(*schema.Set)),
		Snat:                 d.Get("snat").(string),
	}

	return request, nil
}
