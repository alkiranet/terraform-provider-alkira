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
				Description: "Tags for billing.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "ID of credential managed by Credential Manager.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"direct_inter_vpc_communication": {
				Description: "Enable direct inter-vpc communication. Default is set to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
				Type:     schema.TypeList,
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
				Description: "The size of the connector, one of `SMALL`, `MEDIUM`, " +
					"`LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, `10LARGE`, `20LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "2LARGE", "4LARGE", "5LARGE", "10LARGE", "20LARGE"}, false),
			},
			"tgw_connect_enabled": {
				Description: "On setting tgwConnectEnabled to true Alkira will use TGW Connect attachments to build connection to AWS Transit gateway." +
            		"Connect Attachments suppport Generic Routing Encapsulation (GRE) tunnel protocol for high performance," +
            		"and Border Gateway Protocol (BGP) for dynamic routing. This applies to all tgw attachments." + 
            		"This field can be set to true, only if the vpc is in the same AWS region as the Alkira CXP it is being onboarded onto.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"tgw_attachment": {
				Description: "TGW attachment.",
				Type:        schema.TypeSet,
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
				Description: "The list of CIDR attached to the target VPC for routing " +
					"purpose. It could be only specified if `vpc_subnet` is not specified.",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"vpc_subnet"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"vpc_subnet": {
				Description: "The list of subnets of the target VPC for routing purpose. " +
					"It could only specified if `vpc_cidr` is not specified.",
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
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"options": {
							Description: "Routing options, one of `ADVERTISE_DEFAULT_ROUTE`, " +
								"`OVERRIDE_DEFAULT_ROUTE` and `ADVERTISE_CUSTOM_PREFIX`.",
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
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set the state
	d.SetId(string(response.Id))

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
	d.Set("direct_inter_vpc_communication", connector.DirectInterVPCCommunicationEnabled)
	d.Set("enabled", connector.Enabled)
	d.Set("failover_cxps", connector.SecondaryCXPs)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("vpc_id", connector.VpcId)
	d.Set("tgw_connect_enabled", connector.TgwConnectEnabled)

	// Get segment
	numOfSegments := len(connector.Segments)

	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(connector.Segments[0], m)

		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("segment_id", segmentId)
	} else {
		return diag.FromErr(fmt.Errorf("the number of segments are invalid %n", numOfSegments))
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

	return resourceConnectorAwsVpcRead(ctx, d, m)
}

func resourceConnectorAwsVpcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsVpc(client)

	// DELETE
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

// generateConnectorAwsVpcRequest generate request for connector-aws-vpc
func generateConnectorAwsVpcRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAwsVpc, error) {

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	inputPrefixes, err := generateUserInputPrefixes(d.Get("vpc_cidr").([]interface{}), d.Get("vpc_subnet").(*schema.Set))

	if err != nil {
		return nil, err
	}

	exportOptions := alkira.ExportOptions{
		Mode:     "USER_INPUT_PREFIXES",
		Prefixes: inputPrefixes,
	}

	routeTables := expandAwsVpcRouteTables(d.Get("vpc_route_table").(*schema.Set))
	tgwAttachments := expandAwsVpcTgwAttachments(d.Get("tgw_attachment").(*schema.Set))

	vpcRouting := alkira.ConnectorAwsVpcRouting{
		Export: exportOptions,
		Import: alkira.ImportOptions{routeTables},
	}

	request := &alkira.ConnectorAwsVpc{
		BillingTags:                        convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		CXP:                                d.Get("cxp").(string),
		CredentialId:                       d.Get("credential_id").(string),
		CustomerName:                       m.(*alkira.AlkiraClient).Username,
		CustomerRegion:                     d.Get("aws_region").(string),
		DirectInterVPCCommunicationEnabled: d.Get("direct_inter_vpc_communication").(bool),
		Enabled:                            d.Get("enabled").(bool),
		Group:                              d.Get("group").(string),
		Name:                               d.Get("name").(string),
		Segments:                           []string{segmentName},
		SecondaryCXPs:                      convertTypeListToStringList(d.Get("failover_cxps").([]interface{})),
		Size:                               d.Get("size").(string),
		TgwConnectEnabled:                  d.Get("tgw_connect_enabled").(bool),
		TgwAttachments:                     tgwAttachments,
		VpcId:                              d.Get("vpc_id").(string),
		VpcOwnerId:                         d.Get("aws_account_id").(string),
		VpcRouting:                         vpcRouting,
	}

	return request, nil
}
