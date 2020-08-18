package alkira

import (
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
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	name   := d.Get("name").(string)

	log.Printf("[INFO] Group Creating")
	id, err := client.CreateGroup(name)
	log.Printf("[INFO] Group ID: %d", id)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("group_id", id)

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
	groupId   := d.Get("group_id").(int)

	log.Printf("[INFO] Deleting Group %d", groupId)
	err := client.DeleteGroup(groupId)

	if err != nil {
	 	return err
	}

	return nil
}
