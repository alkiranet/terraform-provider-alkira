package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorGcpVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Manage GCP Cloud Connector.",

		CreateContext: resourceConnectorGcpVpcCreate,
		ReadContext:   resourceConnectorGcpVpcRead,
		UpdateContext: resourceConnectorGcpVpcUpdate,
		DeleteContext: resourceConnectorGcpVpcDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceConnectorGcpVpcRead),
		},

		Schema: map[string]*schema.Schema{
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "ID of resource `credential_gcp_vpc`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"failover_cxps": {
				Description: "A list of additional CXPs where the connector " +
					"should be provisioned for failover.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"gcp_project_id": {
				Description: "GCP Project ID.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vpc_subnet": {
				Description: "The list of subnets of the target GCP VPC for " +
					"routing purpose. Given connector supports multiple prefixes " +
					"per subnet, each prefix under a subnet will be a new entry.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "An identifier for the subnetwork " +
								"resource with format " +
								"`projects/{{project}}/regions/{{region}}/subnetworks/{{name}}`.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"cidr": {
							Description: "The CIDR of the subnet.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"gcp_routing": {
				Description: "GCP Routing describes the routes that are to be " +
					"imported to the VPC from the CXP. This essentially controls " +
					"how traffic is routed between the CXP and the VPC. " +
					"When routing option is not provided, the traffic exiting " +
					"the VPC will be sent to the CXP (i.e a default route to " +
					"CXP will be added to all route tables on that VPC)",
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_list_ids": {
							Description: "IDs of prefix lists defined on the " +
								"network.",
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"custom_prefix": {
							Description: "Specifies the source of the routes " +
								"that need to be imported. The value could be " +
								"`ADVERTISE_DEFAULT_ROUTE` and " +
								"`ADVERTISE_CUSTOM_PREFIX`.",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ADVERTISE_DEFAULT_ROUTE",
								"ADVERTISE_CUSTOM_PREFIX",
							}, false),
						},
					},
				},
				Optional: true,
			},
			"gcp_region": {
				Description: "GCP region where VPC resides.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"gcp_vpc_name": {
				Description: "GCP VPC name.",
				Type:        schema.TypeString,
				Required:    true,
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
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_id": {
				Description: "The ID of the segment associated with the " +
					"connector.",
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Description: "The size of the connector, one of `5XSMALL`,`XSMALL`,`SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"customer_asn": {
				Description: "A specific BGP ASN for the connector. This " +
					"field cannot be updated once the connector has been " +
					"provisioned. The ASN can be any private ASN (`64512 " +
					"- 65534`, `4200000000 - 4294967294`) that is not used " +
					"elsewhere in the network.",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  64522,
			},
			"scale_group_id": {
				Description: "The ID of the scale group associated with the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceConnectorGcpVpcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	request, err := generateConnectorGcpVpcRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set states
	d.SetId(string(response.Id))

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorGcpVpcRead(ctx, d, m)
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

	if client.Provision {
		d.Set("provision_state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorGcpVpcRead(ctx, d, m)
}

func resourceConnectorGcpVpcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	// GET
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("credential_id", connector.CredentialId)
	d.Set("gcp_region", connector.CustomerRegion)
	d.Set("enabled", connector.Enabled)
	d.Set("failover_cxps", connector.SecondaryCXPs)
	d.Set("gcp_project_id", connector.ProjectId)
	d.Set("gcp_vpc_name", connector.VpcName)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("customer_asn", connector.CustomerASN)
	d.Set("scale_group_id", connector.ScaleGroupId)
	d.Set("description", connector.Description)
	setGcpRoutingOptions(connector.GcpRouting, d)
	setGcpVpcSubnets(connector.GcpRouting, d)

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
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorGcpVpcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	request, err := generateConnectorGcpVpcRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorGcpVpcRead(ctx, d, m)
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
	if client.Provision {
		d.Set("provision_state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorGcpVpcRead(ctx, d, m)
}

func resourceConnectorGcpVpcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	// DELETE
	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_connector_gcp_vpc (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_connector_gcp_vpc (id=%s)", err, d.Id()))
	}

	d.SetId("")

	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	// Check provision state
	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
