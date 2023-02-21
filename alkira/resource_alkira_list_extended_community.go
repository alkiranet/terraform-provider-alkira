package alkira

import (
	"context"
	"fmt"

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
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
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
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
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

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateListExtendedCommunityRequest(d, m)

	// Send request
	response, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	d.SetId(string(response.Id))
	return resourceListExtendedCommunityRead(d, m)
}

func resourceListExtendedCommunityRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	list, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("values", list.Values)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourceListExtendedCommunityUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateListExtendedCommunityRequest(d, m)

	// Send request to update
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceListExtendedCommunityRead(d, m)
}

func resourceListExtendedCommunityDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete list_extended_community %s, provision failed", d.Id())
	}

	d.SetId("")
	return nil
}

func generateListExtendedCommunityRequest(d *schema.ResourceData, m interface{}) *alkira.List {

	request := &alkira.List{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      convertTypeListToStringList(d.Get("values").([]interface{})),
	}

	return request
}
