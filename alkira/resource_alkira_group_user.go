package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraGroupUser() *schema.Resource {
	return &schema.Resource{
		Description: "Manage user groups\n\n",
		Create:      resourceGroupUser,
		Read:        resourceGroupUserRead,
		Update:      resourceGroupUserUpdate,
		Delete:      resourceGroupUserDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the user group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the user group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceGroupUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	id, err := client.CreateUserGroup(d.Get("name").(string), d.Get("description").(string))

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceGroupUserRead(d, meta)
}

func resourceGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetUserGroupById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", group.Name)
	d.Set("description", group.Description)

	return nil
}

func resourceGroupUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Updating User Group (%s)", d.Id())
	err := client.UpdateUserGroup(d.Id(), d.Get("name").(string), d.Get("description").(string))

	if err != nil {
		return err
	}

	return nil
}

func resourceGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteUserGroup(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
