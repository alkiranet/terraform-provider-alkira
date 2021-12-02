package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraListCommunity() *schema.Resource {
	return &schema.Resource{
		Description: "This list could be used to matches a route when all " +
			"values in the list are present on the route. A route matches " +
			"a list when any of the values match.",
		Create: resourceListCommunity,
		Read:   resourceListCommunityRead,
		Update: resourceListCommunityUpdate,
		Delete: resourceListCommunityDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description for the list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"values": {
				Description: "A list of communities to match on routes. Each " +
					"community in the list is a tag value in the format of " +
					"`AA:NN` format (where AA and NN are `0-65535`). AA " +
					"denotes a AS number.",
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceListCommunity(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	list, err := generateListCommunityRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] failed to generate list community request")
		return err
	}

	id, err := client.CreateList(list, alkira.ListTypeCommunity)

	if err != nil {
		log.Printf("[ERROR] failed to create list community")
		return err
	}

	d.SetId(id)
	return resourceListCommunityRead(d, m)
}

func resourceListCommunityRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	list, err := client.GetListById(d.Id(), alkira.ListTypeCommunity)

	if err != nil {
		log.Printf("[ERROR] failed to get list community %s", d.Id())
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("values", list.Values)

	return nil
}

func resourceListCommunityUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	list, err := generateListCommunityRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing list community %s", d.Id())
	err = client.UpdateList(d.Id(), list, alkira.ListTypeCommunity)

	return resourceListCommunityRead(d, m)
}

func resourceListCommunityDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting list community %s", d.Id())
	return client.DeleteList(d.Id(), alkira.ListTypeCommunity)
}

func generateListCommunityRequest(d *schema.ResourceData, m interface{}) (*alkira.List, error) {

	values := convertTypeListToStringList(d.Get("values").([]interface{}))

	request := &alkira.List{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      values,
	}

	return request, nil
}
