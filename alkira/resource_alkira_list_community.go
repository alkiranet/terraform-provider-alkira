package alkira

import (
	"context"
	"fmt"

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
				Description: "Name of the list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description for the list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
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
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateListCommunityRequest(d, m)

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
	return resourceListCommunityRead(d, m)
}

func resourceListCommunityRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

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

func resourceListCommunityUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateListCommunityRequest(d, m)

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceListCommunityRead(d, m)
}

func resourceListCommunityDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete list_community %s, provision failed", d.Id())
	}

	d.SetId("")
	return nil
}

func generateListCommunityRequest(d *schema.ResourceData, m interface{}) *alkira.List {

	request := &alkira.List{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      convertTypeListToStringList(d.Get("values").([]interface{})),
	}

	return request
}
