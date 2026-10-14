package alkira

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraMultiCloudGateway() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a Multi-Cloud Gateway. Creation waits until the gateway is PROVISIONED.",
		CreateContext: resourceMultiCloudGatewayCreate,
		ReadContext:   resourceMultiCloudGatewayRead,
		UpdateContext: resourceMultiCloudGatewayUpdate,
		DeleteContext: resourceMultiCloudGatewayDelete,
		Importer:      &schema.ResourceImporter{StateContext: importWithReadValidation(resourceMultiCloudGatewayRead)},
		Timeouts:      &schema.ResourceTimeout{Create: schema.DefaultTimeout(20 * time.Minute)},
		Schema: map[string]*schema.Schema{
			"name":        {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Name of the gateway. Changing the name creates a new gateway."},
			"tier":        {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"10GBPS", "50GBPS", "UNLIMITED"}, false), Description: "Gateway bandwidth tier: `10GBPS`, `50GBPS`, or `UNLIMITED`. Can be updated in place."},
			"description": {Type: schema.TypeString, Optional: true, Description: "Description of the gateway. Can be updated in place."},
			"native_id":   {Type: schema.TypeString, Computed: true, Description: "Native gateway ID assigned by Lumen."},
			"state":       {Type: schema.TypeString, Computed: true, Description: "Provisioning state of the gateway."},
		},
	}
}

func multiCloudGatewayRequest(d *schema.ResourceData) *alkira.MultiCloudGateway {
	return &alkira.MultiCloudGateway{Name: d.Get("name").(string), Tier: d.Get("tier").(string), Description: d.Get("description").(string)}
}

func setMultiCloudGateway(d *schema.ResourceData, gateway *alkira.MultiCloudGateway) diag.Diagnostics {
	for key, value := range map[string]string{"name": gateway.Name, "tier": gateway.Tier, "description": gateway.Description, "native_id": gateway.NativeId, "state": gateway.State} {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceMultiCloudGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	api := alkira.NewMultiCloudGateway(multiCloudGatewayClientWithContext(ctx, m.(*alkira.AlkiraClient)))
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
	for {
		if err := ctx.Err(); err != nil {
			return diag.FromErr(fmt.Errorf("waiting for multi-cloud gateway %s to become PROVISIONED: %w", d.Id(), err))
		}
		gateway, _, err := api.GetById(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		switch gateway.State {
		case "PROVISIONED":
			return setMultiCloudGateway(d, gateway)
		case "FAILED":
			return diag.Errorf("multi-cloud gateway %s entered FAILED state", d.Id())
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
		return diag.FromErr(err)
	}
	return setMultiCloudGateway(d, gateway)
}

func resourceMultiCloudGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewMultiCloudGateway(m.(*alkira.AlkiraClient))
	// Description must be present even when cleared; the client model omits empty descriptions.
	request := &multiCloudGatewayUpdateRequest{Name: d.Get("name").(string), Tier: d.Get("tier").(string), Description: d.Get("description").(string)}
	updateAPI := &alkira.AlkiraAPI[multiCloudGatewayUpdateRequest]{Client: api.Client, Uri: api.Uri, Provision: api.Provision}
	_, err, valErr, _ := updateAPI.Update(d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}
	if valErr != nil {
		return diag.FromErr(valErr)
	}
	return resourceMultiCloudGatewayRead(ctx, d, m)
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

// The client API does not accept contexts. Bind requests and retries to the
// create deadline without changing the shared provider client's configuration.
func multiCloudGatewayClientWithContext(ctx context.Context, original *alkira.AlkiraClient) *alkira.AlkiraClient {
	client := *original
	c := original.Client
	client.Client = &retryablehttp.Client{
		HTTPClient: c.HTTPClient, Logger: c.Logger,
		RetryWaitMin: c.RetryWaitMin, RetryWaitMax: c.RetryWaitMax, RetryMax: c.RetryMax,
		ResponseLogHook: c.ResponseLogHook, CheckRetry: c.CheckRetry, Backoff: c.Backoff,
		ErrorHandler: c.ErrorHandler, PrepareRetry: c.PrepareRetry,
		RequestLogHook: func(logger retryablehttp.Logger, request *http.Request, attempt int) {
			*request = *request.WithContext(ctx)
			if c.RequestLogHook != nil {
				c.RequestLogHook(logger, request, attempt)
			}
		},
	}
	return &client
}

type multiCloudGatewayUpdateRequest struct {
	Name        string `json:"name"`
	Tier        string `json:"tier"`
	Description string `json:"description"`
}
