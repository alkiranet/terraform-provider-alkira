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

func resourceGroupUser(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	group := &alkira.UserGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	resource, _, err := api.Create(group)

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))

	return resourceGroupUserRead(d, m)
}

func resourceGroupUserRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	group, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", group.Name)
	d.Set("description", group.Description)

	return nil
}

func resourceGroupUserUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	group := &alkira.UserGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[INFO] Updating User Group (%s)", d.Id())
	_, err := api.Update(d.Id(), group)

	if err != nil {
		return err
	}

	return nil
}

func resourceGroupUserDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	_, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
