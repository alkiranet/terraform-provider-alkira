package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraProbeHTTPS() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Alkira HTTPS Probe.",
		CreateContext: resourceProbeHTTPSCreate,
		ReadContext:   resourceProbeHTTPSRead,
		UpdateContext: resourceProbeHTTPSUpdate,
		DeleteContext: resourceProbeHTTPSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the HTTPS probe.",
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
			"server_name": {
				Description: "The server name for TLS SNI.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"disable_cert_validation": {
				Description: "Whether to disable certificate validation.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"ca_certificate": {
				Description: "Required when certificate validation" +
					" is enabled and certificate is self assigned.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"headers": {
				Description: "HTTP headers to include in the request.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
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
							Type:          schema.TypeString,
							Description:   "Response code the response should have.",
							Optional:      true,
							ConflictsWith: []string{"regex"},
						},
						"regex": {
							Description:   "Regex the response body should match.",
							Type:          schema.TypeString,
							ConflictsWith: []string{"status_code"},
							Optional:      true,
						},
					},
				},
			},
			"failure_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of consecutive failures required to mark the probe as failed.",
			},
			"success_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of consecutive successes required to mark the probe as successful.",
			},
			"period_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "How often (in seconds) to perform the probe.",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of seconds after which the probe times out.",
			},
		},
	}
}

func resourceProbeHTTPSCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, err := generateHTTPSProbeRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, _, err, _ := api.Create(probe)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.ID)

	return resourceProbeHTTPSRead(ctx, d, m)
}

func resourceProbeHTTPSRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, _, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Failed to Get HTTPS Probe",
			Detail:   err.Error(),
		}}
	}

	if probe.Type != "HTTPS" {
		return diag.Errorf("Retrieved probe is not of HTTPS type")
	}

	if err := setHTTPSProbeState(probe, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceProbeHTTPSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, err := generateHTTPSProbeRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err, _ = api.Update(d.Id(), probe)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProbeHTTPSRead(ctx, d, m)
}

func resourceProbeHTTPSDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	_, err, _ := api.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
