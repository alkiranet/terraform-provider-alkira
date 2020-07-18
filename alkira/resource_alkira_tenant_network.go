package alkira

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/terraform-provider-alkira/alkira/internal"
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
				Required: true,
			},
			"segments": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceTenantNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*internal.AlkiraClient)

	client.ProvisionTenantNetwork()
	return resourceTenantNetworkRead(d, meta)
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
