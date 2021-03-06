package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Description: "The segment of the connector belongs to. Currently, only `1` segment is allowed.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "2LARGE"}, false),
			},
			"vpc_id": {
				Description: "The ID of the target VPC.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vpc_cidr": {
				Description:   "The list of CIDR attached to the target VPC for routing purpose. It could be only specified if `vpc_subnet` is not specified.",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"vpc_subnet"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"vpc_subnet": {
				Description:   "The list of subnets of the target VPC for routing purpose. It could only specified if `vpc_cidr` is not specified.",
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
							Description: "Prefix List Ids",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"options": {
							Description:  "Routing options, one of `ADVERTISE_DEFAULT_ROUTE`, `OVERRIDE_DEFAULT_ROUTE` and `ADVERTISE_CUSTOM_PREFIX`.",
							Type:         schema.TypeString,
							Optional:     true,
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

	connector, err := generateConnectorAwsVpcRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateConnectorAwsVpc(connector)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConnectorAwsVpcRead(d, m)
}

func resourceConnectorAwsVpcRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorAwsVpcUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorAwsVpcRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing Connector (AWS-VPC) %s", d.Id())
	err = client.UpdateConnectorAwsVpc(d.Id(), connector)

	return err
}

func resourceConnectorAwsVpcDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (AWS-VPC) %s", d.Id())
	err := client.DeleteConnectorAwsVpc(d.Id())

	if err != nil {
		return err
	}

	return nil
}

// generateConnectorAwsVpcRequest generate request for connector-aws-vpc
func generateConnectorAwsVpcRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAwsVpcRequest, error) {
	client := m.(*alkira.AlkiraClient)
	billingTags := convertTypeListToIntList(d.Get("billing_tags").([]interface{}))
	segments := []string{d.Get("segment").(string)}
	inputPrefixes, err := generateUserInputPrefixes(d.Get("vpc_cidr").([]interface{}), d.Get("vpc_subnet").(*schema.Set))

	if err != nil {
		return nil, err
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

	request := &alkira.ConnectorAwsVpcRequest{
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

	return request, nil
}
