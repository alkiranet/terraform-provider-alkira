package alkira

import (
	"context"
	"fmt"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraMultiCloudGateway() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a Multi-Cloud Gateway. Create and update wait until the gateway is PROVISIONED. An external gateway (created in Lumen Connect rather than from Alkira) cannot be deleted from Alkira, so `terraform destroy` fails for an imported external gateway.",
		CreateContext: resourceMultiCloudGatewayCreate,
		ReadContext:   resourceMultiCloudGatewayRead,
		UpdateContext: resourceMultiCloudGatewayUpdate,
		DeleteContext: resourceMultiCloudGatewayDelete,
		Importer:      &schema.ResourceImporter{StateContext: importWithReadValidation(resourceMultiCloudGatewayRead)},
		Timeouts:      &schema.ResourceTimeout{Create: schema.DefaultTimeout(20 * time.Minute), Update: schema.DefaultTimeout(20 * time.Minute)},
		Schema: map[string]*schema.Schema{
			"name":                   {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Name of the gateway. Changing the name creates a new gateway."},
			"tier":                   {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"10GBPS", "50GBPS", "UNLIMITED"}, false), Description: "Gateway bandwidth tier: `10GBPS`, `50GBPS`, or `UNLIMITED`. Can be updated in place."},
			"description":            {Type: schema.TypeString, Optional: true, Description: "Description of the gateway. Can be updated in place."},
			"billing_account_number": {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true, Description: "The Lumen billing account number (BAN) this gateway bills to. When omitted, the tenant's default Lumen BAN is used. Changing it replaces the gateway."},
			"native_id":              {Type: schema.TypeString, Computed: true, Description: "Native gateway ID assigned by Lumen."},
			"state":                  {Type: schema.TypeString, Computed: true, Description: "Provisioning state of the gateway."},
			"external":               {Type: schema.TypeBool, Computed: true, Description: "True if the gateway was created in Lumen Connect rather than from Alkira."},
		},
	}
}

func multiCloudGatewayRequest(d *schema.ResourceData) *alkira.MultiCloudGateway {
	return &alkira.MultiCloudGateway{Name: d.Get("name").(string), Tier: d.Get("tier").(string), Description: d.Get("description").(string), BillingAccountNumber: d.Get("billing_account_number").(string)}
}

func setMultiCloudGateway(d *schema.ResourceData, gateway *alkira.MultiCloudGateway) diag.Diagnostics {
	return setMultiCloudGatewayFields(d, gateway, "name", "tier", "description", "billing_account_number", "native_id", "state", "external")
}

// setMultiCloudGatewayFields is shared with the data source, which sets only some of the keys.
func setMultiCloudGatewayFields(d *schema.ResourceData, gateway *alkira.MultiCloudGateway, keys ...string) diag.Diagnostics {
	values := map[string]string{"name": gateway.Name, "tier": gateway.Tier, "description": gateway.Description, "billing_account_number": gateway.BillingAccountNumber, "native_id": gateway.NativeId, "state": gateway.State}
	for _, key := range keys {
		if key == "external" {
			if err := d.Set(key, gateway.External); err != nil {
				return diag.FromErr(err)
			}
			continue
		}
		if err := d.Set(key, values[key]); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceMultiCloudGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewMultiCloudGateway(m.(*alkira.AlkiraClient))
	gateway, _, err, valErr, _ := api.Create(multiCloudGatewayRequest(d))
	if err != nil {
		return diag.FromErr(err)
	}
	if gateway.Id == "" {
		return diag.Errorf("creating multi-cloud gateway: API returned no ID")
	}
	d.SetId(gateway.Id.String())
	if valErr != nil {
		return diag.FromErr(valErr)
	}
	return waitForMultiCloudGateway(ctx, d, api, d.Timeout(schema.TimeoutCreate))
}

// waitForMultiCloudGateway polls while the gateway is PENDING or PROVISIONING; only the wait is bound to the timeout.
func waitForMultiCloudGateway(ctx context.Context, d *schema.ResourceData, api *alkira.AlkiraAPI[alkira.MultiCloudGateway], timeout time.Duration) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		gateway, _, err := api.GetById(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		switch gateway.State {
		case "PROVISIONED":
			return setMultiCloudGateway(d, gateway)
		case "PENDING", "PROVISIONING":
		default:
			return diag.Errorf("multi-cloud gateway %s entered %s state", d.Id(), gateway.State)
		}
		timer := time.NewTimer(5 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return diag.FromErr(fmt.Errorf("waiting for multi-cloud gateway %s to become PROVISIONED: %w", d.Id(), ctx.Err()))
		case <-timer.C:
		}
	}
}

func resourceMultiCloudGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	gateway, _, err := alkira.NewMultiCloudGateway(m.(*alkira.AlkiraClient)).GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}
	return setMultiCloudGateway(d, gateway)
}

func resourceMultiCloudGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewMultiCloudGateway(m.(*alkira.AlkiraClient))
	_, err, valErr, _ := api.Update(d.Id(), multiCloudGatewayRequest(d))
	if err != nil {
		return diag.FromErr(err)
	}
	if valErr != nil {
		return diag.FromErr(valErr)
	}
	return waitForMultiCloudGateway(ctx, d, api, d.Timeout(schema.TimeoutUpdate))
}

func resourceMultiCloudGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err, valErr, _ := alkira.NewMultiCloudGateway(m.(*alkira.AlkiraClient)).Delete(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("deleting multi-cloud gateway %s: %w", d.Id(), err))
	}
	if valErr != nil {
		return diag.FromErr(valErr)
	}
	d.SetId("")
	return nil
}
