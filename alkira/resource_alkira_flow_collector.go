package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraFlowCollector() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage flow collector.",
		CreateContext: resourceFlowCollector,
		ReadContext:   resourceFlowCollectorRead,
		UpdateContext: resourceFlowCollectorUpdate,
		DeleteContext: resourceFlowCollectorDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceFlowCollectorRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the flow collector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the flow collector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"collector_type": {
				Description:  "The type of the flow collector.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "GENERIC",
				ValidateFunc: validation.StringInSlice([]string{"GENERIC"}, false),
			},
			"enabled": {
				Description: "Whether the flow collector is enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"segment_id": {
				Description: "The segment on which flow export destination is " +
					"reachable. This should not be specified when destination " +
					"is reachable via internet. Also, segment can only be used " +
					"when `destination_ip` is provided, `destination_fqdn` is " +
					"not supported.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_ip": {
				Description: "The destination IP of the flow collector where " +
					"flow would be sent. Either `destination_ip` or " +
					"`destination_fqdn` are required.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_fqdn": {
				Description: "The destination FQDN of the flow collector where " +
					"flow would be sent. Either `destination_ip` or " +
					"`destination_fqdn` are required.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_port": {
				Description: "The destination port of the flow collector where " +
					"flow would be sent.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"transport_protocol": {
				Description: "The transport protocol to send the flow records " +
					"to destination.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "UDP",
				ValidateFunc: validation.StringInSlice([]string{"UDP"}, false),
			},
			"export_type": {
				Description: "The flow records export type. Only `IPFIX` is " +
					"supported for now.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPFIX",
				ValidateFunc: validation.StringInSlice([]string{"IPFIX"}, false),
			},
			"flow_record_template_id": {
				Description: "The flow records template ID. Currently only " +
					"default template ID `1` is supported",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"cxps": {
				Description: "A list of CXPs where the collector should be " +
					"provisioned for flow collecting.",
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceFlowCollector(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewFlowCollector(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateFlowCollectorRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceFlowCollectorRead(ctx, d, m)
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

	// Set provision state
	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceFlowCollectorRead(ctx, d, m)
}

func resourceFlowCollectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewFlowCollector(m.(*alkira.AlkiraClient))

	// Get
	flowCollector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", flowCollector.Name)
	d.Set("description", flowCollector.Description)
	d.Set("collector_type", flowCollector.CollectorType)
	d.Set("enabled", flowCollector.Enabled)
	d.Set("destination_ip", flowCollector.DestinationIp)
	d.Set("destination_fqdn", flowCollector.DestinationFqdn)
	d.Set("destination_port", flowCollector.DestinationPort)
	d.Set("transport_protocol", flowCollector.TransportProtocol)
	d.Set("export_type", flowCollector.ExportType)
	d.Set("flow_record_template_id", flowCollector.FlowRecordTemplateId)
	d.Set("cxps", flowCollector.Cxps)

	var segmentId string
	if flowCollector.Segment != "" {
		segmentId, err = getSegmentIdByName(flowCollector.Segment, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.Set("segment_id", segmentId)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceFlowCollectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewFlowCollector(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateFlowCollectorRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceFlowCollectorRead(ctx, d, m)
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

func resourceFlowCollectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewFlowCollector(m.(*alkira.AlkiraClient))

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	// Handle validation error
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateFlowCollectorRequest(d *schema.ResourceData, m interface{}) (*alkira.FlowCollector, error) {

	var segmentName string
	var err error
	segmentId := d.Get("segment_id")
	if segmentId != "" {
		segmentName, err = getSegmentNameById(d.Get("segment_id").(string), m)
	}

	if err != nil {
		return nil, err
	}

	request := &alkira.FlowCollector{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		CollectorType:        d.Get("collector_type").(string),
		Enabled:              d.Get("enabled").(bool),
		DestinationIp:        d.Get("destination_ip").(string),
		DestinationFqdn:      d.Get("destination_fqdn").(string),
		DestinationPort:      d.Get("destination_port").(int),
		TransportProtocol:    d.Get("transport_protocol").(string),
		ExportType:           d.Get("export_type").(string),
		FlowRecordTemplateId: d.Get("flow_record_template_id").(int),
		Cxps:                 convertTypeListToStringList(d.Get("cxps").([]interface{})),
		Segment:              segmentName,
	}

	return request, nil
}
