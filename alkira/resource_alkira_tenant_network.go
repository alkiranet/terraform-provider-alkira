package alkira

import (
	"context"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraTenantNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantNetworkCreate,
		ReadContext:   resourceTenantNetworkRead,
		UpdateContext: resourceTenantNetworkUpdate,
		DeleteContext: resourceTenantNetworkDelete,

		Schema: map[string]*schema.Schema{
			"connectors": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"services": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"segments": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceTenantNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)

	state, err := client.ProvisionTenantNetwork()

	if err != nil {
		log.Printf("[ERROR] Failed to provision tenant network: %s", state)
		return diag.FromErr(err)
	}

	// Wait for tenant network provisoning to finish
	stateConf := &resource.StateChangeConf{
		Target:  []string{"SUCCESS"},
		Pending: []string{"IN_PROGRESS", "PENDING"},
		Timeout: d.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			state, err := client.GetTenantNetworkState()

			if err != nil {
				log.Printf("[ERROR] Received error: %#v", err)
				return state, "ERROR", err
			}

			log.Printf("[DEBUG] Tenant Network %s status received: %#v", client.TenantNetworkId, state)
			return state, state, nil
		},
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		log.Printf("[ERROR] Received error: %#v", err)
	}

	d.SetId(client.TenantNetworkId)
	return resourceTenantNetworkRead(ctx, d, m)
}

func resourceTenantNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceTenantNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTenantNetworkRead(ctx, d, m)
}

func resourceTenantNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)

	state, err := client.ProvisionTenantNetwork()

	if err != nil {
		log.Printf("[ERROR] Failed to deprovision tenant network: %s", state)
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Target:  []string{"SUCCESS"},
		Pending: []string{"IN_PROGRESS", "PENDING"},
		Timeout: d.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			state, err := client.GetTenantNetworkState()

			if err != nil {
				log.Printf("[ERROR] Received error: %#v", err)
				return state, "ERROR", err
			}

			log.Printf("[DEBUG] Tenant Network %s status received: %#v", client.TenantNetworkId, state)
			return state, state, nil
		},
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Tenant Network %s deprovisioned", client.TenantNetworkId)

	d.SetId("")
	return nil
}
