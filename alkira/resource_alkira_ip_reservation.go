package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraIpReservation() *schema.Resource {
	return &schema.Resource{
		Description:   "Provide IP Reservation resource.",
		CreateContext: resourceIpReservation,
		ReadContext:   resourceIpReservationRead,
		UpdateContext: resourceIpReservationUpdate,
		DeleteContext: resourceIpReservationDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceIpReservationRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the IP Reservation.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "The type of the IP Reservation. The value could be " +
					"either `PUBLIC` or `OVERLAY`. `PUBLIC` could be only created " +
					"by Alkira.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PUBLIC", "OVERLAY"}, false),
			},
			"prefix": {
				Description: "The IP Prefix of the IP Reservation. If this is " +
					"specified, both `prefix_type` and `prefix_len` will be " +
					"ignored.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"prefix_len": {
				Description: "The IP Prefix length of the IP Reservation.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"prefix_type": {
				Description: "The IP Prefix type of the IP Reservation. The " +
					"value could be `SEGMENT`, `APIPA`, `AZURE_APIPA` and " +
					"`PUBLIC`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"first_ip_assignment": {
				Description: "The value could be either `CUSTOMER` or `CXP`. " +
					"This is required when `prefix_len` is `30` or the " +
					"`prefix` is a `/30`. This field determines which IP from " +
					"the given or the computed `/30` prefix is assigned to the " +
					"customer end of the tunnel and which IP is assigned to the " +
					"CXP end of the tunnel.",
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CUSTOMER", "CXP"}, false),
			},
			"node_id": {
				Description: "The ID of the node that the IP Reservation is " +
					"assigned to. This must be provided when the given or " +
					"computed `prefix` is `/30`. When the `prefix` is `/32`" +
					"then this field determines whether the IP address will " +
					"be assigned to the customer end or the CXP end.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"cxp": {
				Description: "The CXP of the IP Reservation.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"scale_group_id": {
				Description: "The ID of the Scale Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_id": {
				Description: "The segment ID which the IP Reservation is to be " +
					"used.",
				Type:     schema.TypeString,
				Required: true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceIpReservation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewIPReservation(m.(*alkira.AlkiraClient))

	// Segment
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Construct request
	request := &alkira.IPReservation{
		Name:              d.Get("name").(string),
		Type:              d.Get("type").(string),
		Prefix:            d.Get("prefix").(string),
		PrefixLen:         d.Get("prefix_len").(int),
		PrefixType:        d.Get("prefix_type").(string),
		FirstIpAssignedTo: d.Get("first_ip_assignment").(string),
		NodeId:            d.Get("node_id").(string),
		Cxp:               d.Get("cxp").(string),
		ScaleGroupId:      d.Get("scale_group_id").(string),
		Segment:           segmentName,
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.Id)

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceIpReservationRead(ctx, d, m)
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

	return resourceIpReservationRead(ctx, d, m)
}

func resourceIpReservationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewIPReservation(m.(*alkira.AlkiraClient))

	// Get
	reservation, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", reservation.Name)
	d.Set("type", reservation.Type)
	d.Set("prefix", reservation.Prefix)
	d.Set("prefix_len", reservation.PrefixLen)
	d.Set("prefix_type", reservation.PrefixType)
	d.Set("first_ip_assignment", reservation.FirstIpAssignedTo)
	d.Set("node_id", reservation.NodeId)
	d.Set("cxp", reservation.Cxp)
	d.Set("scale_group_id", reservation.ScaleGroupId)

	// Set segment
	segmentId, err := getSegmentIdByName(reservation.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("segment_id", segmentId)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceIpReservationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewIPReservation(m.(*alkira.AlkiraClient))

	// Segment
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Construct request
	request := &alkira.IPReservation{
		Name:              d.Get("name").(string),
		Type:              d.Get("type").(string),
		Prefix:            d.Get("prefix").(string),
		PrefixLen:         d.Get("prefix_len").(int),
		PrefixType:        d.Get("prefix_type").(string),
		FirstIpAssignedTo: d.Get("first_ip_assignment").(string),
		NodeId:            d.Get("node_id").(string),
		Cxp:               d.Get("cxp").(string),
		ScaleGroupId:      d.Get("scale_group_id").(string),
		Segment:           segmentName,
	}

	// Send update request
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceIpReservationRead(ctx, d, m)
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

func resourceIpReservationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewIPReservation(m.(*alkira.AlkiraClient))

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	d.SetId("")

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
