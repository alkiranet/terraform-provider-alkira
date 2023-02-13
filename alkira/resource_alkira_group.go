package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Provide group resource.",
		Create:      resourceGroup,
		Read:        resourceGroupRead,
		Update:      resourceGroupUpdate,
		Delete:      resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceGroup(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	group := &alkira.Group{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[INFO] Group Creating")
	resource, _, err := api.Create(group)

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))

	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	group, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", group.Name)
	d.Set("description", group.Description)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	group := &alkira.Group{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[INFO] Updating Group (%s)", d.Id())
	_, err := api.Update(d.Id(), group)

	if err != nil {
		return err
	}

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	log.Printf("[INFO] Deleting Group (%s)", d.Id())
	_, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Group (%s)", d.Id())
	d.SetId("")
	return nil
}
