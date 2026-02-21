package alkira

import (
	"context"
	"fmt"

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
			StateContext: importWithReadValidation(resourceProbeHTTPSRead),
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
			"description": {
				Description: "The description of the probe.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"network_entity_id": {
				Description: "The ID of the internet application network entity to probe.",
				Type:        schema.TypeString,
				Required:    true,
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
					" `timeout_seconds` should always be less than or equal to `period_seconds`.",
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

	response, _, err, valErr, _ := api.Create(probe)
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation errors
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		// Try to read the resource to preserve any successfully created state
		readDiags := resourceProbeHTTPSRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
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
		return diag.Errorf("Probe type mismatch.")
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

	_, err, valErr, _ := api.Update(d.Id(), probe)
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation errors
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		// Try to read the resource to preserve current state
		readDiags := resourceProbeHTTPSRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (UPDATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	return resourceProbeHTTPSRead(ctx, d, m)
}

func resourceProbeHTTPSDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	_, err, valErr, _ := api.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation errors
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	d.SetId("")
	return nil
}
