package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraProbeTCP() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Alkira TCP Probe.",
		CreateContext: resourceProbeTCPCreate,
		ReadContext:   resourceProbeTCPRead,
		UpdateContext: resourceProbeTCPUpdate,
		DeleteContext: resourceProbeTCPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the TCP probe.",
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
			"port": {
				Description: "The TCP port to probe.",
				Type:        schema.TypeInt,
				Required:    true,
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

func resourceProbeTCPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, err := generateTCPProbeRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, _, err, _ := api.Create(probe)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.ID)

	return resourceProbeTCPRead(ctx, d, m)
}

func resourceProbeTCPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, _, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Failed to Get TCP Probe",
			Detail:   err.Error(),
		}}
	}

	if probe.Type != "TCP" {
		return diag.Errorf("Probe type mismatch.")
	}

	if err := setTCPProbeState(probe, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceProbeTCPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	probe, err := generateTCPProbeRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err, _ = api.Update(d.Id(), probe)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProbeTCPRead(ctx, d, m)
}

func resourceProbeTCPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewProbe(client)

	_, err, _ := api.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
