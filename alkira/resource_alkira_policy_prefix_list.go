package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyPrefixList() *schema.Resource {
	return &schema.Resource{
		Description: "Manage policy prefix list.",
		Create:      resourcePolicyPrefixList,
		Read:        resourcePolicyPrefixListRead,
		Update:      resourcePolicyPrefixListUpdate,
		Delete:      resourcePolicyPrefixListDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the prefix list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the prefix list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"prefixes": {
				Description: "A list of prefixes.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePolicyPrefixList(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyPrefixListRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] failed to generate prefix list request")
		return err
	}

	id, err := client.CreatePolicyPrefixList(request)

	if err != nil {
		log.Printf("[ERROR] failed to create prefix list")
		return err
	}

	d.SetId(id)
	return resourcePolicyPrefixListRead(d, m)
}

func resourcePolicyPrefixListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	list, err := client.GetPolicyPrefixListById(d.Id())

	if err != nil {
		log.Printf("[ERROR] Failed to get policy prefix list %s", d.Id())
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("prefixes", list.Prefixes)

	return nil
}

func resourcePolicyPrefixListUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyPrefixListRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy prefix list request")
		return err
	}

	err = client.UpdatePolicyPrefixList(d.Id(), request)

	if err != nil {
		log.Printf("[ERROR] Failed to update policy prefix list %s", d.Id())
		return err
	}

	return resourcePolicyPrefixListRead(d, m)
}

func resourcePolicyPrefixListDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeletePolicyPrefixList(d.Id())
}

func generatePolicyPrefixListRequest(d *schema.ResourceData, m interface{}) (*alkira.PolicyPrefixList, error) {

	prefixes := convertTypeListToStringList(d.Get("prefixes").([]interface{}))

	list := &alkira.PolicyPrefixList{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		Prefixes:    prefixes,
	}

	return list, nil
}
