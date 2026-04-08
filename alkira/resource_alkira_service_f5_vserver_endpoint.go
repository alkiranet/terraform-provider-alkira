package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// f5VServerInstanceMetadataHash computes a stable hash for instance_metadata entries
// keyed on instance_id, which is the unique identifier per entry. This prevents
// spurious diffs when the API returns instance metadata in a different order.
var f5VServerInstanceMetadataHash = typeSetHash(func(m map[string]interface{}) string {
	return fmt.Sprintf("%v", m["instance_id"])
})

func resourceAlkiraServiceF5vServerEndpoint() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource for managing F5 vServer endpoint. (**BETA**)",
		CreateContext: resourceF5vServerEndpointCreate,
		ReadContext:   resourceF5vServerEndpointRead,
		UpdateContext: resourceF5vServerEndpointUpdate,
		DeleteContext: resourceF5vServerEndpointDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceF5vServerEndpointRead),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the F5 vServer Endpoint.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"f5_service_id": {
				Description: "ID of the F5 service associated with the" +
					"F5 vServer Endpoint.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"f5_service_instance_ids": {
				Description: "An array of F5 service instance IDs." +
					" A maximum of 2 instances are allowed.",
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"type": {
				Description: "The type of endpoint." +
					" Can be `ELB` or `ILB`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ELB", "ILB"}, false),
			},
			"segment_id": {
				Description: "ID of the segment associated with" +
					" the endpoint.",
				Type:     schema.TypeString,
				Required: true,
			},
			"fqdn_prefix": {
				Description: "The FQDN prefix of the endpoint. Required when type is `ELB`",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"protocol": {
				Description: "The portocol used for the endpoint." +
					" Can be one of `TCP` or `UDP`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
			},
			"port_ranges": {
				Description: "An array of ports or port ranges." +
					" Values can be mixed i.e. ['20', '100-200']." +
					" An array with only the value '-1' means any port." +
					" Required when type is `ELB`",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"snat": {
				Description: "SNAT for the endpoint." +
					" Can be `AUTOMAP` or `NONE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTOMAP", "NONE"}, false),
			},
			"destination_endpoint_port_ranges": {
				Description: "An array of ports or port ranges." +
					" Values can be mixed i.e. ['20', '100-200']." +
					" An array with only the value '-1' means any port." +
					" Required when type is `ILB` and snat is `NONE`",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"destination_endpoint_ip_addresses": {
				Description: "An array of ip addresses. Required when type is `ILB` and snat is `NONE`",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance_metadata": {
				Description: "Per-instance metadata populated after provisioning." +
					"ELB instances populate: elastic_ip, secondary_ip, vlan, route_domain_id, ecmp_pool_name. " +
					"ILB instances populate: virtual_ip, route_domain_id, ecmp_pool_name. " +
					"If provisioning occurs out of band (e.g. via the Alkira portal), run " +
					"`terraform apply -refresh-only` to sync this data into state.",
				Type:     schema.TypeSet,
				Set:      f5VServerInstanceMetadataHash,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Description: "F5 service instance ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"elastic_ip": {
							Description: "Elastic IP address assigned to the instance. Populated for ELB type only.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"secondary_ip": {
							Description: "Secondary IP address assigned to the instance. Populated for ELB type only.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"virtual_ip": {
							Description: "Virtual IP address assigned to the instance. Populated for ILB type only.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"vlan": {
							Description: "VLAN assigned to the instance. Populated for ELB type only.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ecmp_pool_name": {
							Description: "ECMP pool name assigned to the instance.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"route_domain_id": {
							Description: "Route domain ID assigned to the instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"segment_id": {
							Description: "Segment ID associated with the instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceF5vServerEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewF5vServerEndpoint(client)

	request, err := generateRequestF5vServerEndpoint(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceF5vServerEndpointRead(ctx, d, m)
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
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "PROVISION (CREATE) FAILED",
					Detail:   fmt.Sprintf("%s", provErr),
				},
			}
		}
	}
	return resourceF5vServerEndpointRead(ctx, d, m)
}

func resourceF5vServerEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewF5vServerEndpoint(m.(*alkira.AlkiraClient))

	f5, provState, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}
	d.Set("name", f5.Name)
	d.Set("type", f5.Type)
	d.Set("fqdn_prefix", f5.FqdnPrefix)
	d.Set("protocol", f5.Protocol)
	d.Set("port_ranges", f5.PortRanges)
	d.Set("snat", f5.Snat)
	d.Set("f5_service_id", f5.F5ServiceId)
	d.Set("f5_service_instance_ids", f5.F5ServiceInstanceIds)

	segmentId, err := getSegmentIdByName(f5.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("segment_id", segmentId)

	// Set post-provisioning instance metadata
	if f5.Metadata != nil && f5.Metadata.InstanceToMetadata != nil {
		var metadataList []map[string]interface{}
		for instanceId, m := range f5.Metadata.InstanceToMetadata {
			metadataList = append(metadataList, map[string]interface{}{
				"instance_id":     instanceId,
				"elastic_ip":      m.ElasticIp,
				"secondary_ip":    m.SecondaryIp,
				"virtual_ip":      m.VirtualIp,
				"vlan":            m.Vlan,
				"ecmp_pool_name":  m.EcmpPoolName,
				"route_domain_id": int(m.RouteDomainId),
				"segment_id":      int(m.SegmentId),
			})
		}
		d.Set("instance_metadata", metadataList)
	}

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil

}

func resourceF5vServerEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewF5vServerEndpoint(m.(*alkira.AlkiraClient))

	request, err := generateRequestF5vServerEndpoint(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceF5vServerEndpointRead(ctx, d, m)
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

	if client.Provision {
		d.Set("provision_state", provState)
		if provState == "FAILED" {
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "PROVISION (UPDATE) FAILED",
					Detail:   fmt.Sprintf("%s", provErr),
				},
			}
		}
	}
	return nil
}

func resourceF5vServerEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewF5vServerEndpoint(m.(*alkira.AlkiraClient))

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_service_f5_vserver_endpoint (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_service_f5_vserver_endpoint (id=%s)", err, d.Id()))
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

func generateRequestF5vServerEndpoint(d *schema.ResourceData, m interface{}) (*alkira.F5vServerEndpoint, error) {

	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)
	if err != nil {
		return nil, err
	}
	request := &alkira.F5vServerEndpoint{
		Name:                 d.Get("name").(string),
		Segment:              segmentName,
		Type:                 d.Get("type").(string),
		F5ServiceId:          d.Get("f5_service_id").(int),
		F5ServiceInstanceIds: convertTypeSetToIntList(d.Get("f5_service_instance_ids").(*schema.Set)),
		FqdnPrefix:           d.Get("fqdn_prefix").(string),
		Protocol:             d.Get("protocol").(string),
		PortRanges:           convertTypeSetToStringList(d.Get("port_ranges").(*schema.Set)),
		Snat:                 d.Get("snat").(string),
	}

	destinationEndpointPortRanges, ok1 := d.Get("destination_endpoint_port_ranges").(*schema.Set)
	hasDEPortRanges := ok1 && destinationEndpointPortRanges.Len() > 0
	destinationEndpointIpAddresses, ok2 := d.Get("destination_endpoint_ip_addresses").(*schema.Set)
	hasDEIpAddresses := ok2 && destinationEndpointIpAddresses.Len() > 0

	if hasDEPortRanges || hasDEIpAddresses {
		request.DestinationEndpoints = &alkira.F5VServerDestinationEndpoints{}
	}
	if hasDEPortRanges {
		request.DestinationEndpoints.PortRanges = convertTypeSetToStringList(destinationEndpointPortRanges)
	}
	if hasDEIpAddresses {
		request.DestinationEndpoints.IpAddresses = convertTypeSetToStringList(destinationEndpointIpAddresses)
	}

	return request, nil
}
