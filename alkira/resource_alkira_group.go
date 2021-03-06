package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Manage groups\n\n" +
			"Groups can contain one or many connectors across different segments. " +
			"This grouping of connectors can be for policy enforcement purposes or " +
			"for monitoring purposes within the network. It allows for easier policy " +
			"assignment by assigning policies to the entire group at the same time " +
			"instead of having to assign them individually.",
		Create: resourceGroup,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Group Creating")
	id, err := client.CreateGroup(d.Get("name").(string), d.Get("description").(string))

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Group Updating")
	err := client.UpdateGroup(d.Id(), d.Get("name").(string), d.Get("description").(string))

	if err != nil {
		return err
	}

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Group %d", d.Id())
	err := client.DeleteGroup(d.Id())

	if err != nil {
		return err
	}

	return nil
}
