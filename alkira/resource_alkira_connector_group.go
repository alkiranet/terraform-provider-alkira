package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraGroupConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Manage connector groups\n\n" +
			"Connector Groups can contain one or many connectors across different segments. " +
			"This grouping of connectors can be for policy enforcement purposes or " +
			"for monitoring purposes within the network. It allows for easier policy " +
			"assignment by assigning policies to the entire group at the same time " +
			"instead of having to assign them individually.",
		Create: resourceGroupConnector,
		Read:   resourceGroupConnectorRead,
		Update: resourceGroupConnectorUpdate,
		Delete: resourceGroupConnectorDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the connector group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceGroupConnector(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	id, err := client.CreateConnectorGroup(d.Get("name").(string), d.Get("description").(string))

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceGroupConnectorRead(d, meta)
}

func resourceGroupConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorGroupById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", group.Name)
	d.Set("description", group.Description)

	return nil
}

func resourceGroupConnectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Updating Connector Group (%s)", d.Id())
	err := client.UpdateConnectorGroup(d.Id(), d.Get("name").(string), d.Get("description").(string))

	if err != nil {
		return err
	}

	return nil
}

func resourceGroupConnectorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteConnectorGroup(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
