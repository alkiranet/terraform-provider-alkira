package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorInternetExit() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Internet Exit Connector.",
		CreateContext: resourceConnectorInternetExitCreate,
		ReadContext:   resourceConnectorInternetExitRead,
		UpdateContext: resourceConnectorInternetExitUpdate,
		DeleteContext: resourceConnectorInternetExitDelete,
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
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"segment_id": {
				Description: "ID of segment associated with the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "The CXP where the connector should be " +
					"provisioned.",
				Type:     schema.TypeString,
				Required: true,
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
			"egress_ips": {
				Description: "The types of egress IPs to use with the connector. " +
					"Current options are `ALKIRA_PUBLIC_IP` or `BYOIP`. If `BYOIP` " +
					"is one of the options provided `byoip_id` must also be set.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ALKIRA_PUBLIC_IP", "BYOIP"}, false),
				},
			},
			"public_ip_number": {
				Description: "The number of the public IPs to the connector. " +
					"Default is `2`.",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			"byoip_id": {
				Description: "ID of the BYOIP to be associated with the " +
					"connector.",
				Type:     schema.TypeInt,
				Optional: true,
			},
			"byoip_public_ips": {
				Description: "Public IPs in BYOIP to be used to access the " +
					"connector. The number of public IPs must be equal to " +
					"`public_ip_number`.",
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"traffic_distribution_algorithm": {
				Description: "The type of the algorithm to be used for traffic distribution." +
					"Currently, only `HASHING` is supported.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HASHING",
				ValidateFunc: validation.StringInSlice([]string{"HASHING"}, false),
			},
			"traffic_distribution_algorithm_attribute": {
				Description: "The attributes depends on the algorithm. For now, " +
					"it's either `DEFAULT` or `SRC_IP`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "SRC_IP"}, false),
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceConnectorInternetExitCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorInternet(m.(*alkira.AlkiraClient))

	request, err := generateConnectorInternetRequest(d, m)

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

	return resourceConnectorInternetExitRead(ctx, d, m)
}

func resourceConnectorInternetExitRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorInternet(m.(*alkira.AlkiraClient))

	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("byoip_id", connector.ByoipId)
	d.Set("byoip_public_ips", connector.PublicIps)
	d.Set("cxp", connector.CXP)
	d.Set("description", connector.Description)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("public_ip_number", connector.NumOfPublicIPs)
	d.Set("egress_ips", connector.EgressIpTypes)

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

	if connector.TrafficDistribution != nil {
		d.Set("traffic_distribution_algorithm", connector.TrafficDistribution.Algorithm)
		d.Set("traffic_distribution_algorithm_attribute", connector.TrafficDistribution.AlgorithmAttributes.Keys)
	}

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorInternetExitUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorInternet(m.(*alkira.AlkiraClient))

	request, err := generateConnectorInternetRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, provErr := api.Update(d.Id(), request)

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

	return nil
}

func resourceConnectorInternetExitDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorInternet(m.(*alkira.AlkiraClient))

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

// generateConnectorInternetRequest generate request for connector-internet
func generateConnectorInternetRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorInternet, error) {

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	algorithmAttributes := alkira.AlgorithmAttributes{
		Keys: d.Get("traffic_distribution_algorithm_attribute").(string),
	}

	trafficDistribution := alkira.TrafficDistribution{
		Algorithm:           d.Get("traffic_distribution_algorithm").(string),
		AlgorithmAttributes: algorithmAttributes,
	}

	if err != nil {
		return nil, err
	}

	request := &alkira.ConnectorInternet{
		BillingTags:         convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		ByoipId:             d.Get("byoip_id").(int),
		CXP:                 d.Get("cxp").(string),
		Description:         d.Get("description").(string),
		Group:               d.Get("group").(string),
		Enabled:             d.Get("enabled").(bool),
		PublicIps:           convertTypeSetToStringList(d.Get("byoip_public_ips").(*schema.Set)),
		Name:                d.Get("name").(string),
		NumOfPublicIPs:      d.Get("public_ip_number").(int),
		Segments:            []string{segmentName},
		TrafficDistribution: &trafficDistribution,
		EgressIpTypes:       convertTypeListToStringList(d.Get("egress_ips").([]interface{})),
	}

	return request, nil
}
