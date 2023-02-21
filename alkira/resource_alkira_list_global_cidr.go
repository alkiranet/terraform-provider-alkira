package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraListGlobalCidr() *schema.Resource {
	return &schema.Resource{
		Description: "A list of CIDRs to be used for services.",
		Create:      resourceListGlobalCidr,
		Read:        resourceListGlobalCidrRead,
		Update:      resourceListGlobalCidrUpdate,
		Delete:      resourceListGlobalCidrDelete,
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
			"cxp": {
				Description: "CXP the list belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"values": {
				Description: "A list of CIDRs, The CIDR must be `/24` and a " +
					"subnet of the following: `10.0.0.0/18`, `172.16.0.0/12`, " +
					"`192.168.0.0/16`, `100.64.0.0/10`.",
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Description: "A list of associated service types.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceListGlobalCidr(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	// Construct requst
	request := generateListGlobalCidrRequest(d, m)

	// Send request
	resource, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	d.SetId(string(resource.Id))
	return resourceListGlobalCidrRead(d, m)
}

func resourceListGlobalCidrRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	list, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("cxp", list.CXP)
	d.Set("values", list.Values)
	d.Set("tags", list.Tags)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourceListGlobalCidrUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateListGlobalCidrRequest(d, m)

	// Send request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceListGlobalCidrRead(d, m)
}

func resourceListGlobalCidrDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete list_global_cidr %s, provision failed", d.Id())
	}

	d.SetId("")
	return nil
}

func generateListGlobalCidrRequest(d *schema.ResourceData, m interface{}) *alkira.GlobalCidrList {

	request := &alkira.GlobalCidrList{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CXP:         d.Get("cxp").(string),
		Values:      convertTypeListToStringList(d.Get("values").([]interface{})),
		Tags:        convertTypeListToStringList(d.Get("tags").([]interface{})),
	}

	return request
}
