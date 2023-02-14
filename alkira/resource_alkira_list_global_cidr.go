package alkira

import (
	"log"

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
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	// Construct requst
	list, err := generateListGlobalCidrRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	resource, _, err := api.Create(list)

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return resourceListGlobalCidrRead(d, m)
}

func resourceListGlobalCidrRead(d *schema.ResourceData, m interface{}) error {
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

	return nil
}

func resourceListGlobalCidrUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	// Construct request
	list, err := generateListGlobalCidrRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	_, err = api.Update(d.Id(), list)

	if err != nil {
		return err
	}

	return resourceListGlobalCidrRead(d, m)
}

func resourceListGlobalCidrDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewGlobalCidrList(m.(*alkira.AlkiraClient))

	_, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	return nil
}

func generateListGlobalCidrRequest(d *schema.ResourceData, m interface{}) (*alkira.GlobalCidrList, error) {

	values := convertTypeListToStringList(d.Get("values").([]interface{}))
	tags := convertTypeListToStringList(d.Get("tags").([]interface{}))

	request := &alkira.GlobalCidrList{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CXP:         d.Get("cxp").(string),
		Values:      values,
		Tags:        tags,
	}

	return request, nil
}
