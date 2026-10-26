package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorMultiCloudGateway() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a Multi-Cloud Gateway connector.",
		CreateContext: resourceConnectorMultiCloudGatewayCreate,
		ReadContext:   resourceConnectorMultiCloudGatewayRead,
		UpdateContext: warnOnFailedStateUpdate(resourceConnectorMultiCloudGatewayUpdate),
		DeleteContext: resourceConnectorMultiCloudGatewayDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return validateMultiCloudGatewayOnRampRegions(d.Get("instance").([]interface{}))
		},
		Importer: &schema.ResourceImporter{StateContext: importWithReadValidation(resourceConnectorMultiCloudGatewayRead)},
		Schema: map[string]*schema.Schema{
			"name":        {Type: schema.TypeString, Required: true, Description: "Name of the connector."},
			"description": {Type: schema.TypeString, Optional: true, Description: "Description of the connector."},
			"cxp": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CXP where the connector is provisioned. Only AWS CXPs support a Multi-Cloud Gateway connector.",
			},
			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "2LARGE", "5LARGE", "10LARGE"}, false),
				Description:  "Size of each instance: `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `5LARGE`, or `10LARGE`.",
			},
			"group":   {Type: schema.TypeString, Optional: true, Description: "Group of the connector."},
			"enabled": {Type: schema.TypeBool, Optional: true, Default: true, Description: "Whether the connector is enabled. Default value is `true`."},
			"provision_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The provision state of the connector.",
			},
			"implicit_group_id": {
				Description: "The ID of the implicit group associated " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"multi_cloud_gateway_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the `alkira_multi_cloud_gateway` this connector attaches to. The gateway must be `PROVISIONED`.",
			},
			"loopback_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(26, 26),
				Description:  "A free /26 of the gateway's loopback pool, used for tunnel endpoints.",
			},
			"billing_tag_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Billing tags associated with the connector. (see resource `alkira_billing_tag`).",
			},
			"segment_options": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Segment options for the connector. Exactly one block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "The ID of the segment.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"advertise_on_prem_routes": {
							Description: "Advertise on-prem routes. Default value is `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"disable_internet_exit": {
							Description: "Disable internet exit for this segment. Default value is `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
			"instance": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Instances of the connector. On update, a block is matched to its stored instance by `name`; a renamed block is treated as a new instance and the old one is removed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the instance, unique in the connector.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"on_ramp_region": {
							Description: "The Lumen on-ramp region of this instance's hosted connection. Can only be changed while the instance has not yet been provisioned.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "ID Alkira assigned to the instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"bandwidth_in_mbps": {
							Description: "Bandwidth of the instance when it was created. Never changes after.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"lumen_connection_id": {
							Description: "Lumen Connect ID to correlate the instance.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceConnectorMultiCloudGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorMultiCloudGateway(client)
	request, err := generateConnectorMultiCloudGatewayRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	connector, provState, err, valErr, provErr := api.Create(request)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(connector.Id.String())
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorMultiCloudGatewayRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorMultiCloudGatewayRead(ctx, d, m)
}

func resourceConnectorMultiCloudGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	connector, provState, err := alkira.NewConnectorMultiCloudGateway(client).GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	if diags := setConnectorMultiCloudGateway(d, connector, m); diags.HasError() {
		return diags
	}

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorMultiCloudGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	request, err := generateConnectorMultiCloudGatewayRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	provState, err, valErr, provErr := alkira.NewConnectorMultiCloudGateway(client).Update(d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorMultiCloudGatewayRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (UPDATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorMultiCloudGatewayRead(ctx, d, m)
}

func resourceConnectorMultiCloudGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	_, err, valErr, provErr := alkira.NewConnectorMultiCloudGateway(client).Delete(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("deleting multi-cloud gateway connector %s: %w", d.Id(), err))
	}
	d.SetId("")

	if valErr != nil {
		return diag.FromErr(valErr)
	}

	if client.Provision && provErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
