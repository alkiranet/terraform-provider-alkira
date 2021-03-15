package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraConnectorAwsVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Manage AWS Cloud Connector.",
		Create:      resourceConnectorAwsVpcCreate,
		Read:        resourceConnectorAwsVpcRead,
		Update:      resourceConnectorAwsVpcUpdate,
		Delete:      resourceConnectorAwsVpcDelete,

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
			"billing_tags": {
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
			"connector_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment": {
				Description: "The segment of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE` or `4LARGE`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vpc_id": {
				Description: "The ID of the VPC the connnector connects to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vpc_cidr": {
				Description: "The CIDR of the VPC the connnector connects to.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{"vpc_subnets"},
			},
			"vpc_subnets": {
				Description: "The subnets of the VPC the connnector connects to.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"vpc_cidr"},
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
							Description: "Prefix List Ids",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},

						"options": {
							Description: "Routing options, one of `ADVERTISE_DEFAULT_ROUTE`, `OVERRIDE_DEFAULT_ROUTE` and `ADVERTISE_CUSTOM_PREFIX`.",
							Type:        schema.TypeString,
							Optional:    true,
							ValidateFunc: validation.StringInSlice([]string{"ADVERTISE_DEFAULT_ROUTE", "OVERRIDE_DEFAULT_ROUTE", "ADVERTISE_CUSTOM_PREFIX"}, false),
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceConnectorAwsVpcCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tags").([]interface{}))

	segments := []string{d.Get("segment").(string)}

	inputPrefixes, err := generateUserInputPrefixes(d.Get("vpc_cidr").(string), d.Get("vpc_subnets").([]interface{}))

	if err != nil {
		return err
	}

	exportOptions := alkira.ExportOptions{
			Mode:     "USER_INPUT_PREFIXES",
			Prefixes: inputPrefixes,
	}

	routeTables := expandAwsVpcRouteTables(d.Get("vpc_route_table").(*schema.Set))

	vpcRouting := alkira.ConnectorAwsVpcRouting{
		Export: exportOptions,
		Import: alkira.ImportOptions{routeTables},
	}

	connectorAwsVpc := &alkira.ConnectorAwsVpcRequest{
		BillingTags:    billingTags,
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		CustomerName:   client.Username,
		CustomerRegion: d.Get("aws_region").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		Segments:       segments,
		Size:           d.Get("size").(string),
		VpcId:          d.Get("vpc_id").(string),
		VpcOwnerId:     d.Get("aws_account_id").(string),
		VpcRouting:     vpcRouting,
	}

	id, err := client.CreateConnectorAwsVpc(connectorAwsVpc)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("connector_id", id)
	return resourceConnectorAwsVpcRead(d, m)
}

func resourceConnectorAwsVpcRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorAwsVpcUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceConnectorAwsVpcRead(d, m)
}

func resourceConnectorAwsVpcDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (AWS-VPC) %s", d.Id())
	err := client.DeleteConnectorAwsVpc(d.Get("connector_id").(int))

	if err != nil {
		return err
	}

	return nil
}
