package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlkiraTenantNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceTenantNetworkCreate,
		Read:   resourceTenantNetworkRead,
		Update: resourceTenantNetworkUpdate,
		Delete: resourceTenantNetworkDelete,

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

func resourceTenantNetworkCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	state, err := client.ProvisionTenantNetwork()

	if err != nil {
		log.Printf("[ERROR] Failed to provision tenant network: %s", state)
		return err
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

			log.Printf("[DEBUG] Tenant Network %d status received: %#v", client.TenantNetworkId, state)
			return state, state, nil
		},
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		log.Printf("[ERROR] Received error: %#v", err)
	}

	d.SetId(client.TenantNetworkId)
	return resourceTenantNetworkRead(d, m)
}

func resourceTenantNetworkRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTenantNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTenantNetworkRead(d, m)
}

func resourceTenantNetworkDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	state, err := client.ProvisionTenantNetwork()

	if err != nil {
		log.Printf("[ERROR] Failed to deprovision tenant network: %s", state)
		return err
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

			log.Printf("[DEBUG] Tenant Network %d status received: %#v", client.TenantNetworkId, state)
			return state, state, nil
		},
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return err
	}

	log.Printf("[INFO] Tenant Network %s deprovisioned", client.TenantNetworkId)

	d.SetId("")
	return nil
}
