package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraInternetApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Internet Application.",
		Create:      resourceInternetApplicationCreate,
		Read:        resourceInternetApplicationRead,
		Update:      resourceInternetApplicationUpdate,
		Delete:      resourceInternetApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"billing_tag_ids": {
				Description: "IDs of billing tags.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"byoip_id": {
				Description: "BYOIP ID.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"connector_id": {
				Description: "Connector ID.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"connector_type": {
				Description: "Connector Type.The value could be `AWS_VPC`, " +
					"`AZURE_VNET`, `GCP_VPC`, `OCI_VCN`, `SD_WAN`, `IP_SEC` " +
					"`ARUBA_EDGE_CONNECT`, `EXPRESS_ROUTE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"fqdn_prefix": {
				Description: "User provided FQDN prefix that will be published on AWS Route 53.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group_id": {
				Description: "ID of the auto generated system group.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"inbound_connector_id": {
				Description: "Inbound connector ID. When `inbound_connector_type` is `DEFAULT`, " +
					"it could be left empty.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"inbound_connector_type": {
				Description: "The inbound connector type specifies how the internet application " +
					"is to be opened up to the external world. By `DEFAULT` the native cloud " +
					"internet connector is used. In this scenario, Alkira takes care of creating " +
					"this inbound internet connector implicitly. If instead inbound access is via " +
					"the `AKAMAI_PROLEXIC` connector, then you need to create and configure " +
					"that connector and use it with the internet application.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "AKAMAI_PROLEXIC"}, false),
			},
			"internet_protocol": {
				Description: "Internet Protocol to be associated with the internet application. " +
					"The value could be: `IPV4`, `IPV6` or `BOTH`. " +
					"In order to use the option IPV6 or BOTH, `enable_ipv6_to_ipv4_translation` " +
					"should be enabled on the associated segment and a valid IP pool range should " +
					"be provided. `IPV6` and `BOTH` options are only available to Internet " +
					"Applications on AWS CXPs. (**BETA**)",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPV4",
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6", "BOTH"}, false),
			},
			"name": {
				Description: "The name of the internet application.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the internet application.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public_ips": {
				Description: "This option pertains to the `AKAMAI_PROLEXIC` inbound_connector_type. " +
					"The public IPs are to be used to access the internet application. These public IPs " +
					"must belong to one of the BYOIP ranges configured for the Akamai Prolexic Connector.",
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"segment_id": {
				Description: "The ID of segment associated with the internet application.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description:  "The size of the internet application, one of `SMALL`, `MEDIUM` and `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"target": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "The type of the target, one of `IP` or `ILB_NAME`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"IP", "ILB_NAME"}, false),
						},
						"value": {
							Description: "IFA ILB name or private IP.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"port_ranges": {
							Description: "list of ports or port ranges. Values can be " +
								"mixed i.e. `[\"20\", \"100-200\"]`. An array with only the " +
								"value `-1` means any port.",
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

func resourceInternetApplicationCreate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateInternetApplicationRequest(d, m)

	if err != nil {
		return err
	}

	// Send request to create
	resource, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("provision_state", provisionState)

	return resourceInternetApplicationRead(d, m)
}

func resourceInternetApplicationRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))

	app, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", app.BillingTags)
	d.Set("connector_id", app.ConnectorId)
	d.Set("connector_type", app.ConnectorType)
	d.Set("fqdn_prefix", app.FqdnPrefix)
	d.Set("name", app.Name)
	d.Set("internet_protocol", app.InternetProtocol)
	d.Set("public_ips", app.PublicIps)
	d.Set("size", app.Size)

	// Segment
	segmentId, err := getSegmentIdByName(app.SegmentName, m)

	if err != nil {
		return err
	}

	d.Set("segment_id", segmentId)

	// targets
	var targets []map[string]interface{}

	for _, target := range app.Targets {
		i := map[string]interface{}{
			"type":       target.Type,
			"value":      target.Value,
			"ports":      target.Ports,
			"portRanges": target.PortRanges,
		}
		targets = append(targets, i)
	}

	d.Set("targets", targets)

	return nil
}

func resourceInternetApplicationUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))

	request, err := generateInternetApplicationRequest(d, m)

	if err != nil {
		return err
	}

	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)

	return resourceInternetApplicationRead(d, m)
}

func resourceInternetApplicationDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if provisionState != "SUCCESS" {
	}

	return nil
}

func generateInternetApplicationRequest(d *schema.ResourceData, m interface{}) (*alkira.InternetApplication, error) {

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// Targets
	//
	targets := expandInternetApplicationTargets(d.Get("target").(*schema.Set))

	// Assemble request
	request := &alkira.InternetApplication{
		BillingTags:          convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		ConnectorId:          d.Get("connector_id").(int),
		ConnectorType:        d.Get("connector_type").(string),
		FqdnPrefix:           d.Get("fqdn_prefix").(string),
		InboundConnectorId:   d.Get("inbound_connector_id").(string),
		InboundConnectorType: d.Get("inbound_connector_type").(string),
		InternetProtocol:     d.Get("internet_protocol").(string),
		Name:                 d.Get("name").(string),
		PublicIps:            convertTypeListToStringList(d.Get("public_ips").([]interface{})),
		SegmentName:          segmentName,
		Size:                 d.Get("size").(string),
		Targets:              targets,
	}

	return request, nil
}
