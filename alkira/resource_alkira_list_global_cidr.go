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
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceListGlobalCidr(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	list, err := generateListGlobalCidrRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] failed to generate list global CIDR request")
		return err
	}

	id, err := client.CreateGlobalCidrList(list)

	if err != nil {
		log.Printf("[ERROR] failed to create list global CIDR")
		return err
	}

	d.SetId(id)
	return resourceListGlobalCidrRead(d, m)
}

func resourceListGlobalCidrRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	list, err := client.GetGlobalCidrListById(d.Id())

	if err != nil {
		log.Printf("[ERROR] failed to get list global CIDR %s", d.Id())
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("cxp", list.CXP)
	d.Set("values", list.Values)

	return nil
}

func resourceListGlobalCidrUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	list, err := generateListGlobalCidrRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing global CIDR list %s", d.Id())
	err = client.UpdateGlobalCidrList(d.Id(), list)

	return resourceListGlobalCidrRead(d, m)
}

func resourceListGlobalCidrDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting global CIDR list %s", d.Id())
	return client.DeleteGlobalCidrList(d.Id())
}

func generateListGlobalCidrRequest(d *schema.ResourceData, m interface{}) (*alkira.GlobalCidrList, error) {

	values := convertTypeListToStringList(d.Get("values").([]interface{}))

	request := &alkira.GlobalCidrList{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CXP:         d.Get("cxp").(string),
		Values:      values,
	}

	return request, nil
}
