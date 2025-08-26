package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAwsVpc() *schema.Resource {
	return &schema.Resource{
		Description:   "Provide AWS VPC Connector resource.",
		CreateContext: resourceConnectorAwsVpcCreate,
		ReadContext:   resourceConnectorAwsVpcRead,
		UpdateContext: resourceConnectorAwsVpcUpdate,
		DeleteContext: resourceConnectorAwsVpcDelete,
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
			"aws_account_id": {
				Description: "AWS Account ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"aws_region": {
				Description: "AWS Region where VPC resides.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "ID of resource `credential_aws_vpc`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "The CXP where the connector should be " +
					"provisioned.",
				Type:     schema.TypeString,
				Required: true,
			},
			"direct_inter_vpc_communication_enabled": {
				Description: "Enable direct inter-vpc communication. " +
					"Default is set to `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"direct_inter_vpc_communication_group": {
				Description: "Direct inter-vpc communication group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Whether the connector is enabled. Default is " +
					"`true`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"failover_cxps": {
				Description: "A list of additional CXPs where the connector " +
					"should be provisioned for failover.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"provision_state": {
				Description: "The provisioning state of connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_id": {
				Description: "The ID of segments associated with the connector. " +
					"Currently, only `1` segment is allowed.",
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Description: "The size of the connector, one of `5XSMALL`,`XSMALL`,`SMALL`, `MEDIUM`, " +
					"`LARGE`, `2LARGE`, `5LARGE`, `10LARGE`, `20LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"tgw_connect_enabled": {
				Description: "When it's set to `true`, Alkira will use TGW Connect " +
					"attachments to build connection to AWS Transit Gateway. " +
					"Connect Attachments suppport GRE tunnel protocol for high " +
					"performance and BGP for dynamic routing. This applies to " +
					"all TGW attachments. This field can be set to `true` only " +
					"if the VPC is in the same AWS region as the Alkira CXP " +
					"it is being deployed onto.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tgw_attachment": {
				Description: "TGW attachment.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Description: "The Id of the subnet.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"az": {
							Description: "The availability zone of the subnet.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"vpc_id": {
				Description: "The ID of the target VPC.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vpc_cidr": {
				Description: "The list of CIDR attached to the target VPC for " +
					"routing purpose. It could be only specified if " +
					"`vpc_subnet` is not specified.",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"vpc_subnet"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"vpc_subnet": {
				Description: "The list of subnets of the target VPC for " +
					"routing purpose. It could only specified if `vpc_cidr` " +
					"is not specified.",
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"vpc_cidr"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The Id of the subnet.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"cidr": {
							Description: "The CIDR of the subnet.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"vpc_route_table": {
				Description: "VPC route table",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The Id of the route table",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"prefix_list_ids": {
							Description: "Prefix List IDs",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"options": {
							Description: "Routing options, one of " +
								"`ADVERTISE_DEFAULT_ROUTE`, " +
								"`OVERRIDE_DEFAULT_ROUTE` or " +
								"`ADVERTISE_CUSTOM_PREFIX`.",
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ADVERTISE_DEFAULT_ROUTE",
								"OVERRIDE_DEFAULT_ROUTE",
								"ADVERTISE_CUSTOM_PREFIX"}, false),
						},
					},
				},
				Optional: true,
			},
			"scale_group_id": {
				Description: "The ID of the scale group associated with the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"overlay_subnets": {
				Description: "Overlay subnet.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceConnectorAwsVpcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsVpc(client)

	request, err := generateConnectorAwsVpcRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set the state
	d.SetId(string(response.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAwsVpcRead(ctx, d, m)
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

	return resourceConnectorAwsVpcRead(ctx, d, m)
}

func resourceConnectorAwsVpcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsVpc(client)

	// GET
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("aws_account_id", connector.VpcOwnerId)
	d.Set("aws_region", connector.CustomerRegion)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("credential_id", connector.CredentialId)
	d.Set("cxp", connector.CXP)
	d.Set("direct_inter_vpc_communication_enabled", connector.DirectInterVPCCommunicationEnabled)
	d.Set("direct_inter_vpc_communication_group", connector.DirectInterVPCCommunicationGroup)
	d.Set("enabled", connector.Enabled)
	d.Set("failover_cxps", connector.SecondaryCXPs)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("vpc_id", connector.VpcId)
	d.Set("tgw_connect_enabled", connector.TgwConnectEnabled)
	d.Set("scale_group_id", connector.ScaleGroupId)
	d.Set("description", connector.Description)

	// Get segment
	numOfSegments := len(connector.Segments)

	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(connector.Segments[0], m)

		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("segment_id", segmentId)
	} else {
		return diag.FromErr(fmt.Errorf("failed to find segment"))
	}

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorAwsVpcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsVpc(client)

	request, err := generateConnectorAwsVpcRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAwsVpcRead(ctx, d, m)
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

	return resourceConnectorAwsVpcRead(ctx, d, m)
}

func resourceConnectorAwsVpcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsVpc(client)

	// DELETE
	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	// Handle validation error
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
