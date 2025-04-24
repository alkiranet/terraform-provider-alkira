package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraProbeHTTP() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Alkira HTTP Probe.",
		CreateContext: resourceProbeHTTPCreate,
		ReadContext:   resourceProbeHTTPRead,
		UpdateContext: resourceProbeHTTPUpdate,
		DeleteContext: resourceProbeHTTPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the HTTP probe.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Whether the probe is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"network_entity": {
				Description: "Network entity configuration.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The type of network entity to probe." +
								" Only `INTERNET_APPLICATION` supported for now.",
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Description: "The ID of the network entity.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"uri": {
				Description: "The URI to probe.",
				Type:        schema.TypeString,
				Required:    true,
			},
			// will be enabled in phase 2.
			// "headers": {
			// 	Description: "HTTP headers to include in the request.",
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Elem:        &schema.Schema{Type: schema.TypeString},
			// },
			"validators": {
				Description: "Validators for the HTTP response.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Type of Validator, can be one of " +
								"`HTTP_STATUS_CODE` or `HTTP_RESPONSE_BODY`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP_STATUS_CODE", "HTTP_RESPONSE_BODY"}, false),
						},
						"status_code": {
							Type:        schema.TypeString,
							Description: "Response code the response should have.",
							Optional:    true,
						},
						"regex": {
							Description: "Regex the response body should match.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"failure_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
				Description: "The number of consecutive failures required to mark the probe as failed." +
					" Default is `3`, and the maximum value allowed is `50`.",
			},
			"success_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				Description: "The number of consecutive successes required to mark the probe as successful." +
					" Default value is `1`, and the maximum value allowed is `50`.",
			},
			"period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
				Description: "How often (in seconds) to perform the probe." +
					" Default value is `60`, and the maximum value allowed is `360`.",
			},
			"timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
				Description: "Number of seconds after which the probe times out." +
					" Default value is `60`, and the maximum value allowed is `360`." +
					" `timeout_seconds` should always be greater than `period_seconds`.",
			},
		},
	}
}

func resourceProbeHTTPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, err := generateHTTPProbeRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, _, err, _ := api.Create(probe)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.ID)

	return resourceProbeHTTPRead(ctx, d, m)
}

func resourceProbeHTTPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, _, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Failed to Get HTTP Probe",
			Detail:   err.Error(),
		}}
	}

	if probe.Type != "HTTP" {
		return diag.Errorf("Probe type mismatch.")
	}

	if err := setHTTPProbeState(probe, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceProbeHTTPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, err := generateHTTPProbeRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err, _ = api.Update(d.Id(), probe)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProbeHTTPRead(ctx, d, m)
}

func resourceProbeHTTPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	_, err, _ := api.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
