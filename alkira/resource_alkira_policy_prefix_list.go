package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyPrefixList() *schema.Resource {
	return &schema.Resource{
		Description: "Manage policy prefix list.",
		Create:      resourcePolicyPrefixList,
		Read:        resourcePolicyPrefixListRead,
		Update:      resourcePolicyPrefixListUpdate,
		Delete:      resourcePolicyPrefixListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the prefix list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the prefix list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"prefixes": {
				Description: "A list of prefixes.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"prefix_range": {
				Description: "A valid prefix range that could be used to " +
					"define a prefix of type `ROUTE`.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Description: "A valid CIDR as prefix in " +
								"`x.x.x.x/m` format.",
							Type:     schema.TypeString,
							Required: true,
						},
						"le": {
							Description: "Integer less than `32` but " +
								"greater than mask `m` in prefix",
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ge": {
							Description: "Integer less than `32` but " +
								"greater than mask `m` in prefix and less than `le`.",
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourcePolicyPrefixList(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyPrefixListRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	response, _, err := api.Create(request)

	if err != nil {
		return err
	}

	d.SetId(string(response.Id))
	return resourcePolicyPrefixListRead(d, m)
}

func resourcePolicyPrefixListRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	list, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("prefixes", list.Prefixes)

	return nil
}

func resourcePolicyPrefixListUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyPrefixListRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	_, err = api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	return resourcePolicyPrefixListRead(d, m)
}

func resourcePolicyPrefixListDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if provisionState != "SUCCESS" {
	}

	d.SetId("")
	return nil
}

func generatePolicyPrefixListRequest(d *schema.ResourceData, m interface{}) (*alkira.PolicyPrefixList, error) {

	prefixes := convertTypeListToStringList(d.Get("prefixes").([]interface{}))
	prefixRanges, err := expandPrefixListPrefixRanges(d.Get("prefix_range").([]interface{}))

	if err != nil {
		return nil, err
	}

	list := &alkira.PolicyPrefixList{
		Description:  d.Get("description").(string),
		Name:         d.Get("name").(string),
		Prefixes:     prefixes,
		PrefixRanges: prefixRanges,
	}

	return list, nil
}
