package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyPrefixList() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyPrefixList,
		Read:   resourcePolicyPrefixListRead,
		Update: resourcePolicyPrefixListUpdate,
		Delete: resourcePolicyPrefixListDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prefixes": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePolicyPrefixList(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	prefixes := convertTypeListToStringList(d.Get("prefixes").([]interface{}))

	list := &alkira.PolicyPrefixList{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		Prefixes:    prefixes,
	}

	log.Printf("[INFO] Policy prefix list Creating")
	id, err := client.CreatePolicyPrefixList(list)

	if err != nil {
		log.Printf("[ERROR] failed to create prefix list")
		return err
	}

	d.SetId(id)
	return resourcePolicyPrefixListRead(d, meta)
}

func resourcePolicyPrefixListRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePolicyPrefixListUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourcePolicyPrefixListRead(d, meta)
}

func resourcePolicyPrefixListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting policy prefix list %s", d.Id())
	return client.DeletePolicyPrefixList(d.Id())
}
