package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraServiceZscaler() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Zscaler firewall service. (**BETA**)",
		Create:      resourceZscaler,
		Read:        resourceZscalerRead,
		Update:      resourceZscalerUpdate,
		Delete:      resourceZscalerDelete,
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
			"connector_internet_exit_id": {
				//
				// NOTE: This field is included to ensure that teardown
				// of the zscaler service happens first.  By including
				// this field we are ensuring a dependency for the
				// alkira zscaler serivce.  Terraform destroys
				// dependencies first.
				//
				Description: "The ID of the `connector_internet_exit` " +
					"associated with the zscaler service.",
				Type:     schema.TypeString,
				Required: true,
			},
			"billing_tag_ids": {
				Description: "Billing tag IDs to associate with the service.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the Zscaler service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ipsec_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				Description: "The IPSEC tunnel configuration. This field " +
					"should only be set when `tunnel_type` is `IPSEC`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"esp_dh_group_number": {
							Description: "The IPSEC phase 2 DH Group to be " +
								"used. Input value must be either `MODP1024`" +
								"or `MODP2048`. The default value is `MODP1024`",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "MODP1024",
							ValidateFunc: validation.StringInSlice([]string{"MODP1024", "MODP2048"}, false),
						},
						"esp_encryption_algorithm": {
							Description: "The IPSEC phase 2 Encryption " +
								"Algorithm to be used. Input value must " +
								"be either `NULL` or `AES256CBC`. The " +
								"default value is `NULL`",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "NULL",
							ValidateFunc: validation.StringInSlice([]string{"NULL", "AES256CBC"}, false),
						},
						"esp_integrity_algorithm": {
							Description: "The IPSEC phase 2 Integrity " +
								"Algorithm to be used. Input value must " +
								"be either `MD5` or `SHA256`. The default " +
								"value is `MD5`.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "MD5",
							ValidateFunc: validation.StringInSlice([]string{"MD5", "SHA256"}, false),
						},
						"health_check_type": {
							Description: "The type of health check. Input " +
								"values must be either `IKE_STATUS` " +
								"`PING_PROBE` or `HTTP_PROBE`",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"IKE_STATUS", "PING_PROBE", "HTTP_PROBE"}, false),
						},
						"http_probe_url": {
							Description: "The url to check connection to " +
								"health, should be provided " +
								"when health check type is 'HTTP_PROBE'",
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_dh_group_number": {
							Description: "The IPSEC phase 1 DH Group to be " +
								"used. Input value must either be `MODP1024` " +
								"or `MODP2048`. The default is `MODP1024`",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "MODP1024",
							ValidateFunc: validation.StringInSlice([]string{"MODP1024", "MODP2048"}, false),
						},
						"ike_encryption_algorithm": {
							Description: "The IPSEC phase 1 Encryption " +
								"Algorithm to be used. Only `AES256CBC` " +
								"is allowed. The default value is `AES256CBC`.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "AES256CBC",
							ValidateFunc: validation.StringInSlice([]string{"AES256CBC"}, false),
						},
						"ike_integrity_algorithm": {
							Description: "The IPSEC phase 1 Integrity " +
								"Algorithm to be used. Only `SHA256` " +
								"is allowed. The default value is `SHA256`.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SHA256",
							ValidateFunc: validation.StringInSlice([]string{"SHA256"}, false),
						},
						"local_fpdn_id": {
							Description: "The local FQDN Id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"pre_shared_key": {
							Description: "The preshared key.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"ping_probe_ip": {
							Description: "The ping destination to check " +
								"connection health. It should be provided " +
								"when `health_check_type` is `PING_PROBE`",
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Description: "The name of the zscaler firewall.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"primary_public_edge_ip": {
				Description: "The IP for closest Zscaler PoP to `cxp` region.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"secondary_public_edge_ip": {
				Description: "The ip for standby Zscaler PoP to `cxp` region.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_ids": {
				Description: "IDs of segment associated with the service.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"size": {
				Description:  "The size of the service one of `SMALL`, `MEDIUM`, `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"tunnel_protocol": {
				Description:  "The type of tunnel protocol to be used to connect to Zscaler PoP.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{"IPSEC", "GRE"}, false),
			},
		},
	}
}

func resourceZscaler(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceZscaler(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateZscalerRequest(d, m)

	if err != nil {
		return err
	}

	// Send create request
	response, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	d.SetId(string(response.Id))
	return resourceZscalerRead(d, m)
}

func resourceZscalerRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceZscaler(m.(*alkira.AlkiraClient))

	// Get the service
	z, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	segmentIds, err := convertSegmentNamesToSegmentIds(z.Segments, m)

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", z.BillingTags)
	d.Set("cxp", z.Cxp)
	d.Set("description", z.Description)
	d.Set("ipsec_configuration", deflateZscalerIpsecConfiguration(z.IpsecConfiguration))
	d.Set("name", z.Name)
	d.Set("primary_public_edge_ip", z.PrimaryPublicEdgeIp)
	d.Set("secondary_public_edge_ip", z.SecondaryPublicEdgeIp)
	d.Set("segment_ids", segmentIds)
	d.Set("size", z.Size)
	d.Set("tunnel_protocol", z.TunnelType)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourceZscalerUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceZscaler(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateZscalerRequest(d, m)

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

	return resourceZscalerRead(d, m)
}

func resourceZscalerDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceZscaler(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete service_zscaler %s, provision failed", d.Id())
	}

	d.SetId("")
	return nil
}

// generateZscalerRequest generate service-zscaler request
func generateZscalerRequest(d *schema.ResourceData, m interface{}) (*alkira.ServiceZscaler, error) {

	cfgs, err := expandZscalerIpsecConfigurations(d.Get("ipsec_configuration").(*schema.Set))

	if err != nil {
		return nil, err
	}

	segmentNames, err := convertSegmentIdsToSegmentNames(d.Get("segment_ids").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	return &alkira.ServiceZscaler{
		BillingTags:           convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		Cxp:                   d.Get("cxp").(string),
		Description:           d.Get("description").(string),
		IpsecConfiguration:    cfgs,
		Name:                  d.Get("name").(string),
		PrimaryPublicEdgeIp:   d.Get("primary_public_edge_ip").(string),
		SecondaryPublicEdgeIp: d.Get("secondary_public_edge_ip").(string),
		Segments:              segmentNames,
		Size:                  d.Get("size").(string),
		TunnelType:            d.Get("tunnel_protocol").(string),
	}, nil
}
