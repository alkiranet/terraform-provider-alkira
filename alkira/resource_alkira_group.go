package alkira

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroup,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group",
			},
		},
	}
}

func resourceGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	name   := d.Get("name").(string)

	log.Printf("[INFO] Group Creating")
	id, statusCode := client.CreateGroup(name)
	log.Printf("[INFO] Group ID: %d", id)

	if statusCode != 200 {
		fmt.Printf("ERROR: failed to create group")
	}

	d.SetId(strconv.Itoa(id))
	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client    := meta.(*alkira.AlkiraClient)
	groupId := d.Id()

	log.Printf("[INFO] Deleting Group %s", groupId)
	statusCode := client.DeleteGroup(groupId)

	if statusCode != 202 {
	 	return fmt.Errorf("failed to delete group %s", groupId)
	}

	return nil
}
