package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorOciVcn() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Oracle Cloud (OCI) Virtual Computing Network (VCN) Cloud Connector.",
		Create:      resourceConnectorOciVcnCreate,
		Read:        resourceConnectorOciVcnRead,
		Update:      resourceConnectorOciVcnUpdate,
		Delete:      resourceConnectorOciVcnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"oci_region": {
				Description: "OCI region of the VCN.",
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
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"credential_id": {
				Description: "ID of OCI-VCN credential.",
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
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created with " +
					"the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_id": {
				Description: "The ID of segments associated with the connector. " +
					"Currently, only `1` segment is allowed.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, " +
					"`10LARGE`, `20LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "2LARGE", "4LARGE", "5LARGE", "10LARGE", "20LARGE"}, false),
			},
			"vcn_id": {
				Description: "The OCID of the VCN.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vcn_cidr": {
				Description: "The list of CIDR attached to the target VCN " +
					"for routing purpose. It could be only specified if " +
					"`vcn_subnet` is not specified.",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"vcn_subnet"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"vcn_subnet": {
				Description: "The list of subnets of the target VCN for " +
					"routing purpose. It could only specified if `vcn_cidr` " +
					"is not specified.",
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"vcn_cidr"},
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
			"vcn_route_table": {
				Description: "VCN route table.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the route table.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"prefix_list_ids": {
							Description: "Prefix List IDs.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"options": {
							Description: "Routing options, one of `ADVERTISE_DEFAULT_ROUTE`, " +
								"`OVERRIDE_DEFAULT_ROUTE` and `ADVERTISE_CUSTOM_PREFIX`.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ADVERTISE_DEFAULT_ROUTE", "OVERRIDE_DEFAULT_ROUTE", "ADVERTISE_CUSTOM_PREFIX"}, false),
						},
					},
				},
				Optional: true,
			},
			"billing_tag_ids": {
				Description: "IDs of billing tags.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceConnectorOciVcnCreate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorOciVcn(m.(*alkira.AlkiraClient))

	// Construct request
	connector, err := generateConnectorOciVcnRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	resource, provisionState, err := api.Create(connector)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)
	d.SetId(string(resource.Id))

	return resourceConnectorOciVcnRead(d, m)
}

func resourceConnectorOciVcnRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorOciVcn(m.(*alkira.AlkiraClient))

	// Read connector
	connector, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("credential_id", connector.CredentialId)
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("failover_cxps", connector.SecondaryCXPs)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("oci_region", connector.CustomerRegion)
	d.Set("size", connector.Size)
	d.Set("vcn_id", connector.VcnId)

	if len(connector.Segments) > 0 {
		segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
		segment, _, err := segmentApi.GetByName(connector.Segments[0])

		if err != nil {
			return err
		}
		d.Set("segment_id", segment.Id)
	}

	return nil
}

func resourceConnectorOciVcnUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorOciVcn(m.(*alkira.AlkiraClient))

	// Construct request
	connector, err := generateConnectorOciVcnRequest(d, m)

	if err != nil {
		return err
	}

	// Send request to update connector
	provisionState, err := api.Update(d.Id(), connector)

	d.Set("provision_state", provisionState)
	return err
}

func resourceConnectorOciVcnDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorOciVcn(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if provisionState != "SUCCESS" {
	}

	d.SetId("")
	return nil
}

// generateConnectorOciVcnRequest generate request for connector-oci-vcn
func generateConnectorOciVcnRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorOciVcn, error) {

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	failoverCXPs := convertTypeListToStringList(d.Get("failover_cxps").([]interface{}))

	//
	// Construct Segment
	//
	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, err := segmentApi.GetById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by ID: %d", d.Get("segment_id"))
		return nil, err
	}

	//
	// Construct Routing Options
	//
	inputPrefixes, err := generateConnectorOciVcnUserInputPrefixes(d.Get("vcn_cidr").([]interface{}), d.Get("vcn_subnet").(*schema.Set))

	if err != nil {
		return nil, err
	}

	exportOptions := alkira.ConnectorOciVcnExportOptions{
		Mode:     "USER_INPUT_PREFIXES",
		Prefixes: inputPrefixes,
	}

	routeTables := expandConnectorOciVcnRouteTables(d.Get("vcn_route_table").(*schema.Set))

	vcnRouting := alkira.ConnectorOciVcnRouting{
		Export: exportOptions,
		Import: alkira.ConnectorOciVcnImportOptions{routeTables},
	}

	//
	// Construct request
	//
	request := &alkira.ConnectorOciVcn{
		BillingTags:    billingTags,
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		CustomerRegion: d.Get("oci_region").(string),
		Enabled:        d.Get("enabled").(bool),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		SecondaryCXPs:  failoverCXPs,
		Segments:       []string{segment.Name},
		Size:           d.Get("size").(string),
		VcnId:          d.Get("vcn_id").(string),
		VcnRouting:     vcnRouting,
	}

	return request, nil
}
