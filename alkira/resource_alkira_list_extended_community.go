package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraListExtendedCommunity() *schema.Resource {
	return &schema.Resource{
		Description: "An extended community list matches a route when all values " +
			"in the list are present on the route. A route matches an extended " +
			"community list when any of the values match.",
		Create: resourceListExtendedCommunity,
		Read:   resourceListExtendedCommunityRead,
		Update: resourceListExtendedCommunityUpdate,
		Delete: resourceListExtendedCommunityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "name of the list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "description for the list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"values": {
				Description: "extended-community values to match on routes. Each " +
					"extended-community in this set is a structured tag value in " +
					"the format of `type:AA:NN` format (where AA is `0-65535` and " +
					"NN is `0-4294967295`) `AA` denotes a AS number or it could be " +
					"in the format of `IPaddr:nn` where IPaddr is a `x.x.x.x` IPv4 " +
					"address and nn is a 2 byte value `0-65535`. Type will only be" +
					"`soo` for now.",
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceListExtendedCommunity(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	// Construct request
	list, err := generateListExtendedCommunityRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	resource, _, err := api.Create(list)

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return resourceListExtendedCommunityRead(d, m)
}

func resourceListExtendedCommunityRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	list, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("values", list.Values)

	return nil
}

func resourceListExtendedCommunityUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	// Construct request
	list, err := generateListExtendedCommunityRequest(d, m)

	if err != nil {
		return err
	}

	// Send request to update
	_, err = api.Update(d.Id(), list)

	if err != nil {
		return err
	}

	return resourceListExtendedCommunityRead(d, m)
}

func resourceListExtendedCommunityDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	_, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func generateListExtendedCommunityRequest(d *schema.ResourceData, m interface{}) (*alkira.List, error) {

	values := convertTypeListToStringList(d.Get("values").([]interface{}))

	request := &alkira.List{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      values,
	}

	return request, nil
}
