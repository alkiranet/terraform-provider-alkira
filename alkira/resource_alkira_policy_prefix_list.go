package alkira

import (
	"log"
	"strconv"

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
			"prefix_list_id": {
				Type:     schema.TypeInt,
				Computed: true,
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
	log.Printf("[INFO] Policy prefix list id: %d", id)

	if err != nil {
		log.Printf("[ERROR] failed to create prefix list")
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("prefix_list_id", id)

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
	id := d.Get("prefix_list_id").(int)

	log.Printf("[INFO] Deleting policy prefix list %d", id)
	err := client.DeletePolicyPrefixList(id)

	if err != nil {
		return err
	}

	return nil
}
