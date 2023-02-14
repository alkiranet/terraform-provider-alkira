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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

	// Construct request
	list, err := generateListCommunityRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	resource, _, err := api.Create(list)

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return resourceListCommunityRead(d, m)
}

func resourceListCommunityRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

	list, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("values", list.Values)

	return nil
}

func resourceListCommunityUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

	list, err := generateListCommunityRequest(d, m)

	if err != nil {
		return err
	}

	_, err = api.Update(d.Id(), list)

	if err != nil {
		return err
	}

	return resourceListCommunityRead(d, m)
}

func resourceListCommunityDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

	_, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	return nil
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
