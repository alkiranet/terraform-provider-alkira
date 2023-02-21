package alkira

import (
	"context"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorGcpVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Manage GCP Cloud Connector.",

		Create: resourceConnectorGcpVpcCreate,
		Read:   resourceConnectorGcpVpcRead,
		Update: resourceConnectorGcpVpcUpdate,
		Delete: resourceConnectorGcpVpcDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"failover_cxps": {
				Description: "A list of additional CXPs where the connector should be provisioned for failover.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"gcp_project_id": {
				Description: "GCP Project ID.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vpc_subnet": {
				Description: "The list of subnets of the target GCP VPC for routing purpose. " +
					"Given GCP VPC supports multiple prefixes per subnet, each prefix under a subnet will be a new entry.",
				Type:     schema.TypeSet,
				Optional: true,
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
			"gcp_routing": {
				Description: "GCP Routing describes the routes that are to be imported to the VPC " +
					"from the CXP. This essentially controls how traffic is routed between the " +
					"CXP and the VPC. gcpRouting provides a customized routing specification. " +
					"When gcpRouting is not provided i.e when gcpRouting is null/empty then all " +
					"traffic exiting the VPC will be sent to the CXP (i.e a default route to " +
					"CXP will be added to all route tables on that VPC)",
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_list_ids": {
							Description: "Ids of prefix lists defined on the network.",
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"custom_prefix": {
							Description: "custom_prefix is an instruction which specifies " +
								"the source of the routes that need to be imported. Only " +
								"`ADVERTISE_DEFAULT_ROUTE` and `ADVERTISE_CUSTOM_PREFIX` are valid inputs.",
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
			"gcp_vpc_id": {
				Description: "GCP VPC ID.",
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
				Description: "The ID of implicit group automaticaly created with the connector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_id": {
				Description: "The ID of the segment associated with the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM` or `LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, `10LARGE`, `20LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "2LARGE", "4LARGE", "5LARGE", "10LARGE", "20LARGE"}, false),
			},
		},
	}
}

func resourceConnectorGcpVpcCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorGcpVpcRequest(d, m)

	if err != nil {
		return err
	}

	// Send create request
	response, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	// Set states
	d.SetId(string(response.Id))

	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceConnectorGcpVpcRead(d, m)
}

func resourceConnectorGcpVpcRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	connector, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("credential_id", connector.CredentialId)
	d.Set("customer_region", connector.CustomerRegion)
	d.Set("enabled", connector.Enabled)
	d.Set("failover_cxps", connector.SecondaryCXPs)
	d.Set("gcp_project_id", connector.ProjectId)
	d.Set("gcp_vpc_id", connector.VpcId)
	d.Set("gcp_vpc_name", connector.VpcName)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	setGcpRoutingOptions(connector.GcpRouting, d)

	// Get segment
	numOfSegments := len(connector.Segments)
	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(connector.Segments[0], m)

		if err != nil {
			return err
		}
		d.Set("segment_id", segmentId)
	} else {
		return fmt.Errorf("the number of segments are invalid %n", numOfSegments)
	}

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return err
}

func resourceConnectorGcpVpcUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorGcpVpcRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceConnectorGcpVpcRead(d, m)
}

func resourceConnectorGcpVpcDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete connector_gcp_vpc %s, provision failed", d.Id())
	}

	d.SetId("")

	return nil
}

func generateConnectorGcpVpcRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorGcpVpc, error) {

	//
	// Routing
	//
	gcpRouting, err := convertGcpRouting(d.Get("gcp_routing").(*schema.Set), d.Get("vpc_subnet").(*schema.Set))
	if err != nil {
		log.Printf("[ERROR] failed to convert gcp routing")
		return nil, err
	}

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	// Assemble request
	connector := &alkira.ConnectorGcpVpc{
		BillingTags:    convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		GcpRouting:     gcpRouting,
		CustomerRegion: d.Get("gcp_region").(string),
		Enabled:        d.Get("enabled").(bool),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		ProjectId:      d.Get("gcp_project_id").(string),
		Segments:       []string{segmentName},
		SecondaryCXPs:  convertTypeListToStringList(d.Get("failover_cxps").([]interface{})),
		Size:           d.Get("size").(string),
		VpcId:          d.Get("gcp_vpc_id").(string),
		VpcName:        d.Get("gcp_vpc_name").(string),
	}

	return connector, nil
}
