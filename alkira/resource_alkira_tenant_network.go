package alkira

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlkiraTenantNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceTenantNetworkCreate,
		Read:   resourceTenantNetworkRead,
		Update: resourceTenantNetworkUpdate,
		Delete: resourceTenantNetworkDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTenantNetworkCreate(d *schema.ResourceData, m interface{}) error {
        return resourceTenantNetworkRead(d, m)
}

func resourceTenantNetworkRead(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceTenantNetworkUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceTenantNetworkRead(d, m)
}

func resourceTenantNetworkDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}
