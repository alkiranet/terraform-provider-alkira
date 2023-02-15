package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraListAsPath() *schema.Resource {
	return &schema.Resource{
		Description: "This list could be used in a policy rule, a route " +
			"will match successfully if any one value from the list is " +
			"included within the AS-PATH of the route.",
		Create: resourceListAsPath,
		Read:   resourceListAsPathRead,
		Update: resourceListAsPathUpdate,
		Delete: resourceListAsPathDelete,
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
				Description: "Value can be regular expression of AS PATH " +
					"or space sparated AS numbers. BGP regular expressions" +
					"are based on POSIX 1003.2 regular expressions.",
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceListAsPath(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListAsPath(m.(*alkira.AlkiraClient))

	// Construct request
	list, err := generateListAsPathRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] failed to generate list as path request")
		return err
	}

	// Send request
	resource, _, err := api.Create(list)

	if err != nil {
		log.Printf("[ERROR] failed to create list as path")
		return err
	}

	d.SetId(string(resource.Id))

	return resourceListAsPathRead(d, m)
}

func resourceListAsPathRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListAsPath(m.(*alkira.AlkiraClient))

	list, err := api.GetById(d.Id())

	if err != nil {
		log.Printf("[ERROR] failed to get list as path %s", d.Id())
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("values", list.Values)

	return nil
}

func resourceListAsPathUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListAsPath(m.(*alkira.AlkiraClient))

	list, err := generateListAsPathRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing list as path %s", d.Id())
	_, err = api.Update(d.Id(), list)

	if err != nil {
		return err
	}

	return resourceListAsPathRead(d, m)
}

func resourceListAsPathDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListAsPath(m.(*alkira.AlkiraClient))

	_, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func generateListAsPathRequest(d *schema.ResourceData, m interface{}) (*alkira.List, error) {

	values := convertTypeListToStringList(d.Get("values").([]interface{}))

	request := &alkira.List{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      values,
	}

	return request, nil
}
